package tui

import (
	"strings"
	"testing"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/gstack"
)

func TestBuildCategoryGroupsIncludesDescriptions(t *testing.T) {
	groups := BuildCategoryGroups(gstack.Prerequisites{HasBun: true})
	if len(groups) == 0 {
		t.Fatal("expected category groups")
	}

	for _, group := range groups {
		if strings.TrimSpace(group.Description) == "" {
			t.Fatalf("group %s should have a description", group.Category)
		}
	}
}

func TestBuildCategoryGroupsIncludesGuidelinesCategory(t *testing.T) {
	groups := BuildCategoryGroups(gstack.Prerequisites{})
	var found bool
	for _, group := range groups {
		if group.Category != gstack.CategoryGuidelines {
			continue
		}
		found = true
		// Must contain the Karpathy skill
		var hasKarpathy bool
		for _, skill := range group.Skills {
			if skill.Key == "core-skills:karpathy-guidelines" {
				hasKarpathy = true
				if skill.Source != "atv" {
					t.Fatalf("karpathy-guidelines should have source atv, got %q", skill.Source)
				}
				if skill.IsGstack {
					t.Fatal("karpathy-guidelines should not be a gstack skill")
				}
			}
		}
		if !hasKarpathy {
			t.Fatal("guidelines category should contain karpathy-guidelines skill")
		}
	}
	if !found {
		t.Fatal("expected guidelines category group in BuildCategoryGroups output")
	}
}

func TestKarpathyGuidedFlowParseSelections(t *testing.T) {
	// Simulate what happens when user selects karpathy-guidelines in Screen 4
	selected := []string{"core-skills:karpathy-guidelines", "core-skills:ce-plan"}
	atvLayers, gstackDirs := ParseSelections(selected)

	if len(gstackDirs) != 0 {
		t.Fatalf("karpathy selection should not produce gstack dirs, got %v", gstackDirs)
	}
	// Both keys map to the same "core-skills" layer, deduplicated
	if len(atvLayers) != 1 || atvLayers[0] != "core-skills" {
		t.Fatalf("expected [core-skills], got %v", atvLayers)
	}
}

func TestBuildCategoryGroupsWarnsWhenQARuntimeUnavailable(t *testing.T) {
	groups := BuildCategoryGroups(gstack.Prerequisites{})
	for _, group := range groups {
		if group.Category != gstack.CategoryQATesting {
			continue
		}
		if !strings.Contains(group.Description, "docs-only") {
			t.Fatalf("QA category should mention docs-only fallback when Bun is missing, got %q", group.Description)
		}
		return
	}
	t.Fatal("expected QA category group")
}
