package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/monitor"
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

// LaunchpadTab represents which signal panel is focused.
type LaunchpadTab int

const (
	TabMemory LaunchpadTab = iota
	TabContext
	TabHealth
	TabMoves
	tabCount
)

func (t LaunchpadTab) String() string {
	switch t {
	case TabMemory:
		return "Memory"
	case TabContext:
		return "Context"
	case TabHealth:
		return "Health"
	case TabMoves:
		return "Moves"
	}
	return ""
}

// stateChangedMsg is sent when the watcher detects filesystem changes.
type stateChangedMsg struct{}

// LaunchpadModel is a live Bubble Tea dashboard backed by the filesystem watcher.
type LaunchpadModel struct {
	root    string
	watcher *monitor.Watcher
	tab     LaunchpadTab

	// Suggest-then-execute state
	selectedRec int
}

// NewLaunchpadModel creates the live dashboard model backed by a watcher.
func NewLaunchpadModel(root string, w *monitor.Watcher) LaunchpadModel {
	return LaunchpadModel{
		root:    root,
		watcher: w,
	}
}

func (m LaunchpadModel) Init() tea.Cmd {
	// Bridge watcher onChange callback to tea.Msg
	ch := make(chan struct{}, 1)
	m.watcher.SetOnChange(func(_ monitor.LiveState) {
		select {
		case ch <- struct{}{}:
		default:
		}
	})

	// Wait for watcher events and convert to tea.Msg
	return func() tea.Msg {
		<-ch
		return stateChangedMsg{}
	}
}

func waitForChange(ch chan struct{}) tea.Cmd {
	return func() tea.Msg {
		<-ch
		return stateChangedMsg{}
	}
}

func (m LaunchpadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab", "right", "l":
			m.tab = (m.tab + 1) % tabCount
		case "shift+tab", "left", "h":
			m.tab = (m.tab - 1 + tabCount) % tabCount
		case "1":
			m.tab = TabMemory
		case "2":
			m.tab = TabContext
		case "3":
			m.tab = TabHealth
		case "4":
			m.tab = TabMoves
		case "j", "down":
			state := m.watcher.State()
			if m.tab == TabMoves && m.selectedRec < len(state.LaunchpadSnapshot.Recommendations)-1 {
				m.selectedRec++
			}
		case "k", "up":
			if m.tab == TabMoves && m.selectedRec > 0 {
				m.selectedRec--
			}
		case "r":
			m.watcher.ForceRefresh()
		}
	case stateChangedMsg:
		// Re-subscribe for next change
		ch := make(chan struct{}, 1)
		m.watcher.SetOnChange(func(_ monitor.LiveState) {
			select {
			case ch <- struct{}{}:
			default:
			}
		})
		return m, waitForChange(ch)
	}
	return m, nil
}

func (m LaunchpadModel) View() string {
	var b strings.Builder

	state := m.watcher.State()

	// Header with online/offline indicator
	b.WriteString("\n")
	b.WriteString(lpAccentStyle.Render("  ⚡") + lpTitleStyle.Render(" ATV Launchpad ") + lpAccentStyle.Render("⚡"))
	b.WriteString(lpDimStyle.Render("  Live dashboard · event-driven"))
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

	switch m.tab {
	case TabMemory:
		m.renderMemory(&b, &state)
	case TabContext:
		m.renderContext(&b, &state)
	case TabHealth:
		m.renderHealth(&b, &state)
	case TabMoves:
		m.renderMoves(&b, &state)
	}

	// Footer
	b.WriteString("\n")
	lastEvent := state.LastFSEvent
	if lastEvent.IsZero() {
		b.WriteString(lpDimStyle.Render("  Last FS event: never"))
	} else {
		ago := time.Since(lastEvent).Truncate(time.Second)
		b.WriteString(lpDimStyle.Render(fmt.Sprintf("  Last FS event: %s ago", ago)))
	}
	b.WriteString(lpDimStyle.Render("  │  "))
	b.WriteString(lpDimStyle.Render("← → tab  1-4 jump  j/k navigate  r refresh  q quit"))
	b.WriteString("\n\n")

	return b.String()
}

// ─── Memory Tab ─────────────────────────────────────────────────────────────

func (m LaunchpadModel) renderMemory(b *strings.Builder, state *monitor.LiveState) {
	b.WriteString(lpHeaderStyle.Render("  Repo Memory Artifacts"))
	b.WriteString("\n\n")

	renderArtifactSection(b, "Brainstorms", state.Brainstorms)
	renderArtifactSection(b, "Plans", state.Plans)
	renderArtifactSection(b, "Solutions", state.Solutions)

	// Memory files
	if state.LaunchpadSnapshot.RepoState.MemoryFileCount > 0 {
		b.WriteString(lpHeaderStyle.Render(fmt.Sprintf("  Copilot Memory Files (%d)", state.LaunchpadSnapshot.RepoState.MemoryFileCount)))
		b.WriteString("\n\n")
	} else {
		b.WriteString(lpHeaderStyle.Render("  Copilot Memory"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("  %s No .copilot-memory/ files yet\n", lpDimStyle.Render("○")))
		b.WriteString(lpDimStyle.Render("    Copilot stores repo-scoped facts here automatically\n"))
		b.WriteString("\n")
	}

	// CE workflow stage
	s := state.LaunchpadSnapshot.RepoState
	if s.HasUncheckedPlan {
		b.WriteString(fmt.Sprintf("  %s Active plan has unchecked work\n", lpWarnStyle.Render("⚠")))
	} else if s.HasCompletedPlan {
		b.WriteString(fmt.Sprintf("  %s Completed plan — ready for /ce-compound\n", lpSuccessStyle.Render("✓")))
	}
}

func renderArtifactSection(b *strings.Builder, title string, artifacts []monitor.ArtifactEntry) {
	b.WriteString(fmt.Sprintf("  %s (%d)\n", lpTitleStyle.Render(title), len(artifacts)))
	if len(artifacts) == 0 {
		b.WriteString(lpDimStyle.Render("    (empty)\n"))
	} else {
		for _, a := range artifacts {
			age := time.Since(a.ModTime).Truncate(time.Minute)
			ageStr := formatAge(age)
			b.WriteString(fmt.Sprintf("    %s %s  %s\n",
				lpDimStyle.Render("•"),
				a.Name,
				lpDimStyle.Render(ageStr),
			))
		}
	}
	b.WriteString("\n")
}

func formatAge(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
	return fmt.Sprintf("%dd ago", int(d.Hours()/24))
}

// ─── Context Tab ────────────────────────────────────────────────────────────

func (m LaunchpadModel) renderContext(b *strings.Builder, state *monitor.LiveState) {
	s := state.LaunchpadSnapshot

	b.WriteString(lpHeaderStyle.Render("  Context Estimate"))
	b.WriteString("\n\n")

	ctx := state.ContextEstimate
	b.WriteString(fmt.Sprintf("  Instruction bytes  %s\n", lpCountStyle.Render(fmt.Sprintf("%d", ctx.TotalInstructionBytes))))
	b.WriteString(fmt.Sprintf("  Estimated tokens   %s\n", lpCountStyle.Render(fmt.Sprintf("~%d", ctx.EstimatedTokens))))
	b.WriteString("\n")

	b.WriteString(lpHeaderStyle.Render("  Capability Matrix"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("  %s agents   %s skills   %s instructions   %s prompts\n",
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.InstalledAgents)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.InstalledSkills)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.InstructionFileCount)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.PromptFileCount)),
	))
	b.WriteString(fmt.Sprintf("  %s MCP servers   %s extensions   %s gstack skills\n",
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.MCPServerCount)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.ExtensionRecommendationCount)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.GstackSkillCount)),
	))
	b.WriteString("\n")

	b.WriteString(lpHeaderStyle.Render("  Copilot Config"))
	b.WriteString("\n\n")
	renderStatus(b, s.RepoState.HasCopilotInstructions, "copilot-instructions.md")
	renderStatus(b, s.RepoState.HasSetupSteps, "copilot-setup-steps.yml")
	renderStatus(b, s.RepoState.HasMCPConfig, fmt.Sprintf("MCP servers (%d configured)", s.RepoState.MCPServerCount))
	renderStatus(b, s.RepoState.HasCELocalConfig, "compound-engineering.local.md")
}

func renderStatus(b *strings.Builder, ok bool, label string) {
	if ok {
		b.WriteString(fmt.Sprintf("  %s %s\n", lpSuccessStyle.Render("●"), label))
	} else {
		b.WriteString(fmt.Sprintf("  %s %s\n", lpDimStyle.Render("○"), label))
	}
}

// ─── Health Tab ─────────────────────────────────────────────────────────────

func (m LaunchpadModel) renderHealth(b *strings.Builder, state *monitor.LiveState) {
	s := state.LaunchpadSnapshot

	b.WriteString(lpHeaderStyle.Render("  Install Intelligence"))
	b.WriteString("\n\n")

	if s.HasManifest {
		b.WriteString(fmt.Sprintf("  %s Manifest    %s\n", lpSuccessStyle.Render("●"), lpDimStyle.Render(s.ManifestPath)))
		if !s.GeneratedAt.IsZero() {
			b.WriteString(fmt.Sprintf("  %s Last run    %s\n", lpDimStyle.Render("│"), s.GeneratedAt.Format("2006-01-02 15:04 MST")))
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

	// Drift entries
	b.WriteString("\n")
	b.WriteString(lpHeaderStyle.Render("  Install Drift"))
	b.WriteString("\n\n")

	if len(state.DriftEntries) == 0 {
		b.WriteString(fmt.Sprintf("  %s No drift detected\n", lpSuccessStyle.Render("✓")))
	} else {
		for _, d := range state.DriftEntries {
			icon := lpWarnStyle.Render("⚠")
			status := string(d.Status)
			if d.Status == monitor.DriftMissing {
				icon = lpFailStyle.Render("✗")
			}
			b.WriteString(fmt.Sprintf("  %s %s  %s\n", icon, d.Path, lpDimStyle.Render(status)))
		}
	}

	// Runtime health
	b.WriteString("\n")
	b.WriteString(lpHeaderStyle.Render("  Runtime"))
	b.WriteString("\n\n")
	renderStatus(b, s.RepoState.HasGstackStaging, "gstack staging")
	renderStatus(b, s.RepoState.HasGstackRuntime, "gstack runtime")
	renderStatus(b, s.RepoState.HasAgentBrowserSkill, "agent-browser skill")
	renderStatus(b, s.RepoState.HasGstackUserConfig, "~/.gstack/ user config")
	renderStatus(b, s.RepoState.HasAgentBrowserSessions, "~/.agent-browser/ sessions")
}

// ─── Moves Tab ──────────────────────────────────────────────────────────────

func (m LaunchpadModel) renderMoves(b *strings.Builder, state *monitor.LiveState) {
	b.WriteString(lpHeaderStyle.Render("  Recommended Next Moves"))
	b.WriteString("\n\n")

	recs := state.LaunchpadSnapshot.CloneRecommendations()
	if len(recs) == 0 {
		b.WriteString(lpSuccessStyle.Render("  ✓ All clear — no recommended actions."))
		b.WriteString("\n")
		return
	}

	for i, rec := range recs {
		prefix := "  "
		if i == m.selectedRec {
			prefix = lpAccentStyle.Render("▸ ")
		}
		priority := lpDimStyle.Render(fmt.Sprintf("P%d", rec.Priority))
		b.WriteString(fmt.Sprintf("%s%s %s  %s\n",
			prefix,
			lpCountStyle.Render(fmt.Sprintf("%d.", i+1)),
			rec.Title,
			priority,
		))
		b.WriteString(fmt.Sprintf("     %s\n", lpDimStyle.Render(rec.Reason)))
		b.WriteString("\n")
	}
}

// ─── Shared helpers ─────────────────────────────────────────────────────────

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
func RunLaunchpad(root string, w *monitor.Watcher) error {
	m := NewLaunchpadModel(root, w)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
