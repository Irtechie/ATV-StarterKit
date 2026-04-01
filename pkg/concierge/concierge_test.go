package concierge

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
)

func TestGetMemorySummaryWithoutManifest(t *testing.T) {
	root := t.TempDir()
	summary := GetMemorySummary(root)

	if summary.Status != "no-manifest" {
		t.Fatalf("expected no-manifest status, got %q", summary.Status)
	}
	if summary.Manifest != nil {
		t.Fatal("should not have manifest info when no manifest exists")
	}
	if summary.Message == "" {
		t.Fatal("degraded state should include a helpful message")
	}
}

func TestGetMemorySummaryWithManifest(t *testing.T) {
	root := t.TempDir()
	writeTestManifest(t, root)

	summary := GetMemorySummary(root)

	if summary.Status != "ok" {
		t.Fatalf("expected ok status, got %q", summary.Status)
	}
	if summary.Manifest == nil {
		t.Fatal("should have manifest info")
	}
	if summary.Manifest.PresetName != "Pro" {
		t.Fatalf("expected Pro preset, got %q", summary.Manifest.PresetName)
	}
	if len(summary.Manifest.StackPacks) != 2 {
		t.Fatalf("expected 2 stack packs, got %d", len(summary.Manifest.StackPacks))
	}
	if summary.Manifest.OutcomeSummary.Done != 1 {
		t.Fatalf("expected 1 done outcome, got %d", summary.Manifest.OutcomeSummary.Done)
	}
}

func TestGetMemorySummaryIncludesRepoState(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "docs", "brainstorms", "idea.md"), "# Idea")
	mustWrite(t, filepath.Join(root, ".github", "copilot-instructions.md"), "# Instructions")
	os.MkdirAll(filepath.Join(root, ".github", "agents", "reviewer"), 0o755)
	os.MkdirAll(filepath.Join(root, ".github", "skills", "brainstorming"), 0o755)

	summary := GetMemorySummary(root)

	if summary.RepoState.BrainstormCount != 1 {
		t.Fatalf("expected 1 brainstorm, got %d", summary.RepoState.BrainstormCount)
	}
	if !summary.RepoState.HasCopilotInstructions {
		t.Fatal("should detect copilot-instructions.md")
	}
	if summary.RepoState.InstalledAgents != 1 {
		t.Fatalf("expected 1 agent, got %d", summary.RepoState.InstalledAgents)
	}
	if summary.RepoState.InstalledSkills != 1 {
		t.Fatalf("expected 1 skill, got %d", summary.RepoState.InstalledSkills)
	}
}

func TestListRecommendationsReturnsLocalDeterministic(t *testing.T) {
	root := t.TempDir()
	list := ListRecommendations(root)

	if list.Status != "ok" {
		t.Fatalf("expected ok status, got %q", list.Status)
	}
	if list.Source != "local-deterministic" {
		t.Fatalf("expected local-deterministic source, got %q", list.Source)
	}
	if len(list.Recommendations) == 0 {
		t.Fatal("empty repo should still produce recommendations")
	}
}

func TestExplainRecommendationFound(t *testing.T) {
	root := t.TempDir()
	detail := ExplainRecommendation(root, "start-brainstorm")

	if detail.Status != "ok" {
		t.Fatalf("expected ok status, got %q", detail.Status)
	}
	if detail.ID != "start-brainstorm" {
		t.Fatalf("expected start-brainstorm id, got %q", detail.ID)
	}
	if detail.SuggestedCmd == "" {
		t.Fatal("should have a suggested command")
	}
}

func TestExplainRecommendationNotFound(t *testing.T) {
	root := t.TempDir()
	detail := ExplainRecommendation(root, "nonexistent-id")

	if detail.Status != "not-found" {
		t.Fatalf("expected not-found status, got %q", detail.Status)
	}
}

func TestOpenArtifactKnown(t *testing.T) {
	root := t.TempDir()

	// Test existing artifact
	mustWrite(t, filepath.Join(root, ".github", "copilot-instructions.md"), "# test")
	info := OpenArtifact(root, "instructions")
	if info.Status != "ok" {
		t.Fatalf("expected ok status, got %q", info.Status)
	}
	if !info.Exists {
		t.Fatal("instructions should exist")
	}
	if info.Type != "file" {
		t.Fatalf("expected file type, got %q", info.Type)
	}

	// Test non-existing but known artifact
	info = OpenArtifact(root, "manifest")
	if info.Status != "ok" {
		t.Fatalf("expected ok status, got %q", info.Status)
	}
	if info.Exists {
		t.Fatal("manifest should not exist yet")
	}
}

func TestOpenArtifactUnknown(t *testing.T) {
	root := t.TempDir()
	info := OpenArtifact(root, "foobar")

	if info.Status != "unknown" {
		t.Fatalf("expected unknown status, got %q", info.Status)
	}
}

func TestRunSuggestedActionKnown(t *testing.T) {
	root := t.TempDir()
	result := RunSuggestedAction(root, "start-brainstorm")

	if result.Status != "ready" {
		t.Fatalf("expected ready status, got %q", result.Status)
	}
	if result.Message == "" {
		t.Fatal("should include a suggested command message")
	}
}

func TestRunSuggestedActionUnknown(t *testing.T) {
	root := t.TempDir()
	result := RunSuggestedAction(root, "nonexistent")

	if result.Status != "unknown" {
		t.Fatalf("expected unknown status, got %q", result.Status)
	}
}

func TestAssistantCannotOverrideDeterministicRanking(t *testing.T) {
	// Verify that ListRecommendations always returns the same order
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "docs", "brainstorms", "idea.md"), "# Idea")
	writeTestManifest(t, root)

	first := ListRecommendations(root)
	second := ListRecommendations(root)

	if len(first.Recommendations) != len(second.Recommendations) {
		t.Fatal("recommendations should be deterministic")
	}
	for i := range first.Recommendations {
		if first.Recommendations[i].ID != second.Recommendations[i].ID {
			t.Fatalf("recommendation order differs at %d", i)
		}
	}
}

func TestCoreValueUnchangedWithoutAssistant(t *testing.T) {
	// The concierge tools must return the same data as the raw installstate
	// functions — the assistant layer adds no new truth.
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "docs", "brainstorms", "idea.md"), "# Idea")
	writeTestManifest(t, root)

	snapshot, err := installstate.BuildLaunchpadSnapshot(root)
	if err != nil {
		t.Fatalf("BuildLaunchpadSnapshot() error = %v", err)
	}

	list := ListRecommendations(root)
	if len(list.Recommendations) != len(snapshot.Recommendations) {
		t.Fatalf("concierge recommendations (%d) differs from raw launchpad (%d)",
			len(list.Recommendations), len(snapshot.Recommendations))
	}
	for i := range list.Recommendations {
		if list.Recommendations[i].ID != snapshot.Recommendations[i].ID {
			t.Fatalf("concierge recommendation %d differs from launchpad", i)
		}
	}
}

func writeTestManifest(t *testing.T, root string) {
	t.Helper()
	manifest := installstate.InstallManifest{
		Requested: installstate.RequestedState{
			StackPacks:          []installstate.StackPack{installstate.StackPackGeneral, installstate.StackPackTypeScript},
			GstackDirs:          []string{"review"},
			IncludeAgentBrowser: false,
			PresetName:          "Pro",
		},
		Outcomes: []installstate.InstallOutcome{
			{Step: "Scaffolding ATV files", Status: installstate.InstallStepDone},
		},
	}
	if err := installstate.WriteManifest(root, manifest); err != nil {
		t.Fatalf("WriteManifest() error = %v", err)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
