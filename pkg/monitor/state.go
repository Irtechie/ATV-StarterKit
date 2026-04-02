package monitor

import (
	"time"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
)

// DriftStatus classifies the relationship between an installed file and its catalog entry.
type DriftStatus string

const (
	DriftMissing      DriftStatus = "missing"       // in manifest but not on disk
	DriftUserModified DriftStatus = "user-modified"  // on disk differs from install-time hash
)

// ActionType classifies the kind of executable action.
type ActionType string

const (
	ActionSlashCommand ActionType = "slash-command"
	ActionInstallerOp  ActionType = "installer-op"
	ActionFileOp       ActionType = "file-op"
	ActionGitOp        ActionType = "git-op"
)

// ActionRisk classifies the risk level of an executable action.
type ActionRisk string

const (
	ActionRiskSafe        ActionRisk = "safe"
	ActionRiskCaution     ActionRisk = "caution"
	ActionRiskDestructive ActionRisk = "destructive"
)

// ArtifactEntry represents a repo memory artifact (brainstorm, plan, solution).
type ArtifactEntry struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	ModTime time.Time `json:"modTime"`
	Size    int64     `json:"size"`
}

// DriftEntry represents a file-level drift classification.
type DriftEntry struct {
	Path        string      `json:"path"`
	Status      DriftStatus `json:"status"`
	DiskHash    string      `json:"diskHash,omitempty"`
	InstallHash string      `json:"installHash,omitempty"`
}

// ContextEstimate approximates Copilot context consumption from local signals.
type ContextEstimate struct {
	TotalInstructionBytes int64 `json:"totalInstructionBytes"`
	SkillCount            int   `json:"skillCount"`
	AgentCount            int   `json:"agentCount"`
	PromptCount           int   `json:"promptCount"`
	MCPServerCount        int   `json:"mcpServerCount"`
	EstimatedTokens       int   `json:"estimatedTokens"`
}

// SDKRecommendation is a context-aware recommendation from the GitHub Copilot SDK.
type SDKRecommendation struct {
	ID             string          `json:"id"`
	Title          string          `json:"title"`
	Reason         string          `json:"reason"`
	Explanation    string          `json:"explanation"`
	ProposedAction *ProposedAction `json:"proposedAction,omitempty"`
	Priority       int             `json:"priority"`
}

// ProposedAction is an executable action proposed by the SDK.
type ProposedAction struct {
	Type        ActionType `json:"type"`
	Command     string     `json:"command"`
	Args        []string   `json:"args,omitempty"`
	Description string     `json:"description"`
	RiskLevel   ActionRisk `json:"riskLevel"`
}

// ActionResult captures the outcome of an executed action.
type ActionResult struct {
	Action     ProposedAction `json:"action"`
	Success    bool           `json:"success"`
	Output     string         `json:"output,omitempty"`
	Error      string         `json:"error,omitempty"`
	ExecutedAt time.Time      `json:"executedAt"`
}

// RuntimeHealthStatus represents the health of a single runtime dependency.
type RuntimeHealthStatus struct {
	Available bool      `json:"available"`
	LastCheck time.Time `json:"lastCheck"`
	Detail    string    `json:"detail,omitempty"`
}

// RuntimeHealthState aggregates health of all runtime dependencies.
type RuntimeHealthState struct {
	Gstack       RuntimeHealthStatus `json:"gstack"`
	AgentBrowser RuntimeHealthStatus `json:"agentBrowser"`
	MCPServers   RuntimeHealthStatus `json:"mcpServers"`
}

// LiveState extends LaunchpadSnapshot with realtime monitoring data.
type LiveState struct {
	installstate.LaunchpadSnapshot

	// Layer 1: Repo memory (live-tracked)
	Brainstorms []ArtifactEntry `json:"brainstorms"`
	Plans       []ArtifactEntry `json:"plans"`
	Solutions   []ArtifactEntry `json:"solutions"`

	// Layer 2: Context approximation
	ContextEstimate ContextEstimate `json:"contextEstimate"`

	// Layer 3: Skills/agents (filesystem approximation for TUI)
	InstalledSkillNames []string `json:"installedSkillNames"`
	InstalledAgentNames []string `json:"installedAgentNames"`

	// Layer 4: Install health
	DriftEntries []DriftEntry `json:"driftEntries"`

	// Layer 5: Runtime health
	RuntimeHealth RuntimeHealthState `json:"runtimeHealth"`

	// SDK intelligence (zero values when offline)
	SDKRecommendations []SDKRecommendation `json:"sdkRecommendations,omitempty"`
	SDKOnline          bool                `json:"sdkOnline"`
	SDKLastQuery       time.Time           `json:"sdkLastQuery,omitempty"`

	// Metadata
	SchemaVersion int       `json:"schemaVersion"`
	LastFSEvent   time.Time `json:"lastFSEvent"`
	WatchedPaths  []string  `json:"watchedPaths"`
}

const stateSchemaVersion = 1
