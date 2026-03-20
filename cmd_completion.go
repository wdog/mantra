package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func cmdCompletion(_ *Config, args []string) error {
	shell := ""
	if len(args) > 0 {
		shell = args[0]
	} else {
		// Auto-detect from $SHELL
		shell = filepath.Base(os.Getenv("SHELL"))
	}

	switch shell {
	case "fish":
		fmt.Print(fishCompletion)
		printInstallHint("fish", "~/.config/fish/completions/mantra.fish")
	case "bash":
		fmt.Print(bashCompletion)
		printInstallHint("bash", "/etc/bash_completion.d/mantra  # or source in ~/.bashrc")
	case "zsh":
		fmt.Print(zshCompletion)
		printInstallHint("zsh", "~/.zsh/completions/_mantra  # ensure $fpath includes that dir")
	default:
		shells := strings.Join([]string{"fish", "bash", "zsh"}, ", ")
		fmt.Fprintln(os.Stderr, red("✗ Unknown shell: ")+shell)
		fmt.Fprintln(os.Stderr, dim("  Supported: "+shells))
		fmt.Fprintln(os.Stderr, dim("  Usage: mantra completion <shell>"))
		return fmt.Errorf("unsupported shell: %s", shell)
	}
	return nil
}

func printInstallHint(shell, path string) {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, dim("# Install "+shell+" completions:"))
	fmt.Fprintln(os.Stderr, dim("#   mantra completion "+shell+" > "+path))
}

const fishCompletion = `# mantra fish completions
# Install: mantra completion fish > ~/.config/fish/completions/mantra.fish

complete -c mantra -e
complete -c mantra -f

set -l cmds add checkout commit conflict diff files help init log ls pull push rebase reset stash status

# Top-level commands
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'status'    -d 'Show working tree status'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'diff'      -d 'Show changes'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'add'       -d 'Stage files'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'commit'    -d 'Commit staged changes'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'push'      -d 'Push to remote'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'pull'      -d 'Pull from remote'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'log'       -d 'Show commit log'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'ls'        -d 'List tracked files'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'files'     -d 'List tracked files'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'conflict'  -d 'Guided conflict resolution'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'stash'     -d 'Stash working changes'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'rebase'    -d 'Rebase onto branch'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'reset'     -d 'Reset HEAD interactively'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'checkout'  -d 'Checkout branch or file'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'init'      -d 'Initialize bare repo'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'help'      -d 'Show help'
complete -c mantra -n "not __fish_seen_subcommand_from $cmds" -a 'completion' -d 'Print shell completion script'

# add / diff / checkout — re-enable path completion
complete -c mantra -n '__fish_seen_subcommand_from add diff checkout' -F
complete -c mantra -n '__fish_seen_subcommand_from add' -s f -d 'Force add ignored file'
complete -c mantra -n '__fish_seen_subcommand_from commit' -s m -d 'Commit message'

# stash sub-commands
complete -c mantra -n '__fish_seen_subcommand_from stash' -a 'pop'  -d 'Apply last stash'
complete -c mantra -n '__fish_seen_subcommand_from stash' -a 'list' -d 'List stashes'

# rebase sub-commands
complete -c mantra -n '__fish_seen_subcommand_from rebase' -a '--continue' -d 'Continue after conflict'
complete -c mantra -n '__fish_seen_subcommand_from rebase' -a '--abort'    -d 'Abort rebase'

# completion sub-commands
complete -c mantra -n '__fish_seen_subcommand_from completion' -a 'fish bash zsh'
`

const bashCompletion = `# mantra bash completions
# Install: mantra completion bash > /etc/bash_completion.d/mantra
#      or: source <(mantra completion bash)   # add to ~/.bashrc

_mantra_complete() {
    local cur prev
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    local cmds="status diff add commit push pull log ls files conflict stash rebase reset checkout init help completion"

    case "$prev" in
        mantra)
            COMPREPLY=( $(compgen -W "$cmds" -- "$cur") )
            return ;;
        add|diff|checkout)
            COMPREPLY=( $(compgen -f -- "$cur") )
            return ;;
        stash)
            COMPREPLY=( $(compgen -W "pop list" -- "$cur") )
            return ;;
        rebase)
            COMPREPLY=( $(compgen -W "--continue --abort" -- "$cur") )
            return ;;
        completion)
            COMPREPLY=( $(compgen -W "fish bash zsh" -- "$cur") )
            return ;;
    esac
}

complete -F _mantra_complete mantra
`

const zshCompletion = `#compdef mantra
# mantra zsh completions
# Install: mantra completion zsh > ~/.zsh/completions/_mantra
#   then add to ~/.zshrc:  fpath=(~/.zsh/completions $fpath) && autoload -Uz compinit && compinit

_mantra() {
    local state

    _arguments \
        '1: :->cmd' \
        '*: :->args'

    case $state in
        cmd)
            local cmds=(
                'status:Show working tree status'
                'diff:Show changes'
                'add:Stage files'
                'commit:Commit staged changes'
                'push:Push to remote'
                'pull:Pull from remote'
                'log:Show commit log'
                'ls:List tracked files'
                'files:List tracked files'
                'conflict:Guided conflict resolution'
                'stash:Stash working changes'
                'rebase:Rebase onto branch'
                'reset:Reset HEAD interactively'
                'checkout:Checkout branch or file'
                'init:Initialize bare repo'
                'help:Show help'
                'completion:Print shell completion script'
            )
            _describe 'command' cmds
            ;;
        args)
            case ${words[2]} in
                add|diff|checkout)
                    _files
                    ;;
                stash)
                    local subs=('pop:Apply last stash' 'list:List stashes')
                    _describe 'subcommand' subs
                    ;;
                rebase)
                    local subs=('--continue:Continue after conflict' '--abort:Abort rebase')
                    _describe 'subcommand' subs
                    ;;
                completion)
                    local shells=('fish' 'bash' 'zsh')
                    _describe 'shell' shells
                    ;;
            esac
            ;;
    esac
}

_mantra
`
