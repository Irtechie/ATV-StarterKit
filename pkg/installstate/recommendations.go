package installstate

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// RepoState captures lightweight local facts used by deterministic recommendations.
type RepoState struct {
	BrainstormCount  int
	PlanCount        int
	SolutionCount    int
	HasUncheckedPlan bool
	HasCompletedPlan bool

	// Installed intelligence classification
	InstalledAgents        int
	InstalledSkills        int
	HasCopilotInstructions bool
	HasGstackStaging       bool
	HasAgentBrowserSkill   bool
}

// ScanRepoState inspects the local docs tree to find brainstorms, plans, and solutions.
func ScanRepoState(root string) RepoState {
	state := RepoState{}
	state.BrainstormCount = countMarkdownFiles(filepath.Join(root, "docs", "brainstorms"))
	state.PlanCount, state.HasUncheckedPlan, state.HasCompletedPlan = scanPlans(filepath.Join(root, "docs", "plans"))
	state.SolutionCount = countMarkdownFiles(filepath.Join(root, "docs", "solutions"))

	// Classify installed intelligence
	state.InstalledAgents = countFilesWithSuffix(filepath.Join(root, ".github", "agents"), ".agent.md")
	state.InstalledSkills = countSubdirs(filepath.Join(root, ".github", "skills"))
	state.HasCopilotInstructions = fileExists(filepath.Join(root, ".github", "copilot-instructions.md"))
	state.HasGstackStaging = dirExists(filepath.Join(root, ".gstack"))
	state.HasAgentBrowserSkill = dirExists(filepath.Join(root, ".github", "skills", "agent-browser"))

	return state
}

// BuildRecommendations computes a small deterministic set of next moves from local facts.
func BuildRecommendations(root string, manifest InstallManifest) []Recommendation {
	state := ScanRepoState(root)
	var recommendations []Recommendation

	if outcome := firstProblemOutcome(manifest.Outcomes); outcome != nil {
		reason := outcome.Reason
		if reason == "" {
			reason = outcome.Detail
		}
		recommendations = append(recommendations, Recommendation{
			ID:       "fix-install-issues",
			Title:    "Fix installer warnings before relying on every capability",
			Reason:   reason,
			Priority: 100,
		})
	}

	switch {
	case state.BrainstormCount == 0:
		recommendations = append(recommendations, Recommendation{
			ID:       "start-brainstorm",
			Title:    `Start with /ce-brainstorm to shape the first feature`,
			Reason:   "No brainstorms were found in docs/brainstorms yet.",
			Priority: 90,
		})
	case state.PlanCount == 0:
		recommendations = append(recommendations, Recommendation{
			ID:       "turn-brainstorm-into-plan",
			Title:    `Turn the brainstorm into a plan with /ce-plan`,
			Reason:   "Brainstorms exist, but no plan files were found in docs/plans.",
			Priority: 85,
		})
	case state.HasUncheckedPlan:
		recommendations = append(recommendations, Recommendation{
			ID:       "execute-active-plan",
			Title:    `Continue the active plan with /ce-work`,
			Reason:   "At least one plan still has unchecked items.",
			Priority: 80,
		})
	case state.HasCompletedPlan && state.SolutionCount == 0:
		recommendations = append(recommendations, Recommendation{
			ID:       "compound-learnings",
			Title:    `Capture what shipped with /ce-compound`,
			Reason:   "Completed plans exist, but docs/solutions is still empty.",
			Priority: 70,
		})
	}

	if manifest.Requested.GstackDirs != nil && len(manifest.Requested.GstackDirs) > 0 && stepUsable(manifest.Outcomes, "gstack") {
		recommendations = append(recommendations, Recommendation{
			ID:       "start-gstack-sprint",
			Title:    `Use /gstack-office-hours for a deeper sprint kickoff`,
			Reason:   "gstack skills were requested and synced successfully enough to use.",
			Priority: 55,
		})
	}

	if manifest.Requested.IncludeAgentBrowser && stepUsable(manifest.Outcomes, "agent-browser") {
		recommendations = append(recommendations, Recommendation{
			ID:       "browser-check",
			Title:    "Open the app in a real browser with agent-browser",
			Reason:   "Browser automation tooling was installed or partially prepared.",
			Priority: 45,
		})
	}

	slices.SortStableFunc(recommendations, func(a, b Recommendation) int {
		if a.Priority == b.Priority {
			return strings.Compare(a.ID, b.ID)
		}
		if a.Priority > b.Priority {
			return -1
		}
		return 1
	})

	if len(recommendations) > 3 {
		return recommendations[:3]
	}
	return recommendations
}

func countMarkdownFiles(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		count++
	}
	return count
}

func scanPlans(dir string) (count int, hasUnchecked bool, hasCompleted bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, false, false
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		count++
		content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		text := string(content)
		if strings.Contains(text, "- [ ]") {
			hasUnchecked = true
		}
		if strings.Contains(text, "status: completed") {
			hasCompleted = true
		}
	}
	return count, hasUnchecked, hasCompleted
}

func firstProblemOutcome(outcomes []InstallOutcome) *InstallOutcome {
	for i := range outcomes {
		if outcomes[i].Status == InstallStepFailed || outcomes[i].Status == InstallStepWarning {
			return &outcomes[i]
		}
	}
	return nil
}

func stepUsable(outcomes []InstallOutcome, contains string) bool {
	for _, outcome := range outcomes {
		if !strings.Contains(strings.ToLower(outcome.Step), strings.ToLower(contains)) {
			continue
		}
		return outcome.Status != InstallStepFailed
	}
	return false
}

// WalkMarkdownFiles counts markdown files recursively. Reserved for later launchpad indexing.
func WalkMarkdownFiles(root string) int {
	count := 0
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		count++
		return nil
	})
	return count
}

func countSubdirs(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			count++
		}
	}
	return count
}

func countFilesWithSuffix(dir string, suffix string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), suffix) {
			count++
		}
	}
	return count
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
