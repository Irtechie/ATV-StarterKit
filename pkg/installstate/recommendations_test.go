package installstate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildRecommendationsForEmptyRepo(t *testing.T) {
	root := t.TempDir()
	recs := BuildRecommendations(root, InstallManifest{})
	if len(recs) == 0 || recs[0].ID != "start-brainstorm" {
		t.Fatalf("expected start-brainstorm recommendation first, got %+v", recs)
	}
}

func TestBuildRecommendationsPrefersActivePlanWork(t *testing.T) {
	root := t.TempDir()
	mustWriteMarkdown(t, filepath.Join(root, "docs", "brainstorms", "idea.md"), "# Idea")
	mustWriteMarkdown(t, filepath.Join(root, "docs", "plans", "work.md"), "status: active\n- [ ] do the work\n")

	manifest := InstallManifest{
		Requested: RequestedState{GstackDirs: []string{"review"}},
		Outcomes:  []InstallOutcome{{Step: "Syncing gstack skills", Status: InstallStepDone}},
	}

	recs := BuildRecommendations(root, manifest)
	if len(recs) == 0 || recs[0].ID != "execute-active-plan" {
		t.Fatalf("expected execute-active-plan recommendation first, got %+v", recs)
	}
	if len(recs) < 2 || recs[1].ID != "start-gstack-sprint" {
		t.Fatalf("expected gstack recommendation second, got %+v", recs)
	}
}

func TestBuildRecommendationsPrioritizesInstallProblems(t *testing.T) {
	root := t.TempDir()
	mustWriteMarkdown(t, filepath.Join(root, "docs", "brainstorms", "idea.md"), "# Idea")

	manifest := InstallManifest{
		Outcomes: []InstallOutcome{{Step: "Installing agent-browser", Status: InstallStepWarning, Reason: "npm not found"}},
	}

	recs := BuildRecommendations(root, manifest)
	if len(recs) == 0 || recs[0].ID != "fix-install-issues" {
		t.Fatalf("expected fix-install-issues first, got %+v", recs)
	}
}

func mustWriteMarkdown(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
