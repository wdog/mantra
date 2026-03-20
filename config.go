package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	GitDir   string
	WorkTree string
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot determine home directory: %w", err)
	}

	cfg := &Config{
		GitDir:   filepath.Join(home, ".mantra.git"),
		WorkTree: home,
	}

	configPath := configFilePath(home)
	f, err := os.Open(configPath)
	if os.IsNotExist(err) {
		if werr := WriteDefaultConfig(); werr == nil {
			printInfo("Created config: " + configPath)
		}
		return cfg, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cannot open config: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = expandHome(val, home)

		switch key {
		case "git-dir":
			cfg.GitDir = val
		case "work-tree":
			cfg.WorkTree = val
		}
	}
	return cfg, scanner.Err()
}

func GitArgs(cfg *Config) []string {
	return []string{"--git-dir", cfg.GitDir, "--work-tree", cfg.WorkTree}
}

func configFilePath(home string) string {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg == "" {
		xdg = filepath.Join(home, ".config")
	}
	return filepath.Join(xdg, "mantra", "config")
}

func expandHome(path, home string) string {
	if path == "~" {
		return home
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	return path
}

func WriteDefaultConfig() error {
	home, _ := os.UserHomeDir()
	path := configFilePath(home)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	content := "# mantra configuration\n" +
		"git-dir=~/.mantra.git\n" +
		"work-tree=~\n"
	return os.WriteFile(path, []byte(content), 0644)
}
