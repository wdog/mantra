package main

import (
	"fmt"
	"strings"
)

func cmdStatus(cfg *Config, args []string) error {
	out, err := RunGitCapture(cfg, "status", "--short", "--branch")
	if err != nil && GitExitCode(err) != 1 {
		return err
	}

	printHeader("Dotfiles Status")
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "##") {
			// Branch line
			branch := strings.TrimPrefix(line, "## ")
			fmt.Println(bold("  ") + cyan(" "+branch))
			continue
		}
		if len(line) < 3 {
			fmt.Println(line)
			continue
		}
		xy := line[:2]
		file := line[3:]
		colored := colorizeStatus(xy, file)
		fmt.Println("  " + colored)
	}
	fmt.Println()
	return nil
}

func colorizeStatus(xy, file string) string {
	x := string(xy[0])
	y := string(xy[1])

	var prefix string
	switch {
	case x == "U" || y == "U" || (x == "A" && y == "A") || (x == "D" && y == "D"):
		prefix = red(xy) + " " + red(file) + dim(" (conflict)")
	case x == "A":
		prefix = green(xy) + " " + green(file)
	case x == "M" || y == "M":
		prefix = yellow(xy) + " " + yellow(file)
	case x == "D" || y == "D":
		prefix = red(xy) + " " + red(file)
	case x == "R":
		prefix = cyan(xy) + " " + cyan(file)
	case x == "?" && y == "?":
		prefix = dim(xy) + " " + dim(file)
	default:
		prefix = xy + " " + file
	}
	return prefix
}

func cmdDiff(cfg *Config, args []string) error {
	gitArgs := []string{"diff"}
	if !noColor {
		gitArgs = append(gitArgs, "--color=always")
	}
	gitArgs = append(gitArgs, args...)
	return RunGit(cfg, gitArgs...)
}

func cmdAdd(cfg *Config, args []string) error {
	if len(args) == 0 {
		input := prompt("Files to add (space-separated, or '.' for all)")
		if input == "" {
			return fmt.Errorf("no files specified")
		}
		args = strings.Fields(input)
	}
	gitArgs := append([]string{"add"}, args...)
	if err := RunGit(cfg, gitArgs...); err != nil {
		return err
	}
	printSuccess("Added: " + strings.Join(args, ", "))
	return nil
}

func cmdAddModified(cfg *Config, args []string) error {
	if err := RunGit(cfg, "add", "-u"); err != nil {
		return err
	}
	printSuccess("Staged all modified tracked files.")
	return nil
}

func cmdPush(cfg *Config, args []string) error {
	printInfo("Pushing...")
	gitArgs := append([]string{"push"}, args...)
	if err := RunGit(cfg, gitArgs...); err != nil {
		printWarn("Push failed. Try 'mantra pull' first if the remote has new commits.")
		return err
	}
	printSuccess("Pushed successfully.")
	return nil
}

func cmdPull(cfg *Config, args []string) error {
	printInfo("Pulling...")
	gitArgs := append([]string{"pull"}, args...)
	if err := RunGit(cfg, gitArgs...); err != nil {
		if gitStateFile(cfg, "MERGE_HEAD") {
			printWarn("Merge conflict detected. Run 'mantra conflict' to resolve.")
		}
		return err
	}
	printSuccess("Pulled successfully.")
	if gitStateFile(cfg, "MERGE_HEAD") {
		printWarn("Merge conflict detected. Run 'mantra conflict' to resolve.")
	}
	return nil
}

func cmdLog(cfg *Config, args []string) error {
	gitArgs := []string{"log", "--oneline", "--graph"}
	if !noColor {
		gitArgs = append(gitArgs, "--color=always")
	}
	gitArgs = append(gitArgs, args...)
	return RunGit(cfg, gitArgs...)
}
