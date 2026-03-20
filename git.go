package main

import (
	"bytes"
	"os"
	"os/exec"
)

func RunGit(cfg *Config, subcmd ...string) error {
	args := append(GitArgs(cfg), subcmd...)
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunGitBareCapture(cfg *Config, subcmd ...string) (string, error) {
	args := append([]string{"--git-dir", cfg.GitDir}, subcmd...)
	cmd := exec.Command("git", args...)
	cmd.Stderr = nil // suppress bare+worktree config warnings
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err := cmd.Run()
	return buf.String(), err
}

func RunGitCapture(cfg *Config, subcmd ...string) (string, error) {
	args := append(GitArgs(cfg), subcmd...)
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err := cmd.Run()
	return buf.String(), err
}

func GitExitCode(err error) int {
	if err == nil {
		return 0
	}
	if ex, ok := err.(*exec.ExitError); ok {
		return ex.ExitCode()
	}
	return -1
}

func gitStateFile(cfg *Config, name string) bool {
	_, err := os.Stat(cfg.GitDir + "/" + name)
	return err == nil
}
