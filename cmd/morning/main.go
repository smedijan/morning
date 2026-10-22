package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/smedijan/morning/internal/morning"
)

func main() {
	if err := run(os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "morning:", err)
		os.Exit(1)
	}
}

func run(stdout, stderr io.Writer) error {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("find home directory: %w", err)
		}
		configHome = filepath.Join(home, ".config")
	}

	scripts, err := morning.Discover(filepath.Join(configHome, "morning.d"))
	if err != nil {
		return err
	}
	return morning.Run(scripts, stdout, stderr)
}
