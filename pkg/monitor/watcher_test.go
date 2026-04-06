package monitor

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewWatcher(t *testing.T) {
	root := t.TempDir()
	w, err := NewWatcher(root, WatcherOptions{})
	if err != nil {
		t.Fatalf("NewWatcher: %v", err)
	}
	defer w.fsWatcher.Close()

	if w.debounceWin != 300*time.Millisecond {
		t.Errorf("debounce = %v, want 300ms", w.debounceWin)
	}
}

func TestNewWatcher_CustomDebounce(t *testing.T) {
	root := t.TempDir()
	w, err := NewWatcher(root, WatcherOptions{DebounceWindow: 100 * time.Millisecond})
	if err != nil {
		t.Fatalf("NewWatcher: %v", err)
	}
	defer w.fsWatcher.Close()

	if w.debounceWin != 100*time.Millisecond {
		t.Errorf("debounce = %v, want 100ms", w.debounceWin)
	}
}

func TestWatcher_StartStop(t *testing.T) {
	root := t.TempDir()
	w, err := NewWatcher(root, WatcherOptions{})
	if err != nil {
		t.Fatalf("NewWatcher: %v", err)
	}

	ctx := context.Background()
	if err := w.Start(ctx); err != nil {
		t.Fatalf("Start: %v", err)
	}

	// Should have initial state
	state := w.State()
	if state.SchemaVersion != stateSchemaVersion {
		t.Errorf("SchemaVersion = %d, want %d", state.SchemaVersion, stateSchemaVersion)
	}

	if err := w.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
}

func TestWatcher_GracefulWithMissingDirs(t *testing.T) {
	root := t.TempDir()
	// Don't create any subdirectories — watcher should start fine
	w, err := NewWatcher(root, WatcherOptions{})
	if err != nil {
		t.Fatalf("NewWatcher: %v", err)
	}

	ctx := context.Background()
	if err := w.Start(ctx); err != nil {
		t.Fatalf("Start with missing dirs: %v", err)
	}

	state := w.State()
	if len(state.Brainstorms) != 0 {
		t.Errorf("expected 0 brainstorms, got %d", len(state.Brainstorms))
	}

	if err := w.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
}

func TestWatcher_FullScan(t *testing.T) {
	root := t.TempDir()

	// Create brainstorms dir with a file
	brainstormDir := filepath.Join(root, "docs", "brainstorms")
	if err := os.MkdirAll(brainstormDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(brainstormDir, "test-brainstorm.md"), []byte("# Test"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create plans dir with a file
	plansDir := filepath.Join(root, "docs", "plans")
	if err := os.MkdirAll(plansDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(plansDir, "test-plan.md"), []byte("# Plan"), 0o644); err != nil {
		t.Fatal(err)
	}

	w, err := NewWatcher(root, WatcherOptions{})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := w.Start(ctx); err != nil {
		t.Fatal(err)
	}
	defer w.Stop() //nolint:errcheck

	state := w.State()
	if len(state.Brainstorms) != 1 {
		t.Errorf("brainstorms = %d, want 1", len(state.Brainstorms))
	}
	if len(state.Plans) != 1 {
		t.Errorf("plans = %d, want 1", len(state.Plans))
	}
	if state.Brainstorms[0].Name != "test-brainstorm.md" {
		t.Errorf("brainstorm name = %q, want %q", state.Brainstorms[0].Name, "test-brainstorm.md")
	}
}

func TestWatcher_ForceRefresh(t *testing.T) {
	root := t.TempDir()

	w, err := NewWatcher(root, WatcherOptions{})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := w.Start(ctx); err != nil {
		t.Fatal(err)
	}
	defer w.Stop() //nolint:errcheck

	// Initially no brainstorms
	state := w.State()
	if len(state.Brainstorms) != 0 {
		t.Fatalf("expected 0 brainstorms initially, got %d", len(state.Brainstorms))
	}

	// Create a brainstorm dir and file
	brainstormDir := filepath.Join(root, "docs", "brainstorms")
	if err := os.MkdirAll(brainstormDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(brainstormDir, "new.md"), []byte("# New"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Force refresh should pick it up
	w.ForceRefresh()

	state = w.State()
	if len(state.Brainstorms) != 1 {
		t.Errorf("after refresh: brainstorms = %d, want 1", len(state.Brainstorms))
	}
}

func TestWatcher_StateFile(t *testing.T) {
	root := t.TempDir()

	w, err := NewWatcher(root, WatcherOptions{})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := w.Start(ctx); err != nil {
		t.Fatal(err)
	}
	defer w.Stop() //nolint:errcheck

	// State file should be written
	stateFile := filepath.Join(root, ".atv", "dashboard-state.json")
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		t.Error("state file not written")
	}

	// Should be valid JSON
	data, err := os.ReadFile(stateFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("state file is empty")
	}
}

func TestWatcher_ContextEstimate(t *testing.T) {
	root := t.TempDir()

	// Create .github with copilot-instructions.md
	githubDir := filepath.Join(root, ".github")
	if err := os.MkdirAll(githubDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "This is copilot instructions\n"
	if err := os.WriteFile(filepath.Join(githubDir, "copilot-instructions.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create skills directory
	skillsDir := filepath.Join(githubDir, "skills")
	if err := os.MkdirAll(filepath.Join(skillsDir, "skill-a"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(skillsDir, "skill-b"), 0o755); err != nil {
		t.Fatal(err)
	}

	w, err := NewWatcher(root, WatcherOptions{})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := w.Start(ctx); err != nil {
		t.Fatal(err)
	}
	defer w.Stop() //nolint:errcheck

	state := w.State()
	if state.ContextEstimate.SkillCount != 2 {
		t.Errorf("skill count = %d, want 2", state.ContextEstimate.SkillCount)
	}
	if state.ContextEstimate.TotalInstructionBytes == 0 {
		t.Error("expected non-zero instruction bytes")
	}
	if state.ContextEstimate.EstimatedTokens == 0 {
		t.Error("expected non-zero estimated tokens")
	}
}

func TestIsIgnoredFile(t *testing.T) {
	tests := []struct {
		path    string
		ignored bool
	}{
		{"/tmp/file.swp", true},
		{"/tmp/file~", true},
		{"/tmp/file.tmp", true},
		{"/tmp/4913", true},
		{"/tmp/.DS_Store", true},
		{"/tmp/.git/HEAD", true},
		{"/tmp/docs/brainstorms/test.md", false},
		{"/tmp/.github/skills/test/SKILL.md", false},
	}
	for _, tt := range tests {
		t.Run(filepath.Base(tt.path), func(t *testing.T) {
			got := isIgnoredFile(tt.path)
			if got != tt.ignored {
				t.Errorf("isIgnoredFile(%q) = %v, want %v", tt.path, got, tt.ignored)
			}
		})
	}
}

func TestClassifyEventLayer(t *testing.T) {
	root := t.TempDir()
	w := &Watcher{root: root}

	tests := []struct {
		path  string
		layer string
	}{
		{filepath.Join(root, "docs", "brainstorms", "test.md"), "memory"},
		{filepath.Join(root, "docs", "plans", "test.md"), "memory"},
		{filepath.Join(root, "docs", "solutions", "test.md"), "memory"},
		{filepath.Join(root, ".github", "copilot-instructions.md"), "context"},
		{filepath.Join(root, ".github", "skills", "test", "SKILL.md"), "context"},
		{filepath.Join(root, ".github", "agents", "test.agent.md"), "context"},
		{filepath.Join(root, ".github", "prompts", "test.prompt.md"), "context"},
		{filepath.Join(root, ".atv", "install-manifest.json"), "health"},
		{filepath.Join(root, ".gstack", "test"), "health"},
		{filepath.Join(root, ".copilot-memory", "test.md"), "memory"},
		{filepath.Join(root, "main.go"), ""},
		{filepath.Join(root, "README.md"), ""},
	}
	for _, tt := range tests {
		t.Run(filepath.Base(tt.path), func(t *testing.T) {
			got := w.classifyEventLayer(tt.path)
			if got != tt.layer {
				t.Errorf("classifyEventLayer(%q) = %q, want %q", tt.path, got, tt.layer)
			}
		})
	}
}

func TestScanArtifactDir(t *testing.T) {
	dir := t.TempDir()

	// Create test files
	for _, f := range []struct {
		name    string
		content string
	}{
		{"a.md", "# A"},
		{"b.md", "# B content here"},
		{"not-md.txt", "skip"},
	} {
		if err := os.WriteFile(filepath.Join(dir, f.name), []byte(f.content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0o755); err != nil {
		t.Fatal(err)
	}

	artifacts := scanArtifactDir(dir)
	if len(artifacts) != 2 {
		t.Fatalf("got %d artifacts, want 2", len(artifacts))
	}

	// Should be sorted by ReadDir (alphabetical)
	if artifacts[0].Name != "a.md" {
		t.Errorf("first artifact = %q, want %q", artifacts[0].Name, "a.md")
	}
	if artifacts[0].Size != 3 { // "# A"
		t.Errorf("first artifact size = %d, want 3", artifacts[0].Size)
	}
}

func TestScanArtifactDir_Missing(t *testing.T) {
	artifacts := scanArtifactDir("/nonexistent/path")
	if artifacts != nil {
		t.Errorf("expected nil for missing dir, got %v", artifacts)
	}
}
