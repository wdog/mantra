<div align="center">

![Mantra Banner](assets/banner.png)

# mantra

**Dotfiles manager built on a bare git repository.**
No external dependencies. Single static binary. Works anywhere git does.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-a855f7?style=flat)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux-6366f1?style=flat&logo=linux&logoColor=white)](https://github.com)

</div>

---

## 📑 Index

- [💡 How it works](#-how-it-works)
  - [Typical usage](#typical-usage)
- [📦 Requirements](#-requirements)
- [🔨 Build](#-build)
- [🚀 Installation](#-installation)
- [⚙️ Configuration](#️-configuration)
- [🛠️ Initial Setup](#️-initial-setup)
  - [New repository](#new-repository)
  - [Existing repository](#existing-repository)
  - [Suppress untracked file noise](#suppress-untracked-file-noise)
- [💾 Saving dotfiles](#-saving-dotfiles)
- [💻 Usage](#-usage)
  - [✨ Interactive Mode (REPL)](#-interactive-mode-repl)
  - [📋 Basic Commands](#-basic-commands)
  - [🐚 Shell Completion](#-shell-completion)
  - [🔀 Conflict & History Management](#-conflict--history-management)
  - [⚡ Conflict Resolution Flow](#-conflict-resolution-flow)
  - [🔄 Interactive Reset](#-interactive-reset)

---

## 💡 How it works

mantra uses a **bare git repository** — a repo with no working tree of its own. Instead, it points git's working tree directly at your home directory (`~`). This means:

- Your `~` **is** the working tree. Every file in your home directory is a potential dotfile.
- The git history lives in `~/.mantra.git` (or wherever `git-dir` points), completely separate from your files.
- No symlinks, no copying, no separate dotfiles folder. Files stay exactly where they are.

```
~/.mantra.git/   ← bare repository (git objects, history, config)
~/               ← working tree (your actual home directory)
```

Under the hood every `mantra` command translates to a plain `git` command with `--git-dir` and `--work-tree` flags injected automatically:

```bash
mantra status   →   git --git-dir=~/.mantra.git --work-tree=~ status
mantra add      →   git --git-dir=~/.mantra.git --work-tree=~ add
mantra push     →   git --git-dir=~/.mantra.git --work-tree=~ push
```

Because it is plain git under the hood, you can always fall back to raw git commands using those flags and everything will work as expected.

### Typical usage

**First machine — start tracking your dotfiles:**

```bash
mantra init                              # create ~/.mantra.git
git --git-dir=~/.mantra.git remote add origin git@github.com:you/dotfiles.git

mantra add -f ~/.bashrc                  # start tracking a file
mantra add -f ~/.config/nvim/init.lua
mantra commit -m "initial dotfiles"
mantra push
```

**Day to day — edit a file, save it:**

```bash
# edit ~/.bashrc as usual, then:
mantra modified                          # stage all modified tracked files
mantra commit -m "update bashrc aliases"
mantra push
```

**Second machine — restore your dotfiles:**

```bash
git clone --bare git@github.com:you/dotfiles.git ~/.mantra.git
mantra pull                              # apply files to ~
```

**Check what is tracked and what has changed:**

```bash
mantra ls                                # list all tracked files
mantra status                            # what changed since last commit
mantra diff                              # show the actual diff
mantra log                               # commit history
```

---

## 📦 Requirements

- **Go** 1.21+ *(only for building from source)*
- **Git**

---

## 🔨 Build

**From source:**

```bash
go build -ldflags="-s -w" -o mantra .
```

**Linux amd64 (pre-built binary):**

Download the latest `mantra-linux-amd64` from the [Releases](../../releases) page.

**With Task** — requires [Task](https://taskfile.dev):

```bash
task build      # build for current platform
task install    # build + copy to ~/.local/bin
task clean      # remove built binaries
```

---

## 🚀 Installation

```bash
cp mantra ~/.local/bin/mantra
chmod +x ~/.local/bin/mantra
```

Make sure `~/.local/bin` is in your `$PATH`.

---

## ⚙️ Configuration

Config is read from `$XDG_CONFIG_HOME/mantra/config` (default: `~/.config/mantra/config`).
Created automatically with defaults on first run.

```ini
# ~/.config/mantra/config
git-dir=~/.mantra.git
work-tree=~
```

| Key | Default | Description |
|---|---|---|
| `git-dir` | `~/.mantra.git` | Path to the bare repository |
| `work-tree` | `~` | Root directory of dotfiles |

---

## 🛠️ Initial Setup

### New repository

```bash
mantra init
```

Then add a remote so you can push your dotfiles to GitHub/GitLab/etc:

```bash
git --git-dir=~/.mantra.git remote add origin <url>
```

### Existing repository

Clone an existing dotfiles repo as a bare repository, then use normally:

```bash
git clone --bare <url> ~/.mantra.git
```

### Suppress untracked file noise

With `work-tree=$HOME`, git sees every file in your home directory as untracked. The recommended fix is to ignore everything by default and only track files explicitly:

```gitignore
# ~/.gitignore
*
```

Then add only the files you want to track:

```bash
mantra add -f ~/.bashrc
mantra add -f ~/.config/nvim/init.lua
```

---

## 💾 Saving dotfiles

The typical workflow for tracking and syncing a dotfile:

**1. Track a new file:**

```bash
mantra add -f ~/.bashrc
mantra commit -m "track bashrc"
mantra push
```

**2. Stage all modified tracked files at once:**

```bash
mantra modified        # or: mantra add -u
mantra commit -m "update dotfiles"
mantra push
```

**3. Pull changes on another machine:**

```bash
mantra pull
```

**4. See what has changed:**

```bash
mantra status          # overview
mantra diff            # full diff
mantra log             # commit history
```

**5. List all tracked dotfiles:**

```bash
mantra ls
```

---

## 💻 Usage

```bash
mantra [command] [args]
```

Running **without arguments** launches the interactive REPL.

### ✨ Interactive Mode (REPL)

```bash
mantra
```

Opens an interactive shell with a live-updating prompt:

```
mantra  main ↑1 +2 ~1 ›
```

**Tab completion** is supported for all commands, subcommands, and file paths:

- `add ~/.<TAB>` — completes file paths inside the work-tree
- `add -f ~/.<TAB>` — same, for force-adding ignored files
- `diff <TAB>` — completes file paths
- `stash <TAB>` — completes `list` / `pop`
- `rebase <TAB>` — completes `--abort` / `--continue`
- `completion <TAB>` — completes `fish` / `bash` / `zsh`

Exit with `q`, `exit`, `quit`, `ESC`, or `Ctrl-D`.

#### Prompt indicators

| Symbol | Meaning |
|:---:|---|
| `↑N` | N commits ahead of remote |
| `↓N` | N commits behind remote |
| `!N` | N conflicted files |
| `+N` | N staged files |
| `~N` | N modified files |
| `?N` | N untracked files |

---

### 📋 Basic Commands

| Command | Description |
|---|---|
| `status` | Show working tree status |
| `diff [file]` | Show changes |
| `add [files]` | Stage files (interactive prompt if omitted) |
| `add -u` / `modified` | Stage all modified tracked files |
| `commit [-m msg]` | Commit staged changes (interactive prompt if `-m` omitted) |
| `push` | Push to remote |
| `pull` | Pull from remote |
| `log` | Show commit log (oneline graph) |
| `ls` / `files` | List all tracked files |
| `init` | Initialize the bare mantra repository |
| `completion [shell]` | Print shell completion script (`fish`, `bash`, `zsh`) |
| `help` / `-h` / `--help` | Show help |

---

### 🐚 Shell Completion

Auto-detects your shell from `$SHELL`, or pass it explicitly:

```bash
mantra completion         # auto-detect
mantra completion fish
mantra completion bash
mantra completion zsh
```

**Fish**
```fish
mantra completion fish > ~/.config/fish/completions/mantra.fish
```

**Bash**
```bash
# one-off
source <(mantra completion bash)
# permanent
mantra completion bash > /etc/bash_completion.d/mantra
```

**Zsh**
```zsh
mantra completion zsh > ~/.zsh/completions/_mantra
# add to ~/.zshrc if not already present:
# fpath=(~/.zsh/completions $fpath) && autoload -Uz compinit && compinit
```

---

### 🔀 Conflict & History Management

| Command | Description |
|---|---|
| `conflict` | Guided conflict resolution flow |
| `stash [msg]` | Stash working changes (interactive prompt if omitted) |
| `stash pop` | Apply last stash |
| `stash list` | List stashes |
| `rebase <branch>` | Rebase onto branch (interactive prompt if omitted) |
| `rebase --continue` | Continue rebase after resolving conflicts |
| `rebase --abort` | Abort rebase |
| `reset` | Reset HEAD with interactive mode selector |
| `checkout <ref>` | Checkout branch or file (interactive prompt if omitted) |

---

### ⚡ Conflict Resolution Flow

```bash
mantra conflict
```

For each conflicted file the diff is shown and an action is requested:

| Option | Behavior |
|---|---|
| `edit` | Open `$EDITOR` to resolve manually |
| `ours` | Keep the local version |
| `theirs` | Keep the remote version |
| `skip` | Skip the file (leaves it conflicted) |

Once all conflicts are resolved, offers to commit immediately.

---

### 🔄 Interactive Reset

```bash
mantra reset
```

Without arguments, presents a selection menu:

| Option | Behavior |
|---|---|
| `--soft HEAD~1` | Keep changes staged |
| `--mixed HEAD~1` | Keep changes unstaged |
| `--hard HEAD~1` | Discard changes *(asks for confirmation)* |
| `custom` | Manually enter mode and ref |
