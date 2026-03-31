package installstate

import (
	"path/filepath"
	"testing"
	"time"
)

func TestBuildLaunchpadSnapshotWithoutManifestStillComputesRecommendations(t *testing.T) {
	root := t.TempDir()
	mustWriteMarkdown(t, filepath.Join(root, "docs", "brainstorms", "idea.md"), "# Idea")

	snapshot, err := BuildLaunchpadSnapshot(root)
	if err != nil {
		t.Fatalf("BuildLaunchpadSnapshot() error = %v", err)
	}
	if snapshot.HasManifest {
		t.Fatal("snapshot should not report a manifest when none exists")
	}
	if len(snapshot.Recommendations) == 0 || snapshot.Recommendations[0].ID != "turn-brainstorm-into-plan" {
		t.Fatalf("unexpected recommendations: %+v", snapshot.Recommendations)
	}
}

func TestBuildLaunchpadSnapshotWithManifestIncludesRequestedState(t *testing.T) {
	root := t.TempDir()
	mustWriteMarkdown(t, filepath.Join(root, "docs", "brainstorms", "idea.md"), "# Idea")
	mustWriteMarkdown(t, filepath.Join(root, "docs", "plans", "work.md"), "status: active\n- [ ] do the work\n")

	manifest := InstallManifest{
		GeneratedAt: time.Date(2026, time.March, 31, 12, 0, 0, 0, time.UTC),
		Requested: RequestedState{
			StackPacks:          []StackPack{StackPackGeneral, StackPackTypeScript},
			GstackDirs:          []string{"review"},
			IncludeAgentBrowser: true,
			PresetName:          "Full",
		},
		Outcomes: []InstallOutcome{
			{Step: "Scaffolding ATV files", Status: InstallStepDone},
			{Step: "Syncing gstack skills", Status: InstallStepWarning},
		},
	}
	if err := WriteManifest(root, manifest); err != nil {
		t.Fatalf("WriteManifest() error = %v", err)
	}

	snapshot, err := BuildLaunchpadSnapshot(root)
	if err != nil {
		t.Fatalf("BuildLaunchpadSnapshot() error = %v", err)
	}
	if !snapshot.HasManifest {
		t.Fatal("snapshot should report manifest presence")
	}
	if snapshot.Requested.PresetName != "Full" {
		t.Fatalf("PresetName = %q, want Full", snapshot.Requested.PresetName)
	}
	if snapshot.OutcomeSummary.Done != 1 || snapshot.OutcomeSummary.Warning != 1 {
		t.Fatalf("unexpected outcome summary: %+v", snapshot.OutcomeSummary)
	}
	if !snapshot.HasGstack() || !snapshot.HasAgentBrowser() {
		t.Fatalf("expected manifest capabilities to be reflected: %+v", snapshot.Requested)
	}
	if len(snapshot.StackPackLabels()) != 2 || snapshot.StackPackLabels()[1] != "TypeScript" {
		t.Fatalf("unexpected stack pack labels: %v", snapshot.StackPackLabels())
	}
}
