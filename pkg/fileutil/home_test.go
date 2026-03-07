package fileutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func setTestHome(t *testing.T, home string) {
	t.Helper()
	// Clear Windows-specific env vars so HOME is the single source of truth in tests.
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", "")
	t.Setenv("HOMEDRIVE", "")
	t.Setenv("HOMEPATH", "")
}

func TestHomeDirPrefersHOME(t *testing.T) {
	setTestHome(t, filepath.Join(string(os.PathSeparator), "tmp", "picoclaw-home"))

	got, err := HomeDir()
	if err != nil {
		t.Fatalf("HomeDir() error = %v", err)
	}
	if got != filepath.Join(string(os.PathSeparator), "tmp", "picoclaw-home") {
		t.Fatalf("HomeDir() = %q, want %q", got, filepath.Join(string(os.PathSeparator), "tmp", "picoclaw-home"))
	}
}

func TestHomeDirWindowsFallbackToDrivePath(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only")
	}

	setTestHome(t, "")
	t.Setenv("HOMEDRIVE", "C:")
	t.Setenv("HOMEPATH", `\Users\Test`)

	got, err := HomeDir()
	if err != nil {
		t.Fatalf("HomeDir() error = %v", err)
	}
	if got != `C:\Users\Test` {
		t.Fatalf("HomeDir() = %q, want %q", got, `C:\Users\Test`)
	}
}

func TestExpandHome(t *testing.T) {
	home := filepath.Join(string(os.PathSeparator), "tmp", "picoclaw-home")
	setTestHome(t, home)

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "tilde", in: "~", want: home},
		{name: "slash path", in: "~/workspace", want: filepath.Join(home, "workspace")},
		{name: "other user unchanged", in: "~alice/workspace", want: "~alice/workspace"},
		{name: "plain unchanged", in: "/var/tmp", want: "/var/tmp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpandHome(tt.in)
			if got != tt.want {
				t.Fatalf("ExpandHome(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestExpandHomeWindowsSeparator(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only")
	}

	setTestHome(t, `C:\Users\Test`)
	got := ExpandHome(`~\workspace\skills`)
	want := filepath.Join(`C:\Users\Test`, "workspace", "skills")
	if got != want {
		t.Fatalf("ExpandHome() = %q, want %q", got, want)
	}
}
