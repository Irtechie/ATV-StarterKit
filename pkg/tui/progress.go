package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
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
	StepWarning
	StepFailed
	StepSkipped
)

// InstallStepResult is the structured result produced by one install action.
type InstallStepResult struct {
	Status     installstate.InstallStepStatus
	Detail     string
	Reason     string
	SkipReason installstate.SkipReason
	Substeps   []installstate.InstallOutcome
	Error      error
}

// InstallStep defines a single step in the install process.
type InstallStep struct {
	Name     string
	Detail   string
	Reason   string
	Duration string
	Status   StepStatus
	Error    error
	Action   func() InstallStepResult // the actual work to perform
}

// stepCompleteMsg is sent when a step finishes.
type stepCompleteMsg struct {
	index    int
	result   InstallStepResult
	duration time.Duration
}

// ProgressModel is a Bubbletea model for showing install progress.
type ProgressModel struct {
	steps      []InstallStep
	current    int
	spinner    spinner.Model
	done       bool
	presetName string
	stackName  string
	outcomes   []installstate.InstallOutcome
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
		result := msg.result
		if result.Status == "" {
			if result.Error != nil {
				result.Status = installstate.InstallStepFailed
			} else {
				result.Status = installstate.InstallStepDone
			}
		}

		m.steps[msg.index].Detail = result.Detail
		m.steps[msg.index].Reason = result.Reason
		m.steps[msg.index].Duration = formatStepDuration(msg.duration)
		if result.Error != nil {
			m.steps[msg.index].Error = result.Error
			if m.steps[msg.index].Reason == "" {
				m.steps[msg.index].Reason = result.Error.Error()
			}
		}

		switch result.Status {
		case installstate.InstallStepWarning:
			m.steps[msg.index].Status = StepWarning
		case installstate.InstallStepFailed:
			m.steps[msg.index].Status = StepFailed
		case installstate.InstallStepSkipped:
			m.steps[msg.index].Status = StepSkipped
		default:
			m.steps[msg.index].Status = StepDone
		}

		m.outcomes = append(m.outcomes, installstate.InstallOutcome{
			Step:       m.steps[msg.index].Name,
			Status:     result.Status,
			Detail:     m.steps[msg.index].Detail,
			Reason:     m.steps[msg.index].Reason,
			Duration:   m.steps[msg.index].Duration,
			SkipReason: result.SkipReason,
			Substeps:   result.Substeps,
		})

		// Advance to next step
		m.current = msg.index + 1
		if m.current >= len(m.steps) {
			m.done = true
			return m, tea.Quit
		}

		// Skip steps with nil actions
		for m.current < len(m.steps) && m.steps[m.current].Action == nil {
			m.steps[m.current].Status = StepSkipped
			m.steps[m.current].Reason = "no action required"
			m.outcomes = append(m.outcomes, installstate.InstallOutcome{
				Step:   m.steps[m.current].Name,
				Status: installstate.InstallStepSkipped,
				Reason: m.steps[m.current].Reason,
			})
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
	fmt.Fprintf(&b, "\n  Installing %s preset for %s...\n\n", m.presetName, m.stackName)

	// Steps
	for _, step := range m.steps {
		icon := m.stepIcon(step)
		label := step.Name
		if step.Detail != "" && (step.Status == StepDone || step.Status == StepWarning || step.Status == StepSkipped) {
			label = fmt.Sprintf("%s (%s)", step.Name, step.Detail)
		}
		if step.Reason != "" && (step.Status == StepWarning || step.Status == StepSkipped) {
			label = fmt.Sprintf("%s — %s", step.Name, step.Reason)
		}
		if step.Status == StepFailed {
			reason := step.Reason
			if reason == "" && step.Error != nil {
				reason = step.Error.Error()
			}
			if reason != "" {
				label = fmt.Sprintf("%s — %s", step.Name, reason)
			}
		}
		if step.Duration != "" && (step.Status == StepDone || step.Status == StepWarning || step.Status == StepSkipped) {
			label = fmt.Sprintf("%s · %s", label, step.Duration)
		}
		fmt.Fprintf(&b, "  %s %s\n", icon, label)
	}

	if m.done {
		b.WriteString("\n")
	}

	return b.String()
}

func (m ProgressModel) stepIcon(step InstallStep) string {
	doneStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	failStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	skipStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("243"))

	switch step.Status {
	case StepDone:
		return doneStyle.Render("✅")
	case StepWarning:
		return warnStyle.Render("⚠️ ")
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
		started := time.Now()
		var result InstallStepResult
		if step.Action != nil {
			result = step.Action()
		} else {
			result = InstallStepResult{Status: installstate.InstallStepSkipped, Reason: "no action required"}
		}
		return stepCompleteMsg{index: idx, result: result, duration: time.Since(started)}
	}
}

func formatStepDuration(duration time.Duration) string {
	if duration <= 0 {
		return ""
	}
	if duration < time.Second {
		return duration.Round(10 * time.Millisecond).String()
	}
	return duration.Round(100 * time.Millisecond).String()
}

// RunProgress runs the Bubbletea progress program and returns structured outcomes when complete.
func RunProgress(steps []InstallStep, presetName, stackName string) ([]installstate.InstallOutcome, error) {
	model := NewProgressModel(steps, presetName, stackName)
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}
	progressModel, ok := finalModel.(ProgressModel)
	if !ok {
		return nil, fmt.Errorf("unexpected progress model type %T", finalModel)
	}
	return progressModel.outcomes, nil
}
