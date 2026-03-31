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
