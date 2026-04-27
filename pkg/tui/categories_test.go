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
		var hasAutoresearch bool
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
			if skill.Key == "core-skills:autoresearch" {
				hasAutoresearch = true
				if skill.Source != "atv" {
					t.Fatalf("autoresearch should have source atv, got %q", skill.Source)
				}
				if skill.IsGstack {
					t.Fatal("autoresearch should not be a gstack skill")
				}
			}
		}
		if !hasKarpathy {
			t.Fatal("guidelines category should contain karpathy-guidelines skill")
		}
		if !hasAutoresearch {
			t.Fatal("guidelines category should contain autoresearch skill")
		}
	}
	if !found {
		t.Fatal("expected guidelines category group in BuildCategoryGroups output")
	}
}

// TestBuildCategoryGroupsIncludesMaintenanceCategory guards that the
// `🩺 Maintenance & Health` category surfaces both atv-doctor and
// atv-update in the customize-mode TUI. Both ship via LayerCoreSkills
// so users would receive them silently without a TUI entry.
func TestBuildCategoryGroupsIncludesMaintenanceCategory(t *testing.T) {
	groups := BuildCategoryGroups(gstack.Prerequisites{})
	var found bool
	for _, group := range groups {
		if group.Category != gstack.CategoryMaintenance {
			continue
		}
		found = true
		var hasDoctor, hasUpdate bool
		for _, skill := range group.Skills {
			switch skill.Key {
			case "core-skills:atv-doctor":
				hasDoctor = true
				if skill.Source != "atv" {
					t.Errorf("atv-doctor should have source atv, got %q", skill.Source)
				}
				if skill.IsGstack {
					t.Error("atv-doctor should not be a gstack skill")
				}
			case "core-skills:atv-update":
				hasUpdate = true
				if skill.Source != "atv" {
					t.Errorf("atv-update should have source atv, got %q", skill.Source)
				}
				if skill.IsGstack {
					t.Error("atv-update should not be a gstack skill")
				}
			}
		}
		if !hasDoctor {
			t.Error("maintenance category should contain atv-doctor")
		}
		if !hasUpdate {
			t.Error("maintenance category should contain atv-update")
		}
	}
	if !found {
		t.Fatal("expected maintenance category group in BuildCategoryGroups output")
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

// TestSecurityCategoryIncludesAtvSecurity ensures the customize-mode TUI
// surfaces atv-security as a toggleable option. atv-security ships via
// LayerCoreSkills, so without a TUI entry users would receive it silently
// with no way to opt out.
//
// Also asserts that the legacy `core-skills:cso` key is ABSENT — `/cso` was
// folded into `/atv-security`. A re-introduction would create the same
// name-collision problem with gstack's `/cso` skill that motivated the merge.
func TestSecurityCategoryIncludesAtvSecurity(t *testing.T) {
	groups := BuildCategoryGroups(gstack.Prerequisites{})

	var security *CategoryGroup
	for i, group := range groups {
		if group.Category == gstack.CategorySecurity {
			security = &groups[i]
			break
		}
	}
	if security == nil {
		t.Fatal("expected security category group in BuildCategoryGroups output")
	}

	var foundAtvSecurity bool
	for _, skill := range security.Skills {
		switch skill.Key {
		case "core-skills:atv-security":
			foundAtvSecurity = true
			if skill.Source != "atv" {
				t.Errorf("core-skills:atv-security should have source=atv, got %q", skill.Source)
			}
		case "core-skills:cso":
			t.Error("core-skills:cso must not reappear in the security category — it was folded into /atv-security to avoid collision with gstack's /cso")
		}
	}
	if !foundAtvSecurity {
		t.Error("security category missing required skill: core-skills:atv-security")
	}
}

// TestShippingCategoryIncludesLandAndTakeoff is the customize-mode counterpart
// to TestCoreLayerShipsLandAndTakeoff in pkg/scaffold. It guards against a
// regression where Land or Takeoff are dropped from the Shipping category in
// the customize-mode TUI even if they remain wired in the catalog.
func TestShippingCategoryIncludesLandAndTakeoff(t *testing.T) {
	groups := BuildCategoryGroups(gstack.Prerequisites{})

	var shipping *CategoryGroup
	for i, group := range groups {
		if group.Category == gstack.CategoryShipping {
			shipping = &groups[i]
			break
		}
	}
	if shipping == nil {
		t.Fatal("expected shipping category group in BuildCategoryGroups output")
	}

	want := map[string]bool{
		"core-skills:takeoff": false,
		"core-skills:land":    false,
	}
	for _, skill := range shipping.Skills {
		if _, ok := want[skill.Key]; ok {
			want[skill.Key] = true
			if skill.Source != "atv" {
				t.Errorf("%s should have source=atv, got %q", skill.Key, skill.Source)
			}
			if skill.IsGstack {
				t.Errorf("%s should not be marked as a gstack skill", skill.Key)
			}
		}
	}
	for key, found := range want {
		if !found {
			t.Errorf("shipping category missing required skill: %s", key)
		}
	}
}
