package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func cmdTimeline(cfg *Config, args []string) error {
	var relPath string

	if len(args) > 0 {
		// Resolve the provided path to a path relative to the work-tree
		abs, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(cfg.WorkTree, abs)
		if err != nil || strings.HasPrefix(rel, "..") {
			return fmt.Errorf("file is outside the work-tree: %s", args[0])
		}
		relPath = rel
	} else {
		// No file given — let the user pick from all tracked files
		out, err := RunGitCapture(cfg, "ls-files")
		if err != nil {
			return fmt.Errorf("cannot list tracked files: %w", err)
		}
		files := strings.Fields(strings.TrimSpace(out))
		if len(files) == 0 {
			printWarn("No tracked files yet.")
			return nil
		}
		relPath = selectOption("Select a file", files)
		if relPath == "" {
			return nil
		}
	}

	// Fetch the commit history for this file
	logOut, err := RunGitCapture(cfg,
		"log", "--format=%h  %ad  %s", "--date=short", "--", relPath)
	if err != nil || strings.TrimSpace(logOut) == "" {
		printWarn("No history found for: " + relPath)
		return nil
	}

	commits := strings.Split(strings.TrimRight(logOut, "\n"), "\n")
	if len(commits) == 0 {
		printWarn("No history found for: " + relPath)
		return nil
	}

	printHeader("Timeline: " + relPath)

	for {
		choice := selectOption("Select a commit to view (or Ctrl-C to exit)", commits)
		if choice == "" {
			break
		}

		// Extract the short hash (first field)
		hash := strings.Fields(choice)[0]
		showContent(cfg, hash, relPath)

		fmt.Println()
		fmt.Print(dim("  Press Enter to return to timeline..."))
		fmt.Scanln()
		fmt.Println()
	}

	return nil
}

// showContent prints the file content at the given commit via a pager if available.
func showContent(cfg *Config, hash, relPath string) {
	// Use git show <hash>:<path> — path must be relative to repo root (= work-tree)
	ref := hash + ":" + relPath

	pager := os.Getenv("PAGER")
	if pager == "" {
		if _, err := exec.LookPath("less"); err == nil {
			pager = "less"
		}
	}

	fmt.Println()
	fmt.Println("  " + bold(cyan(relPath)) + dim("  @"+hash))
	fmt.Println(dim("  " + strings.Repeat("─", 60)))
	fmt.Println()

	if pager != "" {
		// Stream through pager
		gitArgs := append(GitArgs(cfg), "show", ref)
		gitCmd := exec.Command("git", gitArgs...)
		pagerCmd := exec.Command(pager)
		pipe, err := gitCmd.StdoutPipe()
		if err == nil {
			pagerCmd.Stdin = pipe
			pagerCmd.Stdout = os.Stdout
			pagerCmd.Stderr = os.Stderr
			if err := gitCmd.Start(); err == nil {
				_ = pagerCmd.Start()
				_ = gitCmd.Wait()
				_ = pagerCmd.Wait()
				return
			}
		}
	}

	// Fallback: print directly
	_ = RunGit(cfg, "show", ref)
}
