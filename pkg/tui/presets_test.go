package tui

import (
	"strings"
	"testing"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/gstack"
)

func TestPresetBrowserAndRuntimeBoundaries(t *testing.T) {
	if ProPreset.IncludeAgentBrowser {
		t.Fatal("Pro preset should not include agent-browser by default")
	}
	if ProPreset.EnableGstackRuntime {
		t.Fatal("Pro preset should remain markdown-only by default")
	}
	if !FullPreset.IncludeAgentBrowser {
		t.Fatal("Full preset should include agent-browser")
	}
	if !FullPreset.EnableGstackRuntime {
		t.Fatal("Full preset should enable gstack runtime when prerequisites exist")
	}
}

func TestShouldEnableGstackRuntime(t *testing.T) {
	withBun := gstack.Prerequisites{HasBun: true}
	withoutBun := gstack.Prerequisites{}

	if ProPreset.ShouldEnableGstackRuntime(ProPreset.GstackDirs, withBun) {
		t.Fatal("Pro preset should not enable runtime even when Bun is available")
	}
	if FullPreset.ShouldEnableGstackRuntime(nil, withBun) {
		t.Fatal("Full preset should not enable runtime without selected gstack skills")
	}
	if FullPreset.ShouldEnableGstackRuntime(FullPreset.GstackDirs, withoutBun) {
		t.Fatal("Full preset should not enable runtime when Bun is missing")
	}
	if !FullPreset.ShouldEnableGstackRuntime(FullPreset.GstackDirs, withBun) {
		t.Fatal("Full preset should enable runtime when Bun is available and gstack skills are selected")
	}
}

func TestPresetPreviewLabelExplainsFallbacks(t *testing.T) {
	withBun := gstack.Prerequisites{HasGit: true, HasBun: true}
	withoutBun := gstack.Prerequisites{HasGit: true}

	fullReady := FullPreset.PreviewLabel(withBun)
	if !strings.Contains(fullReady, "git + Bun ready") {
		t.Fatalf("Full preset preview should mention ready Bun state, got %q", fullReady)
	}

	fullFallback := FullPreset.PreviewLabel(withoutBun)
	if !strings.Contains(fullFallback, "Bun missing → docs-only gstack fallback") {
		t.Fatalf("Full preset preview should mention Bun fallback, got %q", fullFallback)
	}

	proFallback := ProPreset.PreviewLabel(gstack.Prerequisites{})
	if !strings.Contains(proFallback, "git missing → gstack sync blocked") {
		t.Fatalf("Pro preset preview should mention missing git, got %q", proFallback)
	}
}

func TestPresetSelectionSummaryExplainsCapabilitiesAndDowngrades(t *testing.T) {
	summary := FullPreset.SelectionSummary(gstack.Prerequisites{HasGit: true})
	for _, want := range []string{
		"What it adds:",
		"Prerequisites:",
		"If tools are missing:",
		"Detected now:",
	} {
		if !strings.Contains(summary, want) {
			t.Fatalf("selection summary missing %q in %q", want, summary)
		}
	}
}
