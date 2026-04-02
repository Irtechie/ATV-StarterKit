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
	TabMemory
	TabAgents
	TabSkills
	TabRecommendations
	tabCount
)

func (t LaunchpadTab) String() string {
	switch t {
	case TabOverview:
		return "Overview"
	case TabMemory:
		return "Memory"
	case TabAgents:
		return "Agents"
	case TabSkills:
		return "Skills"
	case TabRecommendations:
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

	// Detailed file lists
	brainstorms []string
	plans       []string
	solutions   []string
	agents      []string
	skills      []string
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
			m.tab = TabMemory
		case "3":
			m.tab = TabAgents
		case "4":
			m.tab = TabSkills
		case "5":
			m.tab = TabRecommendations
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
	case TabMemory:
		m.renderMemory(&b)
	case TabAgents:
		m.renderList(&b, "Installed Agents", m.agents, ".agent.md")
	case TabSkills:
		m.renderList(&b, "Installed Skills", m.skills, "SKILL.md")
	case TabRecommendations:
		m.renderRecommendations(&b)
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(lpDimStyle.Render("  ← → tab  1-5 jump  r refresh  q quit"))
	b.WriteString("\n\n")

	return b.String()
}

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
			lpDimStyle.Render("│"),
			lpSuccessStyle.Render(fmt.Sprintf("%d", s.OutcomeSummary.Done)),
			lpWarnStyle.Render(fmt.Sprintf("%d", s.OutcomeSummary.Warning)),
			lpFailStyle.Render(fmt.Sprintf("%d", s.OutcomeSummary.Failed)),
			lpDimStyle.Render(fmt.Sprintf("%d", s.OutcomeSummary.Skipped)),
		))

		// gstack + agent-browser status
		if s.HasGstack() {
			mode := "markdown-only"
			if s.Requested.GstackRuntime {
				mode = "runtime"
			}
			b.WriteString(fmt.Sprintf("  %s gstack      %d dirs (%s)\n", lpDimStyle.Render("│"), len(s.Requested.GstackDirs), mode))
		}
		if s.HasAgentBrowser() {
			icon := lpSuccessStyle.Render("✓")
			if !s.RepoState.HasAgentBrowserSkill {
				icon = lpWarnStyle.Render("!")
			}
			b.WriteString(fmt.Sprintf("  %s browser     %s agent-browser\n", lpDimStyle.Render("╰"), icon))
		} else {
			b.WriteString(fmt.Sprintf("  %s\n", lpDimStyle.Render("╰")))
		}
	} else {
		b.WriteString(fmt.Sprintf("  %s No manifest yet. Run %s\n",
			lpWarnStyle.Render("⚠"),
			lpKeyStyle.Render("atv-installer init --guided")))
	}

	b.WriteString("\n")
	b.WriteString(lpHeaderStyle.Render("  Quick Stats"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("  %s agents   %s skills   %s brainstorms   %s plans   %s solutions\n",
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.InstalledAgents)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.InstalledSkills)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.BrainstormCount)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.PlanCount)),
		lpCountStyle.Render(fmt.Sprintf("%d", s.RepoState.SolutionCount)),
	))

	// Status indicators
	b.WriteString("  ")
	if s.RepoState.HasCopilotInstructions {
		b.WriteString(lpSuccessStyle.Render("● instructions  "))
	} else {
		b.WriteString(lpDimStyle.Render("○ instructions  "))
	}
	if s.RepoState.HasGstackStaging {
		b.WriteString(lpSuccessStyle.Render("● .gstack  "))
	} else {
		b.WriteString(lpDimStyle.Render("○ .gstack  "))
	}
	if s.RepoState.HasAgentBrowserSkill {
		b.WriteString(lpSuccessStyle.Render("● agent-browser"))
	} else {
		b.WriteString(lpDimStyle.Render("○ agent-browser"))
	}
	b.WriteString("\n")

	if s.RepoState.HasUncheckedPlan {
		b.WriteString(fmt.Sprintf("\n  %s Active plan has unchecked work\n", lpWarnStyle.Render("⚠")))
	} else if s.RepoState.HasCompletedPlan {
		b.WriteString(fmt.Sprintf("\n  %s Completed plan — ready for /ce-compound\n", lpSuccessStyle.Render("✓")))
	}
}

func (m LaunchpadModel) renderMemory(b *strings.Builder) {
	b.WriteString(lpHeaderStyle.Render("  Repo Memory — Detailed View"))
	b.WriteString("\n\n")

	renderFileSection(b, "Brainstorms", "docs/brainstorms/", m.brainstorms)
	renderFileSection(b, "Plans", "docs/plans/", m.plans)
	renderFileSection(b, "Solutions", "docs/solutions/", m.solutions)
}

func renderFileSection(b *strings.Builder, title, dir string, files []string) {
	count := len(files)
	b.WriteString(fmt.Sprintf("  %s %s\n", lpTitleStyle.Render(title), lpDimStyle.Render(fmt.Sprintf("(%d in %s)", count, dir))))
	if count == 0 {
		b.WriteString(lpDimStyle.Render("    (empty)\n"))
	} else {
		for i, f := range files {
			prefix := "├"
			if i == len(files)-1 {
				prefix = "╰"
			}
			b.WriteString(fmt.Sprintf("  %s %s\n", lpDimStyle.Render("  "+prefix), f))
		}
	}
	b.WriteString("\n")
}

func (m LaunchpadModel) renderList(b *strings.Builder, title string, items []string, suffix string) {
	b.WriteString(lpHeaderStyle.Render(fmt.Sprintf("  %s", title)))
	b.WriteString(lpDimStyle.Render(fmt.Sprintf("  (%d installed)", len(items))))
	b.WriteString("\n\n")

	if len(items) == 0 {
		b.WriteString(lpDimStyle.Render("    (none installed)\n"))
		return
	}

	// Render in columns if many items
	cols := 3
	if len(items) < 6 {
		cols = 1
	} else if len(items) < 15 {
		cols = 2
	}

	colWidth := 30
	for i, item := range items {
		if i%cols == 0 {
			b.WriteString("  ")
		}
		label := fmt.Sprintf("  %s %s", lpDimStyle.Render("•"), item)
		// Pad to column width
		for len(label) < colWidth && cols > 1 && i%cols < cols-1 {
			label += " "
		}
		b.WriteString(label)
		if i%cols == cols-1 || i == len(items)-1 {
			b.WriteString("\n")
		}
	}
}

func (m LaunchpadModel) renderRecommendations(b *strings.Builder) {
	b.WriteString(lpHeaderStyle.Render("  Recommended Next Moves"))
	b.WriteString("\n\n")

	recs := m.snapshot.CloneRecommendations()
	if len(recs) == 0 {
		b.WriteString(lpDimStyle.Render("    No recommendations yet. Add repo memory or run guided install.\n"))
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

// RunLaunchpad starts the live launchpad TUI.
func RunLaunchpad(root string) error {
	m := NewLaunchpadModel(root)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
