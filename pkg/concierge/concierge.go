// Package concierge provides typed tool implementations for the optional
// GitHub Copilot SDK assistant layer. Each tool wraps the deterministic
// local state from installstate and returns structured JSON.
//
// The concierge is explicitly secondary — it explains and navigates the
// local memory index and recommendations but never owns the truth.
package concierge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
)

// MemorySummary is the structured response for the getMemorySummary tool.
type MemorySummary struct {
	Status    string                 `json:"status"`
	Manifest  *ManifestInfo          `json:"manifest,omitempty"`
	RepoState installstate.RepoState `json:"repoState"`
	Message   string                 `json:"message,omitempty"`
}

// ManifestInfo is a subset of manifest state safe for assistant consumption.
type ManifestInfo struct {
	GeneratedAt    string                      `json:"generatedAt"`
	PresetName     string                      `json:"presetName"`
	StackPacks     []installstate.StackPack    `json:"stackPacks"`
	GstackDirs     []string                    `json:"gstackDirs,omitempty"`
	GstackRuntime  bool                        `json:"gstackRuntime"`
	AgentBrowser   bool                        `json:"agentBrowser"`
	OutcomeSummary installstate.OutcomeSummary `json:"outcomeSummary"`
}

// RecommendationList is the structured response for the listRecommendations tool.
type RecommendationList struct {
	Status          string                        `json:"status"`
	Recommendations []installstate.Recommendation `json:"recommendations"`
	Source          string                        `json:"source"`
	Message         string                        `json:"message,omitempty"`
}

// RecommendationDetail is the structured response for the explainRecommendation tool.
type RecommendationDetail struct {
	Status       string `json:"status"`
	ID           string `json:"id"`
	Title        string `json:"title"`
	Reason       string `json:"reason"`
	Priority     int    `json:"priority"`
	SuggestedCmd string `json:"suggestedCommand,omitempty"`
	Message      string `json:"message,omitempty"`
}

// ArtifactInfo is the structured response for the openArtifact tool.
type ArtifactInfo struct {
	Status  string `json:"status"`
	Path    string `json:"path"`
	Exists  bool   `json:"exists"`
	Type    string `json:"type"`
	Message string `json:"message,omitempty"`
}

// ActionResult is the structured response for the runSuggestedAction tool.
type ActionResult struct {
	Status  string `json:"status"`
	Action  string `json:"action"`
	Message string `json:"message"`
}

// GetMemorySummary returns a structured overview of the repo's local memory
// and install intelligence. Gracefully degrades when no manifest exists.
func GetMemorySummary(root string) MemorySummary {
	repoState := installstate.ScanRepoState(root)

	manifest, err := installstate.ReadManifest(root)
	if err != nil {
		return MemorySummary{
			Status:    degradedStatus(err),
			RepoState: repoState,
			Message:   degradedMessage(err),
		}
	}

	return MemorySummary{
		Status:    "ok",
		RepoState: repoState,
		Manifest: &ManifestInfo{
			GeneratedAt:    manifest.GeneratedAt.Format("2006-01-02 15:04 MST"),
			PresetName:     manifest.Requested.PresetName,
			StackPacks:     manifest.Requested.StackPacks,
			GstackDirs:     manifest.Requested.GstackDirs,
			GstackRuntime:  manifest.Requested.GstackRuntime,
			AgentBrowser:   manifest.Requested.IncludeAgentBrowser,
			OutcomeSummary: installstate.SummarizeOutcomes(manifest.Outcomes),
		},
	}
}

// ListRecommendations returns the deterministic next-step recommendations.
// The assistant must not reorder or filter these without explanation.
func ListRecommendations(root string) RecommendationList {
	snapshot, err := installstate.BuildLaunchpadSnapshot(root)
	if err != nil {
		return RecommendationList{
			Status:  "error",
			Source:  "local",
			Message: fmt.Sprintf("failed to build snapshot: %v", err),
		}
	}

	return RecommendationList{
		Status:          "ok",
		Source:          "local-deterministic",
		Recommendations: snapshot.Recommendations,
	}
}

// ExplainRecommendation returns a detailed explanation for a single recommendation
// by its ID, including the suggested CLI command to execute it.
func ExplainRecommendation(root string, id string) RecommendationDetail {
	snapshot, err := installstate.BuildLaunchpadSnapshot(root)
	if err != nil {
		return RecommendationDetail{
			Status:  "error",
			ID:      id,
			Message: fmt.Sprintf("failed to build snapshot: %v", err),
		}
	}

	for _, rec := range snapshot.Recommendations {
		if rec.ID == id {
			return RecommendationDetail{
				Status:       "ok",
				ID:           rec.ID,
				Title:        rec.Title,
				Reason:       rec.Reason,
				Priority:     rec.Priority,
				SuggestedCmd: suggestedCommand(rec.ID),
			}
		}
	}

	return RecommendationDetail{
		Status:  "not-found",
		ID:      id,
		Message: fmt.Sprintf("no recommendation with id %q in the current deterministic set", id),
	}
}

// OpenArtifact resolves a logical artifact name to a filesystem path and checks
// whether it exists. This allows the assistant to help users navigate artifacts
// without scraping the filesystem directly.
func OpenArtifact(root string, artifact string) ArtifactInfo {
	path, artType := resolveArtifact(root, artifact)
	if path == "" {
		return ArtifactInfo{
			Status:  "unknown",
			Path:    artifact,
			Type:    "unknown",
			Message: fmt.Sprintf("unrecognized artifact %q; known artifacts: manifest, instructions, brainstorms, plans, solutions, agents, skills", artifact),
		}
	}

	_, err := os.Stat(path)
	return ArtifactInfo{
		Status: "ok",
		Path:   filepath.ToSlash(path),
		Exists: err == nil,
		Type:   artType,
	}
}

// RunSuggestedAction validates and describes a suggested action without
// executing it. The assistant should present the command for user confirmation,
// never execute silently.
func RunSuggestedAction(root string, actionID string) ActionResult {
	cmd := suggestedCommand(actionID)
	if cmd == "" {
		return ActionResult{
			Status:  "unknown",
			Action:  actionID,
			Message: fmt.Sprintf("no known action for %q; the assistant cannot invent actions that bypass deterministic recommendations", actionID),
		}
	}

	return ActionResult{
		Status:  "ready",
		Action:  actionID,
		Message: fmt.Sprintf("Suggested command: %s — present this to the user for confirmation before running.", cmd),
	}
}

func degradedStatus(err error) string {
	if os.IsNotExist(err) {
		return "no-manifest"
	}
	return "error"
}

func degradedMessage(err error) string {
	if os.IsNotExist(err) {
		return "No guided install manifest found. Run 'atv-installer init --guided' first, or the launchpad will operate with repo memory only."
	}
	return fmt.Sprintf("Failed to read manifest: %v", err)
}

func suggestedCommand(id string) string {
	commands := map[string]string{
		"fix-install-issues":        "atv-installer init --guided",
		"start-brainstorm":          `/ce-brainstorm "your feature idea"`,
		"turn-brainstorm-into-plan": "/ce-plan",
		"execute-active-plan":       "/ce-work",
		"compound-learnings":        "/ce-compound",
		"start-gstack-sprint":       "/gstack-office-hours",
		"browser-check":             "agent-browser open https://yourapp.com",
	}
	return commands[id]
}

func resolveArtifact(root string, artifact string) (string, string) {
	artifact = strings.ToLower(strings.TrimSpace(artifact))
	switch artifact {
	case "manifest":
		return installstate.ManifestPath(root), "file"
	case "instructions":
		return filepath.Join(root, ".github", "copilot-instructions.md"), "file"
	case "brainstorms":
		return filepath.Join(root, "docs", "brainstorms"), "directory"
	case "plans":
		return filepath.Join(root, "docs", "plans"), "directory"
	case "solutions":
		return filepath.Join(root, "docs", "solutions"), "directory"
	case "agents":
		return filepath.Join(root, ".github", "agents"), "directory"
	case "skills":
		return filepath.Join(root, ".github", "skills"), "directory"
	default:
		return "", ""
	}
}
