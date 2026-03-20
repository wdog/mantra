package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type repoState struct {
	branch    string
	ahead     int
	behind    int
	staged    int
	modified  int
	untracked int
	conflicts int
}

func getRepoState(cfg *Config) repoState {
	s := repoState{}

	out, err := RunGitCapture(cfg, "status", "--short", "--branch")
	if err != nil {
		return s
	}

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	for i, line := range lines {
		if i == 0 && strings.HasPrefix(line, "##") {
			s.branch, s.ahead, s.behind = parseBranchLine(line)
			continue
		}
		if len(line) < 2 {
			continue
		}
		x := line[0]
		y := line[1]
		switch {
		case x == 'U' || y == 'U' || (x == 'A' && y == 'A') || (x == 'D' && y == 'D'):
			s.conflicts++
		case x != ' ' && x != '?':
			s.staged++
		case y == 'M' || y == 'D':
			s.modified++
		case x == '?' && y == '?':
			s.untracked++
		}
	}
	return s
}

func parseBranchLine(line string) (branch string, ahead, behind int) {
	line = strings.TrimPrefix(line, "## ")

	if idx := strings.Index(line, "["); idx != -1 {
		info := line[idx+1 : strings.Index(line, "]")]
		for _, part := range strings.Split(info, ",") {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "ahead ") {
				ahead, _ = strconv.Atoi(strings.TrimPrefix(part, "ahead "))
			}
			if strings.HasPrefix(part, "behind ") {
				behind, _ = strconv.Atoi(strings.TrimPrefix(part, "behind "))
			}
		}
		line = strings.TrimSpace(line[:idx])
	}

	if idx := strings.Index(line, "..."); idx != -1 {
		branch = line[:idx]
	} else {
		branch = line
	}
	return
}

func renderPrompt(s repoState) string {
	var parts []string

	branchName := s.branch
	if branchName == "" {
		branchName = "?"
	}
	parts = append(parts, bold(cyan(" "+branchName)))

	if s.ahead > 0 {
		parts = append(parts, green("↑"+strconv.Itoa(s.ahead)))
	}
	if s.behind > 0 {
		parts = append(parts, red("↓"+strconv.Itoa(s.behind)))
	}
	if s.conflicts > 0 {
		parts = append(parts, red("!"+strconv.Itoa(s.conflicts)))
	}
	if s.staged > 0 {
		parts = append(parts, green("+"+strconv.Itoa(s.staged)))
	}
	if s.modified > 0 {
		parts = append(parts, yellow("~"+strconv.Itoa(s.modified)))
	}
	if s.untracked > 0 {
		parts = append(parts, dim("?"+strconv.Itoa(s.untracked)))
	}

	return bold(cyan("mantra")) + " " + strings.Join(parts, " ") + " " + cyan("›") + " "
}

// replComplete returns completions for the given input prefix.
func replComplete(cfg *Config, input string) []string {
	// Commands that take file-path arguments
	// Order matters: more specific prefixes ("add -f ") must come before "add "
	// "add -u" and "modified" are complete commands — no path completion needed
	fileArgCmds := []string{"add -f ", "add ", "diff ", "timeline "}
	for _, cmd := range fileArgCmds {
		if strings.HasPrefix(input, cmd) {
			partial := input[len(cmd):]
			var result []string
			for _, p := range completeFilePath(cfg.WorkTree, partial) {
				result = append(result, cmd+p)
			}
			return result
		}
	}

	subCommands := map[string][]string{
		"stash ":      {"list", "pop"},
		"rebase ":     {"--abort", "--continue"},
		"completion ": {"fish", "bash", "zsh"},
	}
	for prefix, subs := range subCommands {
		if strings.HasPrefix(input, prefix) {
			rest := input[len(prefix):]
			var matches []string
			for _, s := range subs {
				if strings.HasPrefix(s, rest) {
					matches = append(matches, prefix+s)
				}
			}
			return matches
		}
	}

	topLevel := []string{
		"add", "add -u", "checkout", "commit", "completion", "conflict", "diff", "files",
		"help", "init", "log", "ls", "modified", "pull", "push",
		"rebase", "reset", "stash", "status", "timeline",
		"exit", "quit",
	}
	var matches []string
	for _, cmd := range topLevel {
		if strings.HasPrefix(cmd, input) {
			matches = append(matches, cmd)
		}
	}
	return matches
}

// completeFilePath returns file/dir completions under workTree for the given partial path.
func completeFilePath(workTree, partial string) []string {
	home, _ := os.UserHomeDir()

	// Expand partial to absolute for filesystem ops; remember display format
	abs := partial
	useTilde := false
	switch {
	case strings.HasPrefix(partial, "~/"):
		abs = filepath.Join(home, partial[2:])
		useTilde = true
	case partial == "~":
		abs = home
		useTilde = true
	case partial == "":
		abs = workTree
		useTilde = workTree == home
	case !filepath.IsAbs(partial):
		abs = filepath.Join(workTree, partial)
	}

	var dir, base string
	if partial == "" || partial == "~" || strings.HasSuffix(partial, "/") {
		dir = abs
		base = ""
	} else {
		dir = filepath.Dir(abs)
		base = filepath.Base(abs)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var matches []string
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") && !strings.HasPrefix(base, ".") {
			continue
		}
		if base != "" && !strings.HasPrefix(name, base) {
			continue
		}
		entryAbs := filepath.Join(dir, name)
		var completion string
		switch {
		case useTilde:
			rel, _ := filepath.Rel(home, entryAbs)
			completion = "~/" + rel
		case filepath.IsAbs(partial):
			completion = entryAbs
		default:
			rel, _ := filepath.Rel(workTree, entryAbs)
			completion = rel
		}
		if e.IsDir() {
			completion += "/"
		}
		matches = append(matches, completion)
	}
	sort.Strings(matches)
	return matches
}

func commonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
		}
	}
	return prefix
}

// termios raw mode support
type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]uint8
	Ispeed uint32
	Ospeed uint32
}

func getTermios(fd uintptr) (*termios, error) {
	t := &termios{}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(t)))
	if errno != 0 {
		return nil, errno
	}
	return t, nil
}

func setTermios(fd uintptr, t *termios) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCSETS, uintptr(unsafe.Pointer(t)))
	if errno != 0 {
		return errno
	}
	return nil
}

// readLine reads a line from stdin in raw mode so ESC exits immediately.
// Supports left/right cursor movement and tab completion via the complete func.
// promptStr is reprinted after displaying completions. Returns (line, eof, escPressed).
func readLine(promptStr string, complete func(string) []string) (string, bool, bool) {
	fd := os.Stdin.Fd()
	orig, err := getTermios(fd)
	if err != nil {
		// Fallback: not a tty, use normal line read
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return "", true, false
		}
		return scanner.Text(), false, false
	}

	raw := *orig
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0
	_ = setTermios(fd, &raw)
	defer setTermios(fd, orig)

	var buf []byte
	cursor := 0
	b := make([]byte, 1)

	for {
		n, err := os.Stdin.Read(b)
		if n == 0 || err != nil {
			return string(buf), true, false
		}
		ch := b[0]
		switch ch {
		case 0x09: // Tab — only complete when cursor is at end
			if cursor < len(buf) {
				fmt.Fprint(os.Stdout, "\x07")
				continue
			}
			cur := string(buf)
			matches := complete(cur)
			switch len(matches) {
			case 0:
				fmt.Fprint(os.Stdout, "\x07") // bell
			case 1:
				ext := matches[0][len(cur):]
				buf = append(buf, []byte(ext)...)
				cursor = len(buf)
				fmt.Fprint(os.Stdout, ext)
			default:
				cp := commonPrefix(matches)
				if len(cp) > len(cur) {
					ext := cp[len(cur):]
					buf = append(buf, []byte(ext)...)
					cursor = len(buf)
					fmt.Fprint(os.Stdout, ext)
				} else {
					fmt.Println()
					for _, m := range matches {
						fmt.Fprint(os.Stdout, "  "+cyan(m))
					}
					fmt.Println()
					fmt.Fprint(os.Stderr, promptStr)
					fmt.Fprint(os.Stdout, string(buf))
				}
			}
		case 0x1b: // ESC or escape sequence (arrow keys, etc.)
			// Peek with a short timeout to tell bare ESC from sequences like \x1b[D
			raw.Cc[syscall.VMIN] = 0
			raw.Cc[syscall.VTIME] = 1 // 100 ms
			_ = setTermios(fd, &raw)
			n2, _ := os.Stdin.Read(b)
			raw.Cc[syscall.VMIN] = 1
			raw.Cc[syscall.VTIME] = 0
			_ = setTermios(fd, &raw)

			if n2 == 0 {
				// Bare ESC — exit REPL
				return "", false, true
			}
			if b[0] != '[' {
				// Unknown escape sequence, ignore
				continue
			}
			n3, _ := os.Stdin.Read(b)
			if n3 == 0 {
				continue
			}
			switch b[0] {
			case 'D': // Left arrow
				if cursor > 0 {
					cursor--
					fmt.Fprint(os.Stdout, "\x1b[D")
				}
			case 'C': // Right arrow
				if cursor < len(buf) {
					cursor++
					fmt.Fprint(os.Stdout, "\x1b[C")
				}
			case 'A', 'B': // Up/Down — ignore (no history yet)
			}
		case 0x01: // Ctrl-A — move to start
			if cursor > 0 {
				fmt.Fprintf(os.Stdout, "\x1b[%dD", cursor)
				cursor = 0
			}
		case 0x05: // Ctrl-E — move to end
			if cursor < len(buf) {
				fmt.Fprintf(os.Stdout, "\x1b[%dC", len(buf)-cursor)
				cursor = len(buf)
			}
		case '\r', '\n':
			fmt.Println()
			return string(buf), false, false
		case 0x7f, 0x08: // Backspace / DEL
			if cursor > 0 {
				buf = append(buf[:cursor-1], buf[cursor:]...)
				cursor--
				suffix := buf[cursor:]
				fmt.Fprint(os.Stdout, "\b"+string(suffix)+" ")
				fmt.Fprintf(os.Stdout, "\x1b[%dD", len(suffix)+1)
			}
		case 0x04: // Ctrl-D
			fmt.Println()
			return "", true, false
		case 0x03: // Ctrl-C
			fmt.Println()
			buf = buf[:0]
			cursor = 0
			return "", false, false
		default:
			if ch >= 0x20 {
				// Insert at cursor position
				buf = append(buf, 0)
				copy(buf[cursor+1:], buf[cursor:])
				buf[cursor] = ch
				cursor++
				// Print char + remainder, then reposition cursor
				fmt.Fprint(os.Stdout, string(buf[cursor-1:]))
				if back := len(buf) - cursor; back > 0 {
					fmt.Fprintf(os.Stdout, "\x1b[%dD", back)
				}
			}
		}
	}
}

func runREPL(cfg *Config) {
	printBanner(cfg)

	stdinReader = bufio.NewReader(os.Stdin)

	for {
		state := getRepoState(cfg)
		prompt := renderPrompt(state)
		fmt.Fprint(os.Stderr, prompt)

		line, eof, esc := readLine(prompt, func(s string) []string { return replComplete(cfg, s) })
		if eof || esc {
			printInfo("Bye!")
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" || line == "q" {
			printInfo("Bye!")
			break
		}

		tokens := strings.Fields(line)

		var h handler
		var rest []string
		if len(tokens) >= 2 {
			if fn, ok := commands[tokens[0]+" "+tokens[1]]; ok {
				h = fn
				rest = tokens[2:]
			}
		}
		if h == nil {
			if fn, ok := commands[tokens[0]]; ok {
				h = fn
				rest = tokens[1:]
			}
		}

		if h == nil {
			fmt.Fprintln(os.Stderr, red("✗ Unknown command: ")+tokens[0])
			fmt.Fprintln(os.Stderr, dim("  Type 'help' for available commands."))
			continue
		}

		if err := h(cfg, rest); err != nil {
			msg := err.Error()
			if !strings.Contains(msg, "exit status") {
				fmt.Fprintln(os.Stderr, red("✗ ")+msg)
			}
		}

		fmt.Println()
		stdinReader = bufio.NewReader(os.Stdin)
	}
}
