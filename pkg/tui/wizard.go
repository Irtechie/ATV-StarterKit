package tui

import (
	"fmt"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/detect"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/gstack"
	"github.com/charmbracelet/huh"
)

// WizardResult holds the user's selections from the guided wizard.
type WizardResult struct {
	Stack               detect.Stack
	Components          []string // legacy ATV layer keys (for backward compat)
	ATVLayers           []string // parsed ATV layer keys
	GstackDirs          []string // selected gstack skill directories
	GstackRuntime       bool     // whether to build the gstack TS binary
	IncludeAgentBrowser bool
	PresetName          string // which preset was selected (for progress display)
}

// component layer keys
const (
	LayerCoreSkills          = "core-skills"
	LayerOrchestrators       = "orchestrators"
	LayerUniversalAgents     = "universal-agents"
	LayerStackAgents         = "stack-agents"
	LayerMCPServers          = "mcp-servers"
	LayerVSCodeExtensions    = "vscode-extensions"
	LayerCopilotInstructions = "copilot-instructions"
	LayerSetupSteps          = "setup-steps"
	LayerFileInstructions    = "file-instructions"
	LayerDocsStructure       = "docs-structure"
	LayerLocalConfig         = "local-config"
)

// AllLayers returns all available component layer keys.
func AllLayers() []string {
	return []string{
		LayerCoreSkills,
		LayerOrchestrators,
		LayerUniversalAgents,
		LayerStackAgents,
		LayerMCPServers,
		LayerVSCodeExtensions,
		LayerCopilotInstructions,
		LayerSetupSteps,
		LayerFileInstructions,
		LayerDocsStructure,
	}
}

// RunWizard runs the interactive guided mode TUI with preset selection.
func RunWizard(detected detect.Environment) (*WizardResult, error) {
	result := &WizardResult{
		Stack: detected.Stack,
	}

	prereqs := gstack.DetectPrerequisites()

	// ── Screen 1: Stack selection ──
	var selectedStack string
	switch detected.Stack {
	case detect.StackTypeScript:
		selectedStack = "typescript"
	case detect.StackPython:
		selectedStack = "python"
	case detect.StackRails:
		selectedStack = "rails"
	default:
		selectedStack = "general"
	}

	stackForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What's your primary stack?").
				Description(fmt.Sprintf("Auto-detected: %s (%s)", detected.Stack, detected.StackHint)).
				Options(
					huh.NewOption("TypeScript", "typescript"),
					huh.NewOption("Python", "python"),
					huh.NewOption("Rails", "rails"),
					huh.NewOption("General", "general"),
				).
				Value(&selectedStack),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := stackForm.Run(); err != nil {
		return nil, fmt.Errorf("wizard cancelled: %w", err)
	}

	// ── Screen 2: Preset selection ──
	var selectedPreset string
	presetOptions := make([]huh.Option[string], 0, 3)
	for _, p := range AllPresets() {
		label := fmt.Sprintf("%s %s — %s", p.Emoji, p.Name, p.Description)
		opt := huh.NewOption(label, p.Key)
		if p.Key == "pro" {
			opt = opt.Selected(true)
		}
		presetOptions = append(presetOptions, opt)
	}

	presetForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your setup level").
				Description(fmt.Sprintf("Prerequisites: %s", prereqs.Summary())).
				Options(presetOptions...).
				Value(&selectedPreset),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := presetForm.Run(); err != nil {
		return nil, fmt.Errorf("wizard cancelled: %w", err)
	}

	// Find the selected preset
	var preset Preset
	for _, p := range AllPresets() {
		if p.Key == selectedPreset {
			preset = p
			break
		}
	}

	// ── Screen 3: Customize? ──
	var wantCustomize bool

	customizeForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Want to customize individual skills?").
				Description(fmt.Sprintf("Selected: %s %s — %s", preset.Emoji, preset.Name, preset.Detail)).
				Affirmative("Yes, let me pick").
				Negative("No, install preset as-is").
				Value(&wantCustomize),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := customizeForm.Run(); err != nil {
		return nil, fmt.Errorf("wizard cancelled: %w", err)
	}

	// Build the final selections
	atvLayers := preset.ATVLayers
	gstackDirs := preset.GstackDirs
	includeAgentBrowser := preset.IncludeAgentBrowser

	// ── Screen 4 (optional): Customize ──
	if wantCustomize {
		groups := BuildCategoryGroups(prereqs)

		// Build pre-selected set from preset
		presetGstackSet := make(map[string]bool)
		for _, d := range preset.GstackDirs {
			presetGstackSet[d] = true
		}

		var skillOptions []huh.Option[string]
		for _, group := range groups {
			for _, skill := range group.Skills {
				label := skill.Label
				if skill.IsGstack {
					label = fmt.Sprintf("[gstack] %s", label)
				}
				if skill.RequiresBun && !prereqs.HasBun {
					label = fmt.Sprintf("%s ⚠️ (requires Bun)", label)
				}

				opt := huh.NewOption(label, skill.Key)
				// Pre-select based on preset
				if skill.IsGstack {
					dir := skill.Key
					if len(dir) > 7 && dir[:7] == "gstack:" {
						dir = dir[7:]
					}
					if presetGstackSet[dir] && (!skill.RequiresBun || prereqs.HasBun) {
						opt = opt.Selected(true)
					}
				} else {
					opt = opt.Selected(true) // ATV skills always pre-selected
				}
				skillOptions = append(skillOptions, opt)
			}
		}

		var infraOptions []huh.Option[string]
		for _, infra := range InfraLayers {
			infraOptions = append(infraOptions, huh.NewOption(infra.Label, infra.Key).Selected(true))
		}

		var selectedSkills []string
		var selectedInfra []string

		customForm := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Workflow skills").
					Description("Toggle individual skills on/off.").
					Options(skillOptions...).
					Value(&selectedSkills),
			),
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Infrastructure & configuration").
					Description("All selected by default. Deselect any you don't need.").
					Options(infraOptions...).
					Value(&selectedInfra),
			),
		).WithTheme(huh.ThemeCatppuccin())

		if err := customForm.Run(); err != nil {
			return nil, fmt.Errorf("wizard cancelled: %w", err)
		}

		parsedATV, parsedGstack := ParseSelections(selectedSkills)
		atvLayers = append(parsedATV, selectedInfra...)
		gstackDirs = parsedGstack

		// Check if agent-browser was selected
		includeAgentBrowser = false
		for _, key := range parsedATV {
			if key == "agent-browser" {
				includeAgentBrowser = true
			}
		}
	}

	// Map stack
	switch selectedStack {
	case "typescript":
		result.Stack = detect.StackTypeScript
	case "python":
		result.Stack = detect.StackPython
	case "rails":
		result.Stack = detect.StackRails
	default:
		result.Stack = detect.StackGeneral
	}

	result.ATVLayers = atvLayers
	result.GstackDirs = gstackDirs
	result.Components = atvLayers
	result.GstackRuntime = prereqs.HasBun && len(gstackDirs) > 0
	result.IncludeAgentBrowser = includeAgentBrowser
	result.PresetName = preset.Name

	return result, nil
}
