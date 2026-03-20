package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

const (
	colorReset       = "\033[0m"
	colorRed         = "\033[31m"
	colorGreen       = "\033[32m"
	colorYellow      = "\033[33m"
	colorBlue        = "\033[34m"
	colorMagenta     = "\033[35m"
	colorCyan        = "\033[36m"
	colorBold        = "\033[1m"
	colorDim         = "\033[2m"
	colorBrightRed    = "\033[91m"
	colorBrightYellow = "\033[93m"
	colorBrightCyan   = "\033[96m"
	colorBrightBlue   = "\033[94m"
	colorBrightMag    = "\033[95m"
	colorBgMagenta    = "\033[45m"
	colorBgRed        = "\033[41m"
	colorBgCyan       = "\033[46m"
	colorWhite        = "\033[97m"
	colorBlack        = "\033[30m"
)

var noColor bool

func init() {
	noColor = !isTerminal(os.Stdout)
}

func isTerminal(f *os.File) bool {
	var ws [4]uint16
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&ws)),
	)
	return errno == 0
}

func c(color, text string) string {
	if noColor {
		return text
	}
	return color + text + colorReset
}

func bold(text string) string        { return c(colorBold, text) }
func red(text string) string         { return c(colorRed, text) }
func green(text string) string       { return c(colorGreen, text) }
func yellow(text string) string      { return c(colorYellow, text) }
func blue(text string) string        { return c(colorBlue, text) }
func magenta(text string) string     { return c(colorMagenta, text) }
func cyan(text string) string        { return c(colorCyan, text) }
func dim(text string) string         { return c(colorDim, text) }
func white(text string) string        { return c(colorWhite, text) }
func brightRed(text string) string    { return c(colorBrightRed, text) }
func brightYellow(text string) string { return c(colorBrightYellow, text) }
func brightCyan(text string) string   { return c(colorBrightCyan, text) }
func brightBlue(text string) string   { return c(colorBrightBlue, text) }
func brightMag(text string) string    { return c(colorBrightMag, text) }

// badge renders text on a red background, like a pill label.
func badge(text string) string {
	if noColor {
		return text
	}
	return colorBold + colorBgRed + colorWhite + text + colorReset
}

func printSuccess(msg string) { fmt.Println(green("✓ ") + msg) }
func printError(msg string)   { fmt.Fprintln(os.Stderr, red("✗ ")+msg) }
func printWarn(msg string)    { fmt.Println(yellow("⚠ ") + msg) }
func printInfo(msg string)    { fmt.Println(cyan("→ ") + msg) }

func printHeader(msg string) {
	fmt.Println()
	fmt.Println(bold(cyan("● " + msg)))
	fmt.Println(dim(strings.Repeat("─", 40)))
}

var stdinReader = bufio.NewReader(os.Stdin)

func prompt(label string) string {
	fmt.Fprintf(os.Stderr, "%s %s", cyan("?"), bold(label+": "))
	line, _ := stdinReader.ReadString('\n')
	return strings.TrimSpace(line)
}

func confirm(label string) bool {
	fmt.Fprintf(os.Stderr, "%s %s %s ", cyan("?"), bold(label), dim("[y/N]"))
	line, _ := stdinReader.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes"
}

func promptMultiline(label string) string {
	fmt.Fprintf(os.Stderr, "%s %s\n%s\n",
		cyan("?"), bold(label+":"),
		dim("  (end with a line containing only '.')"))
	var lines []string
	for {
		fmt.Fprintf(os.Stderr, "  ")
		line, _ := stdinReader.ReadString('\n')
		line = strings.TrimRight(line, "\n")
		if line == "." {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// selectOption shows a numbered menu and returns the chosen string.
// Also accepts single-char shortcut matching the first char of an option.
func selectOption(label string, options []string) string {
	fmt.Fprintf(os.Stderr, "%s %s\n", cyan("?"), bold(label+":"))
	for i, opt := range options {
		fmt.Fprintf(os.Stderr, "  %s %s\n", dim(strconv.Itoa(i+1)+"."), opt)
	}
	for {
		fmt.Fprintf(os.Stderr, "  %s ", cyan("›"))
		line, _ := stdinReader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Try numeric
		if n, err := strconv.Atoi(line); err == nil && n >= 1 && n <= len(options) {
			return options[n-1]
		}
		// Try single-char shortcut
		for _, opt := range options {
			if len(line) == 1 && strings.HasPrefix(strings.ToLower(opt), line) {
				return opt
			}
		}
		fmt.Fprintf(os.Stderr, "  %s\n", red("Invalid choice, try again."))
	}
}
