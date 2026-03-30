package sandbox

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/detect"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/gstack"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/scaffold"
)

func TestAutoModeNoGstack(t *testing.T) {
	// Auto mode should install ATV files only, no gstack
	sandboxDir := t.TempDir()

	// Create a fake TypeScript project
	if err := os.WriteFile(filepath.Join(sandboxDir, "tsconfig.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	env := detect.DetectEnvironment(sandboxDir)
	if env.Stack != detect.StackTypeScript {
		t.Fatalf("expected TypeScript stack, got %s", env.Stack)
	}

	catalog := scaffold.BuildCatalog(env.Stack)
	results := scaffold.WriteAll(sandboxDir, catalog)

	if len(results) == 0 {
		t.Fatal("auto mode should produce results")
	}

	// ATV files should exist
	assertFileExists(t, filepath.Join(sandboxDir, ".github", "copilot-instructions.md"))
	assertFileExists(t, filepath.Join(sandboxDir, ".github", "copilot-setup-steps.yml"))

	// gstack should NOT exist
	gstackDir := filepath.Join(sandboxDir, ".github", "skills", "gstack")
	if _, err := os.Stat(gstackDir); !os.IsNotExist(err) {
		t.Error("auto mode should NOT install gstack")
	}
}

func TestATVFilesCreatedCorrectly(t *testing.T) {
	sandboxDir := t.TempDir()

	// Create a fake Python project
	if err := os.WriteFile(filepath.Join(sandboxDir, "pyproject.toml"), []byte("[project]"), 0644); err != nil {
		t.Fatal(err)
	}

	env := detect.DetectEnvironment(sandboxDir)
	catalog := scaffold.BuildCatalog(env.Stack)
	results := scaffold.WriteAll(sandboxDir, catalog)

	// Count created files
	created := 0
	for _, r := range results {
		if r.Status == scaffold.StatusCreated || r.Status == scaffold.StatusDirCreated {
			created++
		}
	}

	if created == 0 {
		t.Error("should have created files")
	}

	// Check key ATV files
	assertFileExists(t, filepath.Join(sandboxDir, ".github", "copilot-instructions.md"))
	assertFileExists(t, filepath.Join(sandboxDir, ".github", "copilot-setup-steps.yml"))
	assertFileExists(t, filepath.Join(sandboxDir, ".github", "copilot-mcp-config.json"))
	assertDirExists(t, filepath.Join(sandboxDir, ".github", "skills"))
	assertDirExists(t, filepath.Join(sandboxDir, ".github", "agents"))
	assertDirExists(t, filepath.Join(sandboxDir, "docs", "plans"))
	assertDirExists(t, filepath.Join(sandboxDir, "docs", "brainstorms"))
}

func TestGstackInstallMarkdownOnly(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network-dependent test in short mode")
	}

	sandboxDir := t.TempDir()
	gstackDir := filepath.Join(sandboxDir, ".github", "skills", "gstack")

	result := gstack.Install(gstackDir, gstack.ModeMarkdownOnly)

	if result.Error != nil {
		t.Fatalf("gstack install failed: %v", result.Error)
	}

	if !result.Cloned {
		t.Error("expected gstack to be cloned")
	}

	// Should not build in markdown-only mode
	if result.Built {
		t.Error("markdown-only mode should not build")
	}

	// .git should be stripped
	if _, err := os.Stat(filepath.Join(gstackDir, ".git")); !os.IsNotExist(err) {
		t.Error(".git directory should be removed after install")
	}

	// Should have skill dirs
	if len(result.SkillDirs) == 0 {
		t.Error("expected at least some skill directories")
	}
}

func TestGstackIdempotent(t *testing.T) {
	sandboxDir := t.TempDir()
	gstackDir := filepath.Join(sandboxDir, "gstack")

	// Create a fake existing install
	if err := os.MkdirAll(gstackDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gstackDir, "SKILL.md"), []byte("# gstack"), 0644); err != nil {
		t.Fatal(err)
	}

	result := gstack.Install(gstackDir, gstack.ModeMarkdownOnly)
	if result.Error != nil {
		t.Fatalf("idempotent install should not error: %v", result.Error)
	}
	if result.Cloned {
		t.Error("should not clone when already installed")
	}
}

func TestInstructionsContainGstackSection(t *testing.T) {
	sandboxDir := t.TempDir()

	env := detect.DetectEnvironment(sandboxDir) // general stack
	catalog := scaffold.BuildCatalog(env.Stack)
	scaffold.WriteAll(sandboxDir, catalog)

	// Read copilot-instructions.md
	content, err := os.ReadFile(filepath.Join(sandboxDir, ".github", "copilot-instructions.md"))
	if err != nil {
		t.Fatalf("failed to read copilot-instructions.md: %v", err)
	}

	contentStr := string(content)

	// Should contain gstack section
	if !containsStr(contentStr, "gstack Skills") {
		t.Error("copilot-instructions.md should contain gstack Skills section")
	}

	// Should contain ATV override rules
	if !containsStr(contentStr, "ATV Override Rules") {
		t.Error("copilot-instructions.md should contain ATV Override Rules")
	}

	// Should contain protected artifacts
	if !containsStr(contentStr, "Protected artifacts") {
		t.Error("copilot-instructions.md should contain protected artifacts rule")
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", path)
	}
}

func assertDirExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Errorf("expected directory to exist: %s", path)
		return
	}
	if !info.IsDir() {
		t.Errorf("expected %s to be a directory", path)
	}
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
