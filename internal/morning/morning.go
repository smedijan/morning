package morning

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Discover returns non-hidden .sh and .ps1 files directly inside dir in filename order.
func Discover(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read morning scripts directory: %w", err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		extension := strings.ToLower(filepath.Ext(entry.Name()))
		if entry.Name()[0] == '.' || (extension != ".sh" && extension != ".ps1") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		if entry.Type()&os.ModeSymlink != 0 {
			info, err := os.Stat(path)
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("stat morning script %q: %w", path, err)
			}
			if !info.Mode().IsRegular() {
				continue
			}
		} else if !entry.Type().IsRegular() {
			continue
		}
		files = append(files, path)
	}
	sort.Strings(files)
	return files, nil
}

// Run sources scripts in filename order, stopping if one fails. Consecutive scripts
// with the same extension share a shell process and its state.
func Run(scripts []string, stdout, stderr io.Writer) error {
	for i := 0; i < len(scripts); {
		extension := strings.ToLower(filepath.Ext(scripts[i]))
		j := i + 1
		for j < len(scripts) && strings.ToLower(filepath.Ext(scripts[j])) == extension {
			j++
		}

		var cmd *exec.Cmd
		switch extension {
		case ".sh":
			args := append([]string{"-c", `for script do . "$script" || exit $?; done`, "morning"}, scripts[i:j]...)
			cmd = exec.Command("sh", args...) // #nosec G204 -- morning scripts are intentionally sourced as user-provided code.
		case ".ps1":
			args := append([]string{"-NoProfile", "-NonInteractive", "-Command", `$ErrorActionPreference = 'Stop'; foreach ($script in $args) { . $script; if (-not $?) { exit 1 } }`}, scripts[i:j]...)
			if _, err := exec.LookPath("pwsh"); err == nil {
				cmd = exec.Command("pwsh", args...) // #nosec G204 -- morning scripts are intentionally sourced as user-provided code.
			} else if _, err := exec.LookPath("powershell"); err == nil {
				cmd = exec.Command("powershell", args...) // #nosec G204 -- morning scripts are intentionally sourced as user-provided code.
			} else {
				return fmt.Errorf("run morning scripts: PowerShell not found (install pwsh or powershell)")
			}
		default:
			return fmt.Errorf("unsupported morning script extension %q", extension)
		}
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("run morning scripts: %w", err)
		}
		i = j
	}
	return nil
}
