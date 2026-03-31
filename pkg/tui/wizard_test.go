package tui

import (
	"errors"
	"strings"
	"testing"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/gstack"
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

func TestDefaultSelectedSkillKeysSkipsRuntimeSelectionsWithoutBun(t *testing.T) {
	groups := BuildCategoryGroups(gstack.Prerequisites{})
	presetGstackSet := make(map[string]bool)
	for _, dir := range FullPreset.GstackDirs {
		presetGstackSet[dir] = true
	}

	for _, group := range groups {
		if group.Category != gstack.CategoryQATesting {
			continue
		}
		selected := defaultSelectedSkillKeys(group, presetGstackSet, gstack.Prerequisites{})
		for _, key := range selected {
			if strings.HasPrefix(key, "gstack:") {
				t.Fatalf("runtime QA gstack skill %q should not be preselected without Bun", key)
			}
		}
		return
	}
	t.Fatal("expected QA category group")
}

func TestSkillOptionLabelIncludesRuntimeHints(t *testing.T) {
	label := skillOptionLabel(CategorySkill{Label: "QA — browser checks", IsGstack: true, RequiresBun: true}, gstack.Prerequisites{})
	if !strings.Contains(label, "requires Bun") {
		t.Fatalf("expected Bun hint in %q", label)
	}
}
