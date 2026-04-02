package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	lpTitleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("178")).Bold(true)
	lpAccentStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true)
	lpSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true)
	lpWarnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	lpFailStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	lpDimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("136"))
	lpKeyStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	lpCountStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
	lpHeaderStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("178")).Bold(true).Underline(true)
)

// LaunchpadTab represents which section is focused.
type LaunchpadTab int

const (
	TabOverview LaunchpadTab = iota
	TabCopilot
	TabCE
	TabGstack
	TabMoves
	tabCount
)

func (t LaunchpadTab) String() string {
	switch t {
	case TabOverview:
		return "Overview"
	case TabCopilot:
		return "Copilot"
	case TabCE:
		return "CE"
	case TabGstack:
		return "Gstack"
	case TabMoves:
		return "Moves"
	}
	return ""
}

type refreshMsg struct{}

// LaunchpadModel is a live Bubble Tea dashboard.
type LaunchpadModel struct {
	root     string
	snapshot installstate.LaunchpadSnapshot
	tab      LaunchpadTab
	err      error

	// Compound engineering
	brainstorms []string
	plans       []string
	solutions   []string

	// Installed intelligence
	agents []string
	skills []string

	// Copilot context
	instructions []string
	prompts      []string
	mcpServers   []string
	extensions   []string

	// gstack
	gstackSkills []string

	// Memory
	memoryFiles []string
}

// NewLaunchpadModel creates the live dashboard model.
func NewLaunchpadModel(root string) LaunchpadModel {
	m := LaunchpadModel{root: root}
	m.refresh()
	return m
}

func (m *LaunchpadModel) refresh() {
	snapshot, err := installstate.BuildLaunchpadSnapshot(m.root)
	m.snapshot = snapshot
	m.err = err
	m.brainstorms = installstate.ListMarkdownNames(filepath.Join(m.root, "docs", "brainstorms"))
	m.plans = installstate.ListMarkdownNames(filepath.Join(m.root, "docs", "plans"))
	m.solutions = installstate.ListMarkdownNames(filepath.Join(m.root, "docs", "solutions"))
	m.agents = installstate.ListAgentNames(filepath.Join(m.root, ".github", "agents"))
	m.skills = installstate.ListSkillDirs(filepath.Join(m.root, ".github", "skills"))
	m.instructions = installstate.ListInstructionFiles(filepath.Join(m.root, ".github"))
	m.prompts = installstate.ListPromptFiles(filepath.Join(m.root, ".github", "prompts"))
	m.gstackSkills = installstate.ListGstackSkillNames(filepath.Join(m.root, ".gstack"))
	m.mcpServers = installstate.ListMCPServerNames(filepath.Join(m.root, ".github", "copilot-mcp-config.json"))
	m.extensions = installstate.ListExtensionRecommendations(filepath.Join(m.root, ".vscode", "extensions.json"))
	m.memoryFiles = installstate.ListMemoryFiles(m.root)
}

func (m LaunchpadModel) Init() tea.Cmd {
	return tea.Batch(tickRefresh())
}

func tickRefresh() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return refreshMsg{}
	})
}

func (m LaunchpadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "right", "l":
			m.tab = (m.tab + 1) % tabCount
		case "shift+tab", "left", "h":
			m.tab = (m.tab - 1 + tabCount) % tabCount
		case "1":
			m.tab = TabOverview
		case "2":
			m.tab = TabCopilot
		case "3":
			m.tab = TabCE
		case "4":
			m.tab = TabGstack
		case "5":
			m.tab = TabMoves
		case "r":
			m.refresh()
		}
	case refreshMsg:
		m.refresh()
		return m, tickRefresh()
	}
	return m, nil
}

func (m LaunchpadModel) View() string {
	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(lpAccentStyle.Render("  ⚡") + lpTitleStyle.Render(" ATV Launchpad ") + lpAccentStyle.Render("⚡"))
	b.WriteString(lpDimStyle.Render("  Live dashboard · auto-refreshes every 3s"))
	b.WriteString("\n\n")

	// Tab bar
	b.WriteString("  ")
	for i := LaunchpadTab(0); i < tabCount; i++ {
		label := fmt.Sprintf(" %d:%s ", i+1, i.String())
		if i == m.tab {
			b.WriteString(lpAccentStyle.Render("[" + label + "]"))
		} else {
			b.WriteString(lpDimStyle.Render(" " + label + " "))
		}
		if i < tabCount-1 {
			b.WriteString(lpDimStyle.Render("│"))
		}
	}
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(lpFailStyle.Render(fmt.Sprintf("  Error: %v\n", m.err)))
		b.WriteString(lpDimStyle.Render("  Press r to retry, q to quit\n"))
		return b.String()
	}

	switch m.tab {
	case TabOverview:
		m.renderOverview(&b)
	case TabCopilot:
		m.renderCopilot(&b)
	case TabCE:
		m.renderCE(&b)
	case TabGstack:
		m.renderGstack(&b)
	case TabMoves:
		m.renderRecommendations(&b)
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(lpDimStyle.Render("  ← → tab  1-5 jump  r refresh  q quit"))
	b.WriteString("\n\n")

	return b.String()
}

// ─── Overview Tab ───────────────────────────────────────────────────────────

func (m LaunchpadModel) renderOverview(b *strings.Builder) {
	s := m.snapshot

	b.WriteString(lpHeaderStyle.Render("  Install Intelligence"))
	b.WriteString("\n\n")

	if s.HasManifest {
		b.WriteString(fmt.Sprintf("  %s Manifest    %s\n", lpSuccessStyle.Render("●"), lpDimStyle.Render(s.ManifestPath)))
		if !s.GeneratedAt.IsZero() {
			b.WriteString(fmt.Sprintf("  %s Last run    %s\n", lpDimStyle.Render("│"), s.GeneratedAt.Format("2006-01-02 15:04 MST")))
		}
		if s.Requested.PresetName != "" {
			b.WriteString(fmt.Sprintf("  %s Preset      %s\n", lpDimStyle.Render("│"), lpCountStyle.Render(s.Requested.PresetName)))
		}
		if labels := s.StackPackLabels(); len(labels) > 0 {
			b.WriteString(fmt.Sprintf("  %s Stacks      %s\n", lpDimStyle.Render("│"), strings.Join(labels, ", ")))
		}
		b.WriteString(fmt.Sprintf("  %s Outcomes    %s done  %s warn  %s fail  %s skip\n",
			lpDimStyle.Render("╰"),
			lpSuccessStyle.Render(fmt.Sprintf("%d", s.OutcomeSummary.Done)),
			lpWarnStyle.Render(fmt.Sprintf("%d", s.OutcomeSummary.Warning)),
			lpFailStyle.Render(fmt.Sprintf("%d", s.OutcomeSummary.Failed)),
			lpDimStyle.Render(fmt.Sprintf("%d", s.OutcomeSummary.Skipped)),
		))
	} else {
		b.WriteString(fmt.Sprintf("  %s No manifest yet. Run %s\n",
			lpWarnStyle.Render("⚠"),
			lpKeyStyle.Render("atv-installer init --guided")))
	}

	b.WriteString("\n")
	b.WriteString(lpHeaderStyle.Render("  Capability Matrix"))
	b.WriteString("\n\n")

	// Row 1: Intelligence surface
	b.WriteString(fmt.Sprintf("  %s agents   %s skills   %s instructions   %s prompts\n",
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.InstalledAgents)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.InstalledSkills)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.InstructionFileCount)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.PromptFileCount)),
	))

	// Row 2: Compound engineering
	b.WriteString(fmt.Sprintf("  %s brainstorms   %s plans   %s solutions\n",
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.BrainstormCount)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.PlanCount)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.SolutionCount)),
	))

	// Row 3: Infrastructure
	b.WriteString(fmt.Sprintf("  %s MCP servers   %s extensions   %s gstack skills   %s memory files\n",
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.MCPServerCount)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.ExtensionRecommendationCount)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.GstackSkillCount)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.MemoryFileCount)),
	))

	b.WriteString("\n")
	b.WriteString(lpHeaderStyle.Render("  Health"))
	b.WriteString("\n\n")

	// Copilot context health
	renderStatus(b, s.RepoState.HasCopilotInstructions, "copilot-instructions.md")
	renderStatus(b, s.RepoState.HasSetupSteps, "copilot-setup-steps.yml")
	renderStatus(b, s.RepoState.HasMCPConfig, "MCP server config")
	renderStatus(b, s.RepoState.HasCELocalConfig, "compound-engineering.local.md")
	renderStatus(b, s.RepoState.HasGstackStaging, ".gstack staging")
	renderStatus(b, s.RepoState.HasGstackRuntime, "gstack runtime (browse)")
	renderStatus(b, s.RepoState.HasAgentBrowserSkill, "agent-browser skill")
	renderStatus(b, s.RepoState.HasGstackUserConfig, "~/.gstack/ user config")
	renderStatus(b, s.RepoState.HasAgentBrowserSessions, "~/.agent-browser/ sessions")
	b.WriteString("\n")

	// Memory files
	if len(m.memoryFiles) > 0 {
		b.WriteString(lpHeaderStyle.Render(fmt.Sprintf("  Memory Files (%d)", len(m.memoryFiles))))
		b.WriteString("\n\n")
		for _, f := range m.memoryFiles {
			b.WriteString(fmt.Sprintf("    • %s\n", f))
		}
		b.WriteString("\n")
	} else {
		b.WriteString(lpHeaderStyle.Render("  Memory Files"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("  %s No .copilot-memory/ files yet\n", lpDimStyle.Render("○")))
		b.WriteString(lpDimStyle.Render("    Copilot stores repo-scoped facts here automatically\n"))
		b.WriteString("\n")
	}

	if s.RepoState.HasUncheckedPlan {
		b.WriteString(fmt.Sprintf("  %s Active plan has unchecked work\n", lpWarnStyle.Render("⚠")))
	} else if s.RepoState.HasCompletedPlan {
		b.WriteString(fmt.Sprintf("  %s Completed plan — ready for /ce-compound\n", lpSuccessStyle.Render("✓")))
	}
}

func renderStatus(b *strings.Builder, ok bool, label string) {
	if ok {
		b.WriteString(fmt.Sprintf("  %s %s\n", lpSuccessStyle.Render("●"), label))
	} else {
		b.WriteString(fmt.Sprintf("  %s %s\n", lpDimStyle.Render("○"), label))
	}
}

// ─── Copilot Tab ────────────────────────────────────────────────────────────

func (m LaunchpadModel) renderCopilot(b *strings.Builder) {
	s := m.snapshot

	// Copilot core config
	b.WriteString(lpHeaderStyle.Render("  Copilot Config"))
	b.WriteString("\n\n")
	renderStatus(b, s.RepoState.HasCopilotInstructions, "copilot-instructions.md")
	renderStatus(b, s.RepoState.HasSetupSteps, "copilot-setup-steps.yml")
	renderStatus(b, s.RepoState.HasMCPConfig, fmt.Sprintf("MCP servers (%d configured)", len(m.mcpServers)))
	b.WriteString("\n")

	// MCP server names (only if any)
	if len(m.mcpServers) > 0 {
		renderBulletList(b, "MCP Servers", m.mcpServers)
		b.WriteString("\n")
	}

	// File-level instructions
	renderBulletList(b, fmt.Sprintf("File Instructions (%d)", len(m.instructions)), m.instructions)
	b.WriteString("\n")

	// Prompts
	renderBulletList(b, fmt.Sprintf("Prompt Files (%d)", len(m.prompts)), m.prompts)
	b.WriteString("\n")

	// Agents — use column layout since there can be many
	b.WriteString(fmt.Sprintf("  %s\n", lpTitleStyle.Render(fmt.Sprintf("Agents (%d)", len(m.agents)))))
	if len(m.agents) == 0 {
		b.WriteString(lpDimStyle.Render("    (none)\n"))
	} else {
		renderColumnList(b, m.agents)
	}
	b.WriteString("\n")

	// VS Code extensions
	renderBulletList(b, fmt.Sprintf("VS Code Extensions (%d recommended)", len(m.extensions)), m.extensions)
}

// ─── Compound Engineering Tab ───────────────────────────────────────────────

func (m LaunchpadModel) renderCE(b *strings.Builder) {
	b.WriteString(lpHeaderStyle.Render("  Compound Engineering Workflow"))
	b.WriteString("\n\n")

	// Workflow status
	stage := m.ceStage()
	b.WriteString(fmt.Sprintf("  Current stage: %s\n\n", lpAccentStyle.Render(stage)))

	renderFileSection(b, "Brainstorms", "docs/brainstorms/", m.brainstorms)
	renderFileSection(b, "Plans", "docs/plans/", m.plans)
	renderFileSection(b, "Solutions", "docs/solutions/", m.solutions)

	// CE project config
	s := m.snapshot
	b.WriteString(lpHeaderStyle.Render("  Project Config"))
	b.WriteString("\n\n")
	if s.RepoState.HasCELocalConfig {
		b.WriteString(fmt.Sprintf("  %s compound-engineering.local.md\n", lpSuccessStyle.Render("●")))
		if s.RepoState.CEReviewAgentCount > 0 {
			b.WriteString(fmt.Sprintf("  %s %s review agents configured\n",
				lpDimStyle.Render("  ╰"),
				lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.CEReviewAgentCount))))
		}
	} else {
		b.WriteString(fmt.Sprintf("  %s No compound-engineering.local.md — defaults used\n", lpDimStyle.Render("○")))
		b.WriteString(lpDimStyle.Render("    Create with /ce-review to configure per-project review agents\n"))
	}
	b.WriteString("\n")

	// Workflow hints
	b.WriteString(lpDimStyle.Render("  Workflow: brainstorm → plan → work → compound"))
	b.WriteString("\n")
	switch stage {
	case "Ready to brainstorm":
		b.WriteString(fmt.Sprintf("  %s Try: %s\n", lpKeyStyle.Render("→"), lpKeyStyle.Render("/ce-brainstorm \"your feature idea\"")))
	case "Ready to plan":
		b.WriteString(fmt.Sprintf("  %s Try: %s\n", lpKeyStyle.Render("→"), lpKeyStyle.Render("/ce-plan")))
	case "Work in progress":
		b.WriteString(fmt.Sprintf("  %s Try: %s\n", lpKeyStyle.Render("→"), lpKeyStyle.Render("/ce-work")))
	case "Ready to compound":
		b.WriteString(fmt.Sprintf("  %s Try: %s\n", lpKeyStyle.Render("→"), lpKeyStyle.Render("/ce-compound")))
	}
}

func (m LaunchpadModel) ceStage() string {
	s := m.snapshot.RepoState
	switch {
	case s.BrainstormCount == 0:
		return "Ready to brainstorm"
	case s.PlanCount == 0:
		return "Ready to plan"
	case s.HasUncheckedPlan:
		return "Work in progress"
	case s.HasCompletedPlan && s.SolutionCount == 0:
		return "Ready to compound"
	case s.SolutionCount > 0:
		return "Compounding knowledge"
	default:
		return "Active"
	}
}

// ─── Gstack Tab ─────────────────────────────────────────────────────────────

func (m LaunchpadModel) renderGstack(b *strings.Builder) {
	s := m.snapshot

	b.WriteString(lpHeaderStyle.Render("  Gstack Status"))
	b.WriteString("\n\n")

	if !s.RepoState.HasGstackStaging {
		b.WriteString(fmt.Sprintf("  %s Gstack not installed. Run with %s\n",
			lpDimStyle.Render("○"),
			lpKeyStyle.Render("atv-installer init --guided")))
		return
	}

	// Runtime status
	if s.RepoState.HasGstackRuntime {
		b.WriteString(fmt.Sprintf("  %s Runtime built (browse binary ready)\n", lpSuccessStyle.Render("●")))
	} else {
		b.WriteString(fmt.Sprintf("  %s Runtime not built — docs-only mode\n", lpWarnStyle.Render("⚠")))
		b.WriteString(lpDimStyle.Render("    Re-run installer with --guided to retry build\n"))
	}

	if s.Requested.GstackRuntime {
		b.WriteString(fmt.Sprintf("  %s Mode: %s\n", lpDimStyle.Render("│"), lpCountStyle.Render("full runtime")))
	} else {
		b.WriteString(fmt.Sprintf("  %s Mode: %s\n", lpDimStyle.Render("│"), lpDimStyle.Render("markdown-only")))
	}

	if s.RepoState.HasAgentBrowserSkill {
		b.WriteString(fmt.Sprintf("  %s agent-browser: %s\n", lpDimStyle.Render("╰"), lpSuccessStyle.Render("installed")))
	} else {
		b.WriteString(fmt.Sprintf("  %s agent-browser: %s\n", lpDimStyle.Render("╰"), lpDimStyle.Render("not installed")))
	}

	// Gstack skills (inside .gstack/)
	b.WriteString("\n")

	// User-global gstack state
	b.WriteString(lpHeaderStyle.Render("  User State"))
	b.WriteString("\n\n")
	if s.RepoState.HasGstackUserConfig {
		b.WriteString(fmt.Sprintf("  %s ~/.gstack/ %s\n",
			lpSuccessStyle.Render("●"),
			lpDimStyle.Render(fmt.Sprintf("(%d session dirs)", s.RepoState.GstackSessionCount))))
	} else {
		b.WriteString(fmt.Sprintf("  %s ~/.gstack/ not found\n", lpDimStyle.Render("○")))
	}
	if s.RepoState.HasAgentBrowserSessions {
		b.WriteString(fmt.Sprintf("  %s ~/.agent-browser/sessions/ %s\n",
			lpSuccessStyle.Render("●"),
			lpDimStyle.Render(fmt.Sprintf("(%d sessions)", s.RepoState.AgentBrowserSessionCount))))
	} else {
		b.WriteString(fmt.Sprintf("  %s ~/.agent-browser/sessions/ not found\n", lpDimStyle.Render("○")))
	}

	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s %s\n", lpTitleStyle.Render("Gstack Skills"), lpDimStyle.Render(fmt.Sprintf("(%d in .gstack/)", len(m.gstackSkills)))))
	if len(m.gstackSkills) == 0 {
		b.WriteString(lpDimStyle.Render("    (none synced)\n"))
	} else {
		renderColumnList(b, m.gstackSkills)
	}

	// .github/skills (core skills)
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s %s\n", lpTitleStyle.Render("Core Skills"), lpDimStyle.Render(fmt.Sprintf("(%d in .github/skills/)", len(m.skills)))))
	if len(m.skills) == 0 {
		b.WriteString(lpDimStyle.Render("    (none installed)\n"))
	} else {
		renderColumnList(b, m.skills)
	}
}

// ─── Moves (Recommendations) Tab ────────────────────────────────────────────

func (m LaunchpadModel) renderRecommendations(b *strings.Builder) {
	b.WriteString(lpHeaderStyle.Render("  Recommended Next Moves"))
	b.WriteString("\n\n")

	recs := m.snapshot.CloneRecommendations()
	if len(recs) == 0 {
		b.WriteString(lpSuccessStyle.Render("  ✓ All clear — no recommended actions."))
		b.WriteString("\n")
		return
	}

	for i, rec := range recs {
		priority := lpDimStyle.Render(fmt.Sprintf("P%d", rec.Priority))
		b.WriteString(fmt.Sprintf("  %s %s  %s\n",
			lpCountStyle.Render(fmt.Sprintf("%d.", i+1)),
			rec.Title,
			priority,
		))
		b.WriteString(fmt.Sprintf("     %s\n", lpDimStyle.Render(rec.Reason)))
		b.WriteString("\n")
	}
}

// ─── Shared helpers ─────────────────────────────────────────────────────────

func renderFileSection(b *strings.Builder, title, dir string, files []string) {
	count := len(files)
	b.WriteString(fmt.Sprintf("  %s\n", lpTitleStyle.Render(fmt.Sprintf("%s (%d in %s)", title, count, dir))))
	if count == 0 {
		b.WriteString(lpDimStyle.Render("    (empty)\n"))
	} else {
		for _, f := range files {
			b.WriteString(fmt.Sprintf("    %s %s\n", lpDimStyle.Render("•"), f))
		}
	}
	b.WriteString("\n")
}

func renderBulletList(b *strings.Builder, title string, items []string) {
	b.WriteString(fmt.Sprintf("  %s\n", lpTitleStyle.Render(title)))
	if len(items) == 0 {
		b.WriteString(lpDimStyle.Render("    (none)\n"))
	} else {
		for _, item := range items {
			b.WriteString(fmt.Sprintf("    %s %s\n", lpDimStyle.Render("•"), item))
		}
	}
}

func renderColumnList(b *strings.Builder, items []string) {
	cols := 3
	if len(items) < 6 {
		cols = 1
	} else if len(items) < 15 {
		cols = 2
	}
	colWidth := 28
	for i, item := range items {
		if i%cols == 0 {
			b.WriteString("    ")
		}
		// Pad raw item name (not styled text) for correct alignment
		padded := item
		for len(padded) < colWidth && cols > 1 && i%cols < cols-1 {
			padded += " "
		}
		b.WriteString(fmt.Sprintf("%s %s", lpDimStyle.Render("•"), padded))
		if i%cols == cols-1 || i == len(items)-1 {
			b.WriteString("\n")
		}
	}
}

// RunLaunchpad starts the live launchpad TUI.
func RunLaunchpad(root string) error {
	m := NewLaunchpadModel(root)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
