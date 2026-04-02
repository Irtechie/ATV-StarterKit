package installstate

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeStackPacks(t *testing.T) {
	normalized, err := NormalizeStackPacks([]StackPack{
		StackPackPython,
		StackPackGeneral,
		StackPackPython,
		StackPackRails,
	})
	if err != nil {
		t.Fatalf("NormalizeStackPacks() error = %v", err)
	}

	want := []StackPack{StackPackGeneral, StackPackPython, StackPackRails}
	if len(normalized) != len(want) {
		t.Fatalf("expected %d packs, got %d: %v", len(want), len(normalized), normalized)
	}
	for i := range want {
		if normalized[i] != want[i] {
			t.Fatalf("pack %d = %s, want %s", i, normalized[i], want[i])
		}
	}
}

func TestValidateStackPacksRequiresAtLeastOne(t *testing.T) {
	err := ValidateStackPacks(nil)
	if !errors.Is(err, ErrNoStackPacks) {
		t.Fatalf("ValidateStackPacks(nil) error = %v, want %v", err, ErrNoStackPacks)
	}
}

func TestWriteManifestRoundTrip(t *testing.T) {
	root := t.TempDir()
	manifest := InstallManifest{
		Requested: RequestedState{
			StackPacks:          []StackPack{StackPackGeneral, StackPackTypeScript},
			ATVLayers:           []string{"core-skills", "setup-steps"},
			GstackDirs:          []string{"review"},
			GstackRuntime:       false,
			IncludeAgentBrowser: false,
			PresetName:          "Pro",
		},
		Outcomes: []InstallOutcome{{
			Step:   "Scaffolding ATV files",
			Status: InstallStepDone,
		}},
	}

	if err := WriteManifest(root, manifest); err != nil {
		t.Fatalf("WriteManifest() error = %v", err)
	}

	path := ManifestPath(root)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("manifest path should exist: %v", err)
	}

	loaded, err := ReadManifest(root)
	if err != nil {
		t.Fatalf("ReadManifest() error = %v", err)
	}
	if loaded.Version != ManifestVersion {
		t.Fatalf("Version = %d, want %d", loaded.Version, ManifestVersion)
	}
	if loaded.RerunPolicy != RerunPolicyAdditiveOnly {
		t.Fatalf("RerunPolicy = %s, want %s", loaded.RerunPolicy, RerunPolicyAdditiveOnly)
	}
	if loaded.GeneratedAt.IsZero() {
		t.Fatal("GeneratedAt should be populated")
	}
	if len(loaded.Requested.StackPacks) != 2 {
		t.Fatalf("expected 2 stack packs, got %d", len(loaded.Requested.StackPacks))
	}
	if loaded.Requested.StackPacks[0] != StackPackGeneral || loaded.Requested.StackPacks[1] != StackPackTypeScript {
		t.Fatalf("unexpected stack packs: %v", loaded.Requested.StackPacks)
	}
	if len(loaded.Outcomes) != 1 || loaded.Outcomes[0].Status != InstallStepDone {
		t.Fatalf("unexpected outcomes: %+v", loaded.Outcomes)
	}
}

func TestManifestPathUsesRepoLocalATVDirectory(t *testing.T) {
	root := filepath.Join("C:\\repo", "example")
	got := ManifestPath(root)
	want := filepath.Join(root, ".atv", "install-manifest.json")
	if got != want {
		t.Fatalf("ManifestPath() = %s, want %s", got, want)
	}
}
