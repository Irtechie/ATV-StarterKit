package installstate

import (
	"os"
	"path/filepath"
	"slices"
	"time"
)

// OutcomeSummary aggregates structured install outcomes for dashboard rendering.
type OutcomeSummary struct {
	Done    int
	Warning int
	Failed  int
	Skipped int
}

// LaunchpadSnapshot is the deterministic local state consumed by the reopen command.
type LaunchpadSnapshot struct {
	Root            string
	ManifestPath    string
	HasManifest     bool
	GeneratedAt     time.Time
	Requested       RequestedState
	RepoState       RepoState
	OutcomeSummary  OutcomeSummary
	Recommendations []Recommendation
}

// BuildLaunchpadSnapshot assembles the local launchpad view model from repo state and the manifest.
func BuildLaunchpadSnapshot(root string) (LaunchpadSnapshot, error) {
	snapshot := LaunchpadSnapshot{
		Root:         root,
		ManifestPath: filepath.ToSlash(ManifestPath(root)),
		RepoState:    ScanRepoState(root),
	}

	manifest, err := ReadManifest(root)
	if err == nil {
		snapshot.HasManifest = true
		snapshot.GeneratedAt = manifest.GeneratedAt
		snapshot.Requested = manifest.Requested
		snapshot.OutcomeSummary = SummarizeOutcomes(manifest.Outcomes)
		snapshot.Recommendations = BuildRecommendations(root, manifest)
		return snapshot, nil
	}
	if !os.IsNotExist(err) {
		return LaunchpadSnapshot{}, err
	}

	snapshot.Recommendations = BuildRecommendations(root, InstallManifest{})
	return snapshot, nil
}

// SummarizeOutcomes counts outcome statuses in deterministic buckets.
func SummarizeOutcomes(outcomes []InstallOutcome) OutcomeSummary {
	var summary OutcomeSummary
	for _, outcome := range outcomes {
		switch outcome.Status {
		case InstallStepWarning:
			summary.Warning++
		case InstallStepFailed:
			summary.Failed++
		case InstallStepSkipped:
			summary.Skipped++
		default:
			summary.Done++
		}
	}
	return summary
}

// StackPackLabels returns deterministic human-readable labels for selected stack packs.
func (s LaunchpadSnapshot) StackPackLabels() []string {
	labels := make([]string, 0, len(s.Requested.StackPacks))
	for _, pack := range s.Requested.StackPacks {
		labels = append(labels, stackPackLabel(pack))
	}
	return labels
}

// HasGstack reports whether the manifest requested any gstack skills.
func (s LaunchpadSnapshot) HasGstack() bool {
	return len(s.Requested.GstackDirs) > 0
}

// HasAgentBrowser reports whether the manifest requested agent-browser.
func (s LaunchpadSnapshot) HasAgentBrowser() bool {
	return s.Requested.IncludeAgentBrowser
}

// CloneRecommendations returns a copy of recommendations for safe rendering.
func (s LaunchpadSnapshot) CloneRecommendations() []Recommendation {
	return slices.Clone(s.Recommendations)
}

func stackPackLabel(pack StackPack) string {
	switch pack {
	case StackPackRails:
		return "Rails"
	case StackPackPython:
		return "Python"
	case StackPackTypeScript:
		return "TypeScript"
	default:
		return "General"
	}
}

// ListBrainstorms returns the names of brainstorm files in docs/brainstorms/.
func ListMarkdownNames(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			names = append(names, entry.Name())
		}
	}
	return names
}

// ListAgentNames returns the names of installed agent files (.agent.md).
func ListAgentNames(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			name := entry.Name()
			// Strip .agent.md suffix for cleaner display
			if len(name) > 9 && name[len(name)-9:] == ".agent.md" {
				name = name[:len(name)-9]
			}
			names = append(names, name)
		}
	}
	return names
}

// ListSkillDirs returns the names of installed skill directories.
func ListSkillDirs(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names
}
