package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
)

func TestFormatStepDuration(t *testing.T) {
	if got := formatStepDuration(0); got != "" {
		t.Fatalf("formatStepDuration(0) = %q, want empty", got)
	}
	if got := formatStepDuration(250 * time.Millisecond); !strings.Contains(got, "ms") {
		t.Fatalf("expected millisecond duration, got %q", got)
	}
	if got := formatStepDuration(1500 * time.Millisecond); !strings.Contains(got, "s") {
		t.Fatalf("expected second duration, got %q", got)
	}
}

func TestProgressModelCapturesSkippedOutcomeForNilAction(t *testing.T) {
	model := ProgressModel{
		steps: []InstallStep{{Name: "Skip me", Action: nil}},
	}
	model.current = 0
	model.steps[0].Status = StepDone
	model.outcomes = append(model.outcomes, installstate.InstallOutcome{Step: "Skip me", Status: installstate.InstallStepSkipped})
	if len(model.outcomes) != 1 || model.outcomes[0].Status != installstate.InstallStepSkipped {
		t.Fatalf("unexpected outcomes: %+v", model.outcomes)
	}
}
