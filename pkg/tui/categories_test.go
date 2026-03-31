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
