package main

import (
	"fmt"
	"os"
)

func cmdInit(cfg *Config, args []string) error {
	printHeader("Initialize mantra repo")
	printInfo("Git dir:   " + cfg.GitDir)
	printInfo("Work tree: " + cfg.WorkTree)

	if _, err := os.Stat(cfg.GitDir); err == nil {
		printWarn("Git dir already exists.")
		return nil
	}

	if !confirm("Initialize bare repo at " + cfg.GitDir + "?") {
		printInfo("Aborted.")
		return nil
	}

	if err := RunGit(cfg, "init", "--bare", cfg.GitDir); err != nil {
		return err
	}
	printSuccess("Initialized.")
	printInfo("Tip: Add a remote with 'git --git-dir=" + cfg.GitDir + " remote add origin <url>'")
	printInfo("Tip: Add '*' to ~/.gitignore then use 'mantra add -f <file>' to track specific files.")
	return nil
}

func cmdHelp(cfg *Config, args []string) error {
	fmt.Println()
	fmt.Println(bold(cyan("mantra")) + dim(" — dotfiles manager"))
	fmt.Println()
	fmt.Println(bold("Usage:") + "  mantra <command> [args]")
	fmt.Println()
	fmt.Println(bold("Commands:"))

	cmds := [][]string{
		{"status", "Show working tree status"},
		{"diff [file]", "Show changes"},
		{"add [files]", "Stage files (interactive if omitted)"},
		{"commit [-m msg]", "Commit staged changes"},
		{"push", "Push to remote"},
		{"pull", "Pull from remote"},
		{"log", "Show commit log"},
		{"ls / files", "List all tracked files"},
		{"", ""},
		{"conflict", "Guided conflict resolution"},
		{"stash [msg]", "Stash working changes"},
		{"stash pop", "Apply last stash"},
		{"stash list", "List stashes"},
		{"rebase <branch>", "Rebase onto branch"},
		{"rebase --continue", "Continue rebase after conflict"},
		{"rebase --abort", "Abort rebase"},
		{"reset", "Reset HEAD (interactive mode selector)"},
		{"checkout <ref>", "Checkout branch or file"},
		{"", ""},
		{"init", "Initialize bare mantra repo"},
		{"help", "Show this help"},
	}

	for _, row := range cmds {
		if row[0] == "" {
			fmt.Println()
			continue
		}
		fmt.Printf("  %-26s %s\n", cyan(row[0]), dim(row[1]))
	}

	fmt.Println()
	fmt.Println(dim("Config: ~/.config/mantra/config"))
	fmt.Println(dim("  git-dir=~/.mantra.git"))
	fmt.Println(dim("  work-tree=~"))
	fmt.Println()
	return nil
}
