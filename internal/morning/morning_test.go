package morning

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestDiscover(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"20-later.sh", "15-middle.PS1", "10-first.sh", "30-ignored.txt", ".hidden.sh", "directory"} {
		path := filepath.Join(dir, name)
		if name == "directory" {
			if err := os.Mkdir(path, 0o755); err != nil {
				t.Fatal(err)
			}
			continue
		}
		if err := os.WriteFile(path, []byte(""), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	files, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 3 || filepath.Base(files[0]) != "10-first.sh" || filepath.Base(files[1]) != "15-middle.PS1" || filepath.Base(files[2]) != "20-later.sh" {
		t.Fatalf("Discover() = %v", files)
	}
}

func TestDiscoverFollowsSymlinksToRegularFiles(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(t.TempDir(), "target")
	if err := os.WriteFile(target, []byte(""), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "10-linked.sh")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(dir, "20-directory.sh")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "missing"), filepath.Join(dir, "30-broken.sh")); err != nil {
		t.Fatal(err)
	}

	files, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0] != link {
		t.Fatalf("Discover() = %v, want [%s]", files, link)
	}
}

func TestDiscoverMissingDirectory(t *testing.T) {
	files, err := Discover(filepath.Join(t.TempDir(), "missing"))
	if err != nil || len(files) != 0 {
		t.Fatalf("Discover() = %v, %v; want no files and no error", files, err)
	}
}

func TestRunSourcesScriptsInOrderInOneShell(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "10-first.sh")
	second := filepath.Join(dir, "20-second.sh")
	if err := os.WriteFile(first, []byte("export MORNING_TEST=value\necho first"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(second, []byte("echo \"$MORNING_TEST second\""), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if err := Run([]string{first, second}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if got, want := stdout.String(), "first\nvalue second\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestRunStopsOnScriptFailure(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "10-bad.sh")
	later := filepath.Join(dir, "20-later.sh")
	if err := os.WriteFile(bad, []byte("return 7"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(later, []byte("echo should-not-run"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if err := Run([]string{bad, later}, &stdout, &stderr); err == nil {
		t.Fatal("Run() succeeded, want script failure")
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}
