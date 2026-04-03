package monitor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
)

func TestComputeDrift_NoChecksums(t *testing.T) {
	entries := ComputeDrift(t.TempDir(), installstate.InstallManifest{})
	if entries != nil {
		t.Errorf("expected nil for empty checksums, got %v", entries)
	}
}

func TestComputeDrift_Missing(t *testing.T) {
	root := t.TempDir()
	manifest := installstate.InstallManifest{
		FileChecksums: map[string]string{
			"file-that-was-deleted.md": hashString("original"),
		},
	}

	entries := ComputeDrift(root, manifest)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Status != DriftMissing {
		t.Errorf("status = %q, want %q", entries[0].Status, DriftMissing)
	}
	if entries[0].Path != "file-that-was-deleted.md" {
		t.Errorf("path = %q, want %q", entries[0].Path, "file-that-was-deleted.md")
	}
}

func TestComputeDrift_UserModified(t *testing.T) {
	root := t.TempDir()
	filePath := filepath.Join(root, "test.md")
	if err := os.WriteFile(filePath, []byte("modified content"), 0o644); err != nil {
		t.Fatal(err)
	}

	manifest := installstate.InstallManifest{
		FileChecksums: map[string]string{
			"test.md": hashString("original content"),
		},
	}

	entries := ComputeDrift(root, manifest)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Status != DriftUserModified {
		t.Errorf("status = %q, want %q", entries[0].Status, DriftUserModified)
	}
}

func TestComputeDrift_NoChange(t *testing.T) {
	root := t.TempDir()
	content := []byte("stable content")
	filePath := filepath.Join(root, "test.md")
	if err := os.WriteFile(filePath, content, 0o644); err != nil {
		t.Fatal(err)
	}

	manifest := installstate.InstallManifest{
		FileChecksums: map[string]string{
			"test.md": hashString("stable content"),
		},
	}

	entries := ComputeDrift(root, manifest)
	if len(entries) != 0 {
		t.Errorf("expected 0 entries (no drift), got %d: %v", len(entries), entries)
	}
}

func TestComputeDrift_IgnorePattern(t *testing.T) {
	root := t.TempDir()

	// Create drift-ignore file
	atvDir := filepath.Join(root, ".atv")
	if err := os.MkdirAll(atvDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(atvDir, "drift-ignore"), []byte("*.log\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	manifest := installstate.InstallManifest{
		FileChecksums: map[string]string{
			"debug.log": hashString("log content"),
			"readme.md": hashString("readme content"),
		},
	}

	entries := ComputeDrift(root, manifest)

	// debug.log should be ignored, readme.md should be detected as missing
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry (ignored *.log), got %d: %v", len(entries), entries)
	}
	if entries[0].Path != "readme.md" {
		t.Errorf("path = %q, want %q", entries[0].Path, "readme.md")
	}
}

func TestComputeDrift_Sorted(t *testing.T) {
	root := t.TempDir()
	manifest := installstate.InstallManifest{
		FileChecksums: map[string]string{
			"z-file.md": hashString("z"),
			"a-file.md": hashString("a"),
			"m-file.md": hashString("m"),
		},
	}

	entries := ComputeDrift(root, manifest)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Path != "a-file.md" || entries[1].Path != "m-file.md" || entries[2].Path != "z-file.md" {
		t.Errorf("entries not sorted: %v, %v, %v", entries[0].Path, entries[1].Path, entries[2].Path)
	}
}

func TestHashFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("hello world"), 0o644); err != nil {
		t.Fatal(err)
	}

	hash := hashFile(path)
	expected := hashString("hello world")
	if hash != expected {
		t.Errorf("hash = %q, want %q", hash, expected)
	}
}

func TestHashFile_Missing(t *testing.T) {
	hash := hashFile("/nonexistent/file")
	if hash != "" {
		t.Errorf("expected empty hash for missing file, got %q", hash)
	}
}

func TestHashString(t *testing.T) {
	h1 := hashString("hello")
	h2 := hashString("hello")
	h3 := hashString("world")

	if h1 != h2 {
		t.Error("same content should produce same hash")
	}
	if h1 == h3 {
		t.Error("different content should produce different hash")
	}
	if len(h1) != 64 {
		t.Errorf("hash length = %d, want 64", len(h1))
	}
}

func TestLoadDriftIgnore(t *testing.T) {
	dir := t.TempDir()
	atvDir := filepath.Join(dir, ".atv")
	if err := os.MkdirAll(atvDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content := "# Comments are ignored\n*.log\ndocs/brainstorms/*\n\n.vscode/settings.json\n"
	if err := os.WriteFile(filepath.Join(atvDir, "drift-ignore"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	patterns := loadDriftIgnore(dir)
	if len(patterns) != 3 {
		t.Fatalf("expected 3 patterns, got %d: %v", len(patterns), patterns)
	}
	if patterns[0] != "*.log" {
		t.Errorf("patterns[0] = %q, want %q", patterns[0], "*.log")
	}
}

func TestLoadDriftIgnore_Missing(t *testing.T) {
	dir := t.TempDir()
	patterns := loadDriftIgnore(dir)
	if patterns != nil {
		t.Errorf("expected nil for missing file, got %v", patterns)
	}
}

func TestIsIgnoredByDrift(t *testing.T) {
	patterns := []string{"*.log", ".vscode/settings.json"}

	tests := []struct {
		path    string
		ignored bool
	}{
		{"debug.log", true},
		{".vscode/settings.json", true},
		{"docs/plans/test.md", false},
		{".github/copilot-instructions.md", false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isIgnoredByDrift(tt.path, patterns)
			if got != tt.ignored {
				t.Errorf("isIgnoredByDrift(%q) = %v, want %v", tt.path, got, tt.ignored)
			}
		})
	}
}
