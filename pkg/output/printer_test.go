package output

import (
	"strings"
	"testing"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
)

func TestGuidedSummaryTextIncludesManifestAndOutcomeReasons(t *testing.T) {
	text := guidedSummaryText([]installstate.InstallOutcome{
		{Step: "Scaffolding ATV files", Status: installstate.InstallStepDone, Detail: "4 files created", Duration: "120ms"},
		{Step: "Syncing gstack skills", Status: installstate.InstallStepWarning, Reason: "fell back to markdown-only sync"},
	}, ".atv/install-manifest.json")

	for _, want := range []string{
		"Guided install summary",
		"Scaffolding ATV files",
		"Syncing gstack skills",
		"fell back to markdown-only sync",
		".atv/install-manifest.json",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("summary missing %q in %q", want, text)
		}
	}
}

func TestLaunchpadTextIncludesManifestAndRecommendations(t *testing.T) {
	text := launchpadText(installstate.LaunchpadSnapshot{
		Root:         "C:/repo/example",
		ManifestPath: ".atv/install-manifest.json",
		HasManifest:  true,
		Requested: installstate.RequestedState{
			StackPacks:          []installstate.StackPack{installstate.StackPackGeneral, installstate.StackPackTypeScript},
			GstackDirs:          []string{"review"},
			IncludeAgentBrowser: true,
			PresetName:          "Full",
		},
		RepoState:      installstate.RepoState{BrainstormCount: 1, PlanCount: 2, SolutionCount: 0, HasUncheckedPlan: true},
		OutcomeSummary: installstate.OutcomeSummary{Done: 2, Warning: 1},
		Recommendations: []installstate.Recommendation{{
			ID:       "execute-active-plan",
			Title:    "Continue the active plan with /ce-work",
			Reason:   "At least one plan still has unchecked items.",
			Priority: 80,
		}},
	})

	for _, want := range []string{
		"ATV Launchpad",
		"Installed intelligence",
		"Repo memory snapshot",
		"Recommended next moves",
		"Continue the active plan with /ce-work",
		"atv-installer launchpad",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("launchpad missing %q in %q", want, text)
		}
	}
}

func TestLaunchpadTextHandlesMissingManifest(t *testing.T) {
	text := launchpadText(installstate.LaunchpadSnapshot{
		ManifestPath: ".atv/install-manifest.json",
		RepoState:    installstate.RepoState{BrainstormCount: 0},
		Recommendations: []installstate.Recommendation{{
			ID:       "start-brainstorm",
			Title:    "Start with /ce-brainstorm to shape the first feature",
			Reason:   "No brainstorms were found in docs/brainstorms yet.",
			Priority: 90,
		}},
	})

	for _, want := range []string{
		"No guided manifest found yet",
		"atv-installer init --guided",
		"Start with /ce-brainstorm to shape the first feature",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("launchpad missing %q in %q", want, text)
		}
	}
}
