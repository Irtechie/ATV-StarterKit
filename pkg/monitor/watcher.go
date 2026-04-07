package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
	"github.com/fsnotify/fsnotify"
)

// ignoredPatterns lists filename patterns to filter from FS events.
var ignoredPatterns = []string{
	".swp", "~", ".tmp", ".DS_Store",
	"4913", // Vim temp file
}

// WatcherOptions configures the filesystem watcher.
type WatcherOptions struct {
	DebounceWindow time.Duration // default 300ms
}

// Watcher monitors filesystem events and maintains live state.
type Watcher struct {
	root        string
	fsWatcher   *fsnotify.Watcher
	state       LiveState
	mu          sync.RWMutex
	debounceWin time.Duration
	onChange    func(LiveState) // notify TUI of state changes
	stateFile   string         // .atv/dashboard-state.json
	cancel      context.CancelFunc
	done        chan struct{}
}

// NewWatcher creates a filesystem watcher for the given repo root.
func NewWatcher(root string, opts WatcherOptions) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create fsnotify watcher: %w", err)
	}

	debounce := opts.DebounceWindow
	if debounce == 0 {
		debounce = 300 * time.Millisecond
	}

	w := &Watcher{
		root:        root,
		fsWatcher:   fsw,
		debounceWin: debounce,
		stateFile:   filepath.Join(root, ".atv", "dashboard-state.json"),
		done:        make(chan struct{}),
	}

	return w, nil
}

// SetOnChange registers a callback invoked on every state update.
func (w *Watcher) SetOnChange(fn func(LiveState)) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.onChange = fn
}

// Start begins watching filesystem events and maintaining state.
func (w *Watcher) Start(ctx context.Context) error {
	ctx, w.cancel = context.WithCancel(ctx)

	// Initial full scan
	w.fullScan()

	// Add watch targets (silently skip dirs that don't exist)
	w.addWatchTargets()

	// Start event loop
	go w.eventLoop(ctx)

	return nil
}

// Stop stops the watcher and cleans up resources.
func (w *Watcher) Stop() error {
	if w.cancel != nil {
		w.cancel()
	}
	<-w.done
	return w.fsWatcher.Close()
}

// State returns a threadsafe snapshot of the current live state.
func (w *Watcher) State() LiveState {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.state
}

// Root returns the watched root directory path.
func (w *Watcher) Root() string {
	return w.root
}

// ForceRefresh triggers a full re-scan of all watched directories.
func (w *Watcher) ForceRefresh() {
	w.fullScan()
	w.addWatchTargets() // pick up new dirs
}

// addWatchTargets adds known directories to the fsnotify watcher.
func (w *Watcher) addWatchTargets() {
	targets := []string{
		filepath.Join(w.root, "docs", "brainstorms"),
		filepath.Join(w.root, "docs", "plans"),
		filepath.Join(w.root, "docs", "solutions"),
		filepath.Join(w.root, "docs"),
		filepath.Join(w.root, ".github"),
		filepath.Join(w.root, ".github", "agents"),
		filepath.Join(w.root, ".github", "skills"),
		filepath.Join(w.root, ".github", "prompts"),
		filepath.Join(w.root, ".atv"),
		filepath.Join(w.root, ".gstack"),
		filepath.Join(w.root, ".vscode"),
		filepath.Join(w.root, ".copilot-memory"),
	}

	// User-global paths
	if home, err := os.UserHomeDir(); err == nil {
		targets = append(targets,
			filepath.Join(home, ".gstack"),
			filepath.Join(home, ".agent-browser", "sessions"),
		)
	}

	watched := make([]string, 0, len(targets))
	for _, t := range targets {
		if info, err := os.Stat(t); err == nil && info.IsDir() {
			_ = w.fsWatcher.Add(t)
			watched = append(watched, t)

			// For skills dir, add subdirectories (fsnotify is not recursive)
			if strings.HasSuffix(t, "skills") || strings.HasSuffix(t, ".gstack") {
				w.addSubdirs(t)
			}
		}
	}

	w.mu.Lock()
	w.state.WatchedPaths = watched
	w.mu.Unlock()
}

// addSubdirs adds immediate subdirectories to the watcher.
func (w *Watcher) addSubdirs(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			_ = w.fsWatcher.Add(filepath.Join(dir, entry.Name()))
		}
	}
}

// eventLoop processes fsnotify events with debouncing.
func (w *Watcher) eventLoop(ctx context.Context) {
	defer close(w.done)

	var debounceTimer *time.Timer
	var pendingLayers map[string]bool

	resetDebounce := func(layer string) {
		if pendingLayers == nil {
			pendingLayers = make(map[string]bool)
		}
		pendingLayers[layer] = true

		if debounceTimer != nil {
			debounceTimer.Stop()
		}
		debounceTimer = time.AfterFunc(w.debounceWin, func() {
			w.processPendingLayers(pendingLayers)
			pendingLayers = nil
		})
	}

	for {
		select {
		case <-ctx.Done():
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			return
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}
			if isIgnoredFile(event.Name) {
				continue
			}
			layer := w.classifyEventLayer(event.Name)
			if layer != "" {
				resetDebounce(layer)
			}
		case _, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			// Log errors but continue watching
		}
	}
}

// classifyEventLayer determines which monitoring layer an FS event belongs to.
func (w *Watcher) classifyEventLayer(path string) string {
	rel, err := filepath.Rel(w.root, path)
	if err != nil {
		// Could be a user-global path
		if strings.Contains(path, ".gstack") || strings.Contains(path, ".agent-browser") {
			return "runtime"
		}
		return ""
	}
	rel = filepath.ToSlash(rel)

	switch {
	case strings.HasPrefix(rel, "docs/brainstorms/") ||
		strings.HasPrefix(rel, "docs/plans/") ||
		strings.HasPrefix(rel, "docs/solutions/"):
		return "memory"

	case strings.HasPrefix(rel, ".github/copilot-instructions") ||
		strings.HasPrefix(rel, ".github/skills/") ||
		strings.HasPrefix(rel, ".github/agents/") ||
		strings.HasPrefix(rel, ".github/prompts/") ||
		strings.HasPrefix(rel, ".github/copilot-mcp-config"):
		return "context"

	case strings.HasPrefix(rel, ".atv/"):
		return "health"

	case strings.HasPrefix(rel, ".gstack/") ||
		strings.HasPrefix(rel, ".vscode/"):
		return "health"

	case strings.HasPrefix(rel, ".copilot-memory/"):
		return "memory"

	default:
		return ""
	}
}

// processPendingLayers re-scans the affected monitoring layers.
func (w *Watcher) processPendingLayers(layers map[string]bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if layers["memory"] {
		w.scanMemoryLayer()
	}
	if layers["context"] {
		w.scanContextLayer()
	}
	if layers["health"] {
		w.scanHealthLayer()
	}
	if layers["runtime"] {
		w.scanRuntimeLayer()
	}

	w.state.LastFSEvent = time.Now()

	// Write state file for VS Code extension
	w.writeStateFile()

	// Notify TUI
	if w.onChange != nil {
		stateCopy := w.state
		go w.onChange(stateCopy)
	}
}

// fullScan performs a complete scan of all layers.
func (w *Watcher) fullScan() {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Rebuild base snapshot
	snapshot, _ := installstate.BuildInstallSnapshot(w.root)
	w.state.InstallSnapshot = snapshot
	w.state.SchemaVersion = stateSchemaVersion

	w.scanMemoryLayer()
	w.scanContextLayer()
	w.scanHealthLayer()
	w.scanRuntimeLayer()

	w.state.LastFSEvent = time.Now()
	w.writeStateFile()
}

// scanMemoryLayer scans brainstorms, plans, and solutions.
func (w *Watcher) scanMemoryLayer() {
	w.state.Brainstorms = scanArtifactDir(filepath.Join(w.root, "docs", "brainstorms"))
	w.state.Plans = scanArtifactDir(filepath.Join(w.root, "docs", "plans"))
	w.state.Solutions = scanArtifactDir(filepath.Join(w.root, "docs", "solutions"))

	// Also update snapshot counts
	w.state.InstallSnapshot.RepoState.BrainstormCount = len(w.state.Brainstorms)
	w.state.InstallSnapshot.RepoState.PlanCount = len(w.state.Plans)
	w.state.InstallSnapshot.RepoState.SolutionCount = len(w.state.Solutions)
}

// scanContextLayer scans Copilot context surface (skills, agents, instructions, prompts).
func (w *Watcher) scanContextLayer() {
	githubDir := filepath.Join(w.root, ".github")
	w.state.InstalledSkillNames = installstate.ListSkillDirs(filepath.Join(githubDir, "skills"))
	w.state.InstalledAgentNames = installstate.ListAgentNames(filepath.Join(githubDir, "agents"))

	// Context estimate
	var totalBytes int64
	skillCount := len(w.state.InstalledSkillNames)
	agentCount := len(w.state.InstalledAgentNames)
	promptCount := countFilesInDir(filepath.Join(githubDir, "prompts"), ".prompt.md")
	mcpCount := 0

	// Sum instruction file sizes
	instrDir := githubDir
	entries, err := os.ReadDir(instrDir)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".instructions.md") {
				if info, err := e.Info(); err == nil {
					totalBytes += info.Size()
				}
			}
		}
	}
	// Add copilot-instructions.md
	if info, err := os.Stat(filepath.Join(githubDir, "copilot-instructions.md")); err == nil {
		totalBytes += info.Size()
	}
	// Count MCP servers
	mcpPath := filepath.Join(githubDir, "copilot-mcp-config.json")
	mcpCount = installstate.CountJSONKeyCount(mcpPath, "servers")

	w.state.ContextEstimate = ContextEstimate{
		TotalInstructionBytes: totalBytes,
		SkillCount:            skillCount,
		AgentCount:            agentCount,
		PromptCount:           promptCount,
		MCPServerCount:        mcpCount,
		EstimatedTokens:       int(totalBytes / 4), // rough approximation
	}

	// Update snapshot fields
	w.state.InstallSnapshot.RepoState.InstalledSkills = skillCount
	w.state.InstallSnapshot.RepoState.InstalledAgents = agentCount
	w.state.InstallSnapshot.RepoState.PromptFileCount = promptCount
	w.state.InstallSnapshot.RepoState.MCPServerCount = mcpCount
}

// scanHealthLayer scans install manifest and drift.
func (w *Watcher) scanHealthLayer() {
	// Re-read manifest to update outcomes
	manifest, err := installstate.ReadManifest(w.root)
	if err == nil {
		w.state.InstallSnapshot.HasManifest = true
		w.state.InstallSnapshot.GeneratedAt = manifest.GeneratedAt
		w.state.InstallSnapshot.Requested = manifest.Requested
		w.state.InstallSnapshot.OutcomeSummary = installstate.SummarizeOutcomes(manifest.Outcomes)
		w.state.InstallSnapshot.Recommendations = installstate.BuildRecommendations(w.root, manifest)

		// Compute drift
		w.state.DriftEntries = ComputeDrift(w.root, manifest)
	}

	// Learning pipeline state
	repoState := installstate.ScanRepoState(w.root)
	w.state.InstinctCount = repoState.InstinctCount
	w.state.ObservationCount = repoState.ObservationCount
	w.state.HasObserverHooks = repoState.HasObserverHooks
}

// scanRuntimeLayer probes runtime tool availability.
func (w *Watcher) scanRuntimeLayer() {
	now := time.Now()

	// gstack: check for .gstack/ dir and runtime binary
	gstackAvail := dirExists(filepath.Join(w.root, ".gstack"))
	w.state.RuntimeHealth.Gstack = RuntimeHealthStatus{
		Available: gstackAvail,
		LastCheck: now,
	}
	if gstackAvail && dirExists(filepath.Join(w.root, ".gstack", "browse", "dist")) {
		w.state.RuntimeHealth.Gstack.Detail = "runtime built"
	} else if gstackAvail {
		w.state.RuntimeHealth.Gstack.Detail = "staging only"
	}

	// Agent browser: check for sessions dir
	home, _ := os.UserHomeDir()
	if home != "" {
		abDir := filepath.Join(home, ".agent-browser", "sessions")
		abAvail := dirExists(abDir)
		w.state.RuntimeHealth.AgentBrowser = RuntimeHealthStatus{
			Available: abAvail,
			LastCheck: now,
		}
	}

	// MCP servers: check config exists
	mcpAvail := fileExists(filepath.Join(w.root, ".github", "copilot-mcp-config.json"))
	w.state.RuntimeHealth.MCPServers = RuntimeHealthStatus{
		Available: mcpAvail,
		LastCheck: now,
	}
}

// writeStateFile atomically writes state to .atv/dashboard-state.json for VS Code.
func (w *Watcher) writeStateFile() {
	dir := filepath.Dir(w.stateFile)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}

	data, err := json.MarshalIndent(w.state, "", "  ")
	if err != nil {
		return
	}
	data = append(data, '\n')

	// Atomic write via temp file + rename (matches WriteManifest pattern)
	tmp, err := os.CreateTemp(dir, "dashboard-state-*.json")
	if err != nil {
		return
	}
	tmpPath := tmp.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return
	}
	if err := tmp.Close(); err != nil {
		return
	}
	_ = os.Rename(tmpPath, w.stateFile)
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func scanArtifactDir(dir string) []ArtifactEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var artifacts []ArtifactEntry
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		artifacts = append(artifacts, ArtifactEntry{
			Name:    entry.Name(),
			Path:    filepath.ToSlash(filepath.Join(dir, entry.Name())),
			ModTime: info.ModTime(),
			Size:    info.Size(),
		})
	}
	return artifacts
}

func isIgnoredFile(path string) bool {
	base := filepath.Base(path)
	for _, pat := range ignoredPatterns {
		if strings.HasSuffix(base, pat) || base == pat {
			return true
		}
	}
	// Ignore .git directory events
	cleanPath := filepath.ToSlash(path)
	return strings.Contains(cleanPath, "/.git/")
}

func countFilesInDir(dir, suffix string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), suffix) {
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
