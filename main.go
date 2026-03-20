package main

import (
	"fmt"
	"os"
	"strings"
)

type handler func(cfg *Config, args []string) error

var commands = map[string]handler{
	"status":            cmdStatus,
	"diff":              cmdDiff,
	"add":               cmdAdd,
	"add -u":            cmdAddModified,
	"modified":          cmdAddModified,
	"commit":            cmdCommit,
	"push":              cmdPush,
	"pull":              cmdPull,
	"log":               cmdLog,
	"conflict":          cmdConflict,
	"stash":             cmdStash,
	"stash pop":         cmdStashPop,
	"stash list":        cmdStashList,
	"rebase":            cmdRebase,
	"rebase --continue": func(cfg *Config, _ []string) error { return RunGit(cfg, "rebase", "--continue") },
	"rebase --abort":    func(cfg *Config, _ []string) error { return RunGit(cfg, "rebase", "--abort") },
	"reset":             cmdReset,
	"checkout":          cmdCheckout,
	"ls":                cmdLsFiles,
	"files":             cmdLsFiles,
	"init":              cmdInit,
	"help":              cmdHelp,
	"completion":        cmdCompletion,
	"timeline":          cmdTimeline,
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, red("✗ Config error: ")+err.Error())
		os.Exit(1)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		runREPL(cfg)
		os.Exit(0)
	}

	// --help / -h anywhere → show help
	if args[0] == "--help" || args[0] == "-h" {
		cmdHelp(cfg, nil)
		os.Exit(0)
	}

	// Try two-word command first (e.g. "stash pop")
	subcmd := args[0]
	rest := args[1:]
	if len(args) >= 2 {
		two := args[0] + " " + args[1]
		if h, ok := commands[two]; ok {
			if err := h(cfg, args[2:]); err != nil {
				fmt.Fprintln(os.Stderr, red("✗ ")+err.Error())
				os.Exit(1)
			}
			return
		}
	}

	h, ok := commands[subcmd]
	if !ok {
		fmt.Fprintln(os.Stderr, red("✗ Unknown command: ")+subcmd)
		fmt.Fprintln(os.Stderr, dim("  Run 'mantra help' for usage."))
		os.Exit(1)
	}

	if err := h(cfg, rest); err != nil {
		msg := err.Error()
		if !strings.Contains(msg, "exit status") {
			fmt.Fprintln(os.Stderr, red("✗ ")+msg)
		}
		os.Exit(1)
	}
}
