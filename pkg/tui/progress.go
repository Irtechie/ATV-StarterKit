package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StepStatus represents the state of an install step.
type StepStatus int

const (
	StepPending StepStatus = iota
	StepRunning
	StepDone
	StepFailed
	StepSkipped
)

// InstallStep defines a single step in the install process.
type InstallStep struct {
	Name   string
	Detail string
	Status StepStatus
	Error  error
	Action func() error // the actual work to perform
}

// stepCompleteMsg is sent when a step finishes.
type stepCompleteMsg struct {
	index int
	err   error
}

// ProgressModel is a Bubbletea model for showing install progress.
type ProgressModel struct {
	steps      []InstallStep
	current    int
	spinner    spinner.Model
	done       bool
	presetName string
	stackName  string
}

// NewProgressModel creates a progress model with the given steps.
func NewProgressModel(steps []InstallStep, presetName, stackName string) ProgressModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	return ProgressModel{
		steps:      steps,
		current:    0,
		spinner:    s,
		presetName: presetName,
		stackName:  stackName,
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.runCurrentStep())
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

	case stepCompleteMsg:
		if msg.err != nil {
			m.steps[msg.index].Status = StepFailed
			m.steps[msg.index].Error = msg.err
		} else {
			m.steps[msg.index].Status = StepDone
		}

		// Advance to next step
		m.current = msg.index + 1
		if m.current >= len(m.steps) {
			m.done = true
			return m, tea.Quit
		}

		// Skip steps with nil actions
		for m.current < len(m.steps) && m.steps[m.current].Action == nil {
			m.steps[m.current].Status = StepSkipped
			m.current++
		}
		if m.current >= len(m.steps) {
			m.done = true
			return m, tea.Quit
		}

		m.steps[m.current].Status = StepRunning
		return m, m.runCurrentStep()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m ProgressModel) View() string {
	var b strings.Builder

	// Header
	b.WriteString(fmt.Sprintf("\n  Installing %s preset for %s...\n\n", m.presetName, m.stackName))

	// Steps
	for _, step := range m.steps {
		icon := m.stepIcon(step)
		label := step.Name
		if step.Detail != "" && step.Status == StepDone {
			label = fmt.Sprintf("%s (%s)", step.Name, step.Detail)
		}
		if step.Status == StepFailed && step.Error != nil {
			label = fmt.Sprintf("%s — %v", step.Name, step.Error)
		}
		b.WriteString(fmt.Sprintf("  %s %s\n", icon, label))
	}

	if m.done {
		b.WriteString("\n")
	}

	return b.String()
}

func (m ProgressModel) stepIcon(step InstallStep) string {
	doneStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	failStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	skipStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("243"))

	switch step.Status {
	case StepDone:
		return doneStyle.Render("✅")
	case StepFailed:
		return failStyle.Render("❌")
	case StepRunning:
		return m.spinner.View()
	case StepSkipped:
		return skipStyle.Render("⏭️ ")
	default:
		return dimStyle.Render("○ ")
	}
}

func (m ProgressModel) runCurrentStep() tea.Cmd {
	idx := m.current
	if idx >= len(m.steps) {
		return nil
	}
	step := m.steps[idx]
	m.steps[idx].Status = StepRunning

	return func() tea.Msg {
		var err error
		if step.Action != nil {
			err = step.Action()
		}
		return stepCompleteMsg{index: idx, err: err}
	}
}

// RunProgress runs the Bubbletea progress program and returns when complete.
func RunProgress(steps []InstallStep, presetName, stackName string) error {
	model := NewProgressModel(steps, presetName, stackName)
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}
