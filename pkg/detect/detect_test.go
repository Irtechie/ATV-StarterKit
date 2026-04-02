package detect

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
)

func TestDetectEnvironmentFindsMultipleLikelyPacks(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "Gemfile"), "source 'https://rubygems.org'")
	mustWriteFile(t, filepath.Join(dir, "config", "routes.rb"), "Rails.application.routes.draw do end")
	mustWriteFile(t, filepath.Join(dir, "tsconfig.json"), "{}")

	env := DetectEnvironment(dir)
	if env.Stack != StackRails {
		t.Fatalf("primary stack = %s, want %s", env.Stack, StackRails)
	}
	want := []installstate.StackPack{installstate.StackPackTypeScript, installstate.StackPackRails}
	if len(env.DetectedPacks) != len(want) {
		t.Fatalf("detected packs = %v, want %v", env.DetectedPacks, want)
	}
	for i := range want {
		if env.DetectedPacks[i] != want[i] {
			t.Fatalf("detected pack %d = %s, want %s", i, env.DetectedPacks[i], want[i])
		}
	}
}

func TestDetectEnvironmentFallsBackToGeneralPack(t *testing.T) {
	env := DetectEnvironment(t.TempDir())
	if env.Stack != StackGeneral {
		t.Fatalf("primary stack = %s, want %s", env.Stack, StackGeneral)
	}
	if len(env.DetectedPacks) != 1 || env.DetectedPacks[0] != installstate.StackPackGeneral {
		t.Fatalf("detected packs = %v, want [general]", env.DetectedPacks)
	}
}

func TestPrimaryStackForPacksUsesPreferredWhenAvailable(t *testing.T) {
	got := PrimaryStackForPacks(
		[]installstate.StackPack{installstate.StackPackGeneral, installstate.StackPackTypeScript, installstate.StackPackRails},
		StackTypeScript,
	)
	if got != StackTypeScript {
		t.Fatalf("PrimaryStackForPacks() = %s, want %s", got, StackTypeScript)
	}
}

func TestPrimaryStackForPacksFallsBackDeterministically(t *testing.T) {
	got := PrimaryStackForPacks(
		[]installstate.StackPack{installstate.StackPackGeneral, installstate.StackPackPython},
		StackRails,
	)
	if got != StackPython {
		t.Fatalf("PrimaryStackForPacks() = %s, want %s", got, StackPython)
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
