package tui

import (
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
