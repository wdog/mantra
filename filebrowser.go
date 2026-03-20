package main

// FileBrowser handles interactive file selection from the tracked dotfiles.
// Keeping it separate makes it easy to swap the UI (e.g. fzf, TUI) later.

import (
	"strings"
)

// BrowseTrackedFiles shows an interactive list of all tracked files and
// returns the selected path relative to the work-tree, or "" if cancelled.
func BrowseTrackedFiles(cfg *Config, label string) (string, error) {
	out, err := RunGitBareCapture(cfg, "ls-files")
	if err != nil {
		return "", err
	}
	files := strings.Split(strings.TrimRight(out, "\n"), "\n")
	// filter empty lines
	var clean []string
	for _, f := range files {
		if f != "" {
			clean = append(clean, f)
		}
	}
	if len(clean) == 0 {
		return "", nil
	}
	return selectOption(label, clean), nil
}
