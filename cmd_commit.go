package main

import (
	"fmt"
	"strings"
)

func cmdCommit(cfg *Config, args []string) error {
	// If -m is already provided, pass through
	for _, a := range args {
		if a == "-m" {
			return RunGit(cfg, append([]string{"commit"}, args...)...)
		}
	}

	// Check if there's anything staged
	staged, err := RunGitCapture(cfg, "diff", "--cached", "--name-only")
	if err != nil {
		return err
	}
	if strings.TrimSpace(staged) == "" {
		printWarn("Nothing staged. Use 'mantra add' first.")
		return nil
	}

	// Show staged files
	printHeader("Staged files")
	for _, f := range strings.Split(strings.TrimSpace(staged), "\n") {
		if f != "" {
			fmt.Println("  " + green("+ ") + f)
		}
	}
	fmt.Println()

	// Get commit message
	msg := promptMultiline("Commit message")
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return fmt.Errorf("commit message cannot be empty")
	}

	if err := RunGit(cfg, "commit", "-m", msg); err != nil {
		return err
	}
	printSuccess("Committed.")
	return nil
}
