package tui

import (
	"errors"
	"testing"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
)

func TestParseStackPackSelectionsNormalizesOrder(t *testing.T) {
	packs, err := parseStackPackSelections([]string{"rails", "general", "typescript", "rails"})
	if err != nil {
		t.Fatalf("parseStackPackSelections() error = %v", err)
	}
	want := []installstate.StackPack{
		installstate.StackPackGeneral,
		installstate.StackPackTypeScript,
		installstate.StackPackRails,
	}
	if len(packs) != len(want) {
		t.Fatalf("expected %d packs, got %d: %v", len(want), len(packs), packs)
	}
	for i := range want {
		if packs[i] != want[i] {
			t.Fatalf("pack %d = %s, want %s", i, packs[i], want[i])
		}
	}
}

func TestParseStackPackSelectionsRequiresAtLeastOne(t *testing.T) {
	_, err := parseStackPackSelections(nil)
	if !errors.Is(err, installstate.ErrNoStackPacks) {
		t.Fatalf("parseStackPackSelections(nil) error = %v, want %v", err, installstate.ErrNoStackPacks)
	}
}
