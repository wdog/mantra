package main

import (
	"fmt"
	"strings"
)

func cmdStash(cfg *Config, args []string) error {
	msg := ""
	if len(args) > 0 {
		msg = strings.Join(args, " ")
	} else {
		msg = prompt("Stash message (optional, enter to skip)")
	}
	gitArgs := []string{"stash", "push"}
	if msg != "" {
		gitArgs = append(gitArgs, "-m", msg)
	}
	if err := RunGit(cfg, gitArgs...); err != nil {
		return err
	}
	printSuccess("Stashed.")
	return nil
}

func cmdStashPop(cfg *Config, args []string) error {
	if err := RunGit(cfg, "stash", "pop"); err != nil {
		if gitStateFile(cfg, "MERGE_HEAD") {
			printWarn("Conflict during stash pop. Run 'mantra conflict' to resolve.")
		}
		return err
	}
	printSuccess("Stash applied.")
	return nil
}

func cmdStashList(cfg *Config, args []string) error {
	return RunGit(cfg, "stash", "list")
}

func cmdRebase(cfg *Config, args []string) error {
	// Handle --continue / --abort
	for _, a := range args {
		if a == "--continue" {
			return RunGit(cfg, "rebase", "--continue")
		}
		if a == "--abort" {
			return RunGit(cfg, "rebase", "--abort")
		}
	}

	branch := ""
	if len(args) > 0 {
		branch = args[0]
	} else {
		branch = prompt("Branch to rebase onto")
	}
	if branch == "" {
		return fmt.Errorf("branch name required")
	}

	if err := RunGit(cfg, "rebase", branch); err != nil {
		if gitStateFile(cfg, "REBASE_HEAD") {
			printWarn("Conflict during rebase.")
			printInfo("Resolve conflicts, then run: mantra rebase --continue")
			printInfo("To abort: mantra rebase --abort")
		}
		return err
	}
	printSuccess("Rebased onto " + branch + ".")
	return nil
}

func cmdReset(cfg *Config, args []string) error {
	mode := ""
	ref := ""

	if len(args) >= 2 {
		mode = args[0]
		ref = args[1]
	} else if len(args) == 1 {
		mode = args[0]
		ref = "HEAD~1"
	} else {
		choice := selectOption("Reset mode", []string{
			"--soft HEAD~1  (keep changes staged)",
			"--mixed HEAD~1 (keep changes unstaged)",
			"--hard HEAD~1  (discard changes)",
			"custom",
		})
		switch {
		case strings.HasPrefix(choice, "--soft"):
			mode, ref = "--soft", "HEAD~1"
		case strings.HasPrefix(choice, "--mixed"):
			mode, ref = "--mixed", "HEAD~1"
		case strings.HasPrefix(choice, "--hard"):
			mode, ref = "--hard", "HEAD~1"
		case choice == "custom":
			mode = prompt("Mode (--soft/--mixed/--hard)")
			ref = prompt("Target ref (e.g. HEAD~1, abc1234)")
		}
	}

	if strings.Contains(mode, "--hard") {
		printWarn("This will discard working tree changes permanently.")
		if !confirm("Proceed with --hard reset?") {
			printInfo("Aborted.")
			return nil
		}
	}

	if err := RunGit(cfg, "reset", mode, ref); err != nil {
		return err
	}
	printSuccess(fmt.Sprintf("Reset %s %s", mode, ref))
	return nil
}

func cmdCheckout(cfg *Config, args []string) error {
	if len(args) == 0 {
		input := prompt("Branch or file to checkout")
		if input == "" {
			return fmt.Errorf("target required")
		}
		args = strings.Fields(input)
	}
	gitArgs := append([]string{"checkout"}, args...)
	return RunGit(cfg, gitArgs...)
}

func cmdLsFiles(cfg *Config, args []string) error {
	printHeader("Tracked files")
	out, err := RunGitBareCapture(cfg, append([]string{"ls-files"}, args...)...)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line != "" {
			fmt.Println("  " + dim("·") + " " + line)
		}
	}
	return nil
}
