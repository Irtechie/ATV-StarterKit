package gstack

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectPrerequisites(t *testing.T) {
	prereqs := DetectPrerequisites()

	// Git should be available in any dev environment
	if !prereqs.HasGit {
		t.Log("git not found — this is expected in some CI environments")
	} else {
		if prereqs.GitVersion == "" || prereqs.GitVersion == "unknown" {
			t.Error("git detected but version is empty or unknown")
		}
	}

	// Bun and Node are optional
	t.Logf("Prerequisites: %s", prereqs.Summary())
}

func TestPrerequisitesSummary(t *testing.T) {
	p := Prerequisites{
		HasGit:      true,
		HasBun:      true,
		HasNode:     true,
		GitVersion:  "2.44.0",
		BunVersion:  "1.1.0",
		NodeVersion: "v20.0.0",
	}

	summary := p.Summary()
	if summary == "" {
		t.Error("summary should not be empty")
	}
	if !contains(summary, "git 2.44.0") {
		t.Errorf("summary missing git version: %s", summary)
	}
	if !contains(summary, "bun 1.1.0") {
		t.Errorf("summary missing bun version: %s", summary)
	}
}

func TestPrerequisitesSummaryMissing(t *testing.T) {
	p := Prerequisites{
		HasGit: false,
		HasBun: false,
	}

	summary := p.Summary()
	if !contains(summary, "git (missing)") {
		t.Errorf("summary should show git as missing: %s", summary)
	}
	if !contains(summary, "bun (missing)") {
		t.Errorf("summary should show bun as missing: %s", summary)
	}
}

func TestRuntimeAvailable(t *testing.T) {
	tests := []struct {
		name     string
		prereqs  Prerequisites
		expected bool
	}{
		{"bun only", Prerequisites{HasBun: true}, true},
		{"node only", Prerequisites{HasNode: true}, true},
		{"both", Prerequisites{HasBun: true, HasNode: true}, true},
		{"neither", Prerequisites{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.prereqs.RuntimeAvailable(); got != tt.expected {
				t.Errorf("RuntimeAvailable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStripGit(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a fake .git directory
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatal(err)
	}
	// Write a file inside .git
	if err := os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main"), 0644); err != nil {
		t.Fatal(err)
	}

	// Verify .git exists
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Fatal(".git should exist before strip")
	}

	// Strip it
	if err := StripGit(tmpDir); err != nil {
		t.Fatalf("StripGit failed: %v", err)
	}

	// Verify .git is gone
	if _, err := os.Stat(gitDir); !os.IsNotExist(err) {
		t.Error(".git should not exist after strip")
	}
}

func TestListSkillDirs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create fake skill directories
	for _, name := range []string{"review", "qa", "ship"} {
		skillDir := filepath.Join(tmpDir, name)
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# "+name), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create a non-skill directory (no SKILL.md)
	if err := os.MkdirAll(filepath.Join(tmpDir, "lib"), 0755); err != nil {
		t.Fatal(err)
	}

	dirs := listSkillDirs(tmpDir)
	if len(dirs) != 3 {
		t.Errorf("expected 3 skill dirs, got %d: %v", len(dirs), dirs)
	}
}

func TestInstallIdempotent(t *testing.T) {
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "gstack")

	// Create a fake existing gstack install
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(targetDir, "SKILL.md"), []byte("# gstack"), 0644); err != nil {
		t.Fatal(err)
	}

	// Install should skip when target already exists with SKILL.md
	result := Install(targetDir, ModeMarkdownOnly)
	if result.Error != nil {
		t.Fatalf("idempotent install should not error: %v", result.Error)
	}
	if result.Warning == "" {
		t.Error("idempotent install should produce a warning about existing install")
	}
	if result.Cloned {
		t.Error("should not clone when already installed")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
