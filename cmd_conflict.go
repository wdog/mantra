package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func cmdConflict(cfg *Config, args []string) error {
	out, _ := RunGitCapture(cfg, "diff", "--name-only", "--diff-filter=U")
	files := nonEmpty(strings.Split(strings.TrimSpace(out), "\n"))

	if len(files) == 0 {
		printSuccess("No conflicts detected.")
		return nil
	}

	printHeader(fmt.Sprintf("Conflicted files (%d)", len(files)))
	for _, f := range files {
		fmt.Println("  " + red("✗ ") + f)
	}

	resolved := 0
	skipped := 0
	options := []string{"edit", "ours", "theirs", "skip"}

	for i, file := range files {
		fmt.Printf("\n%s %s\n", bold(fmt.Sprintf("  [%d/%d]", i+1, len(files))), yellow(file))

		// Show a compact diff
		_ = RunGit(cfg, "diff", "--color=always", "--", file)

		choice := selectOption("Action", options)

		switch choice {
		case "edit":
			if err := openEditor(file); err != nil {
				printWarn("Could not open editor: " + err.Error())
				printInfo("Edit the file manually, then press Enter to continue.")
				prompt("Press Enter when done")
			}
			_ = RunGit(cfg, "add", "--", file)
			resolved++

		case "ours":
			if err := RunGit(cfg, "checkout", "--ours", "--", file); err != nil {
				return err
			}
			_ = RunGit(cfg, "add", "--", file)
			printSuccess("Kept ours: " + file)
			resolved++

		case "theirs":
			if err := RunGit(cfg, "checkout", "--theirs", "--", file); err != nil {
				return err
			}
			_ = RunGit(cfg, "add", "--", file)
			printSuccess("Kept theirs: " + file)
			resolved++

		case "skip":
			printWarn("Skipped: " + file)
			skipped++
		}
	}

	fmt.Println()
	printInfo(fmt.Sprintf("Resolved: %d  Skipped: %d", resolved, skipped))

	if skipped == 0 && resolved > 0 {
		if confirm("All conflicts resolved. Commit now?") {
			return cmdCommit(cfg, nil)
		}
	}
	return nil
}

func openEditor(file string) error {
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func nonEmpty(ss []string) []string {
	var out []string
	for _, s := range ss {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
