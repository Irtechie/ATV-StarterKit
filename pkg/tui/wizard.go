package tui

import (
	"fmt"
	"strings"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/detect"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/gstack"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
	"github.com/charmbracelet/huh"
)

// WizardResult holds the user's selections from the guided wizard.
type WizardResult struct {
	Stack               detect.Stack             // primary stack for root templates and progress display
	StackPacks          []installstate.StackPack // additive selected packs for stack-specific assets
	Components          []string                 // legacy ATV layer keys (for backward compat)
	ATVLayers           []string                 // parsed ATV layer keys
	GstackDirs          []string                 // selected gstack skill directories
	GstackRuntime       bool                     // whether to build the gstack TS binary
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
		Stack:      detected.Stack,
		StackPacks: installstate.AllStackPacks(),
	}

	prereqs := gstack.DetectPrerequisites()

	// ── Screen 1: Stack selection ──
	selectedStackPacks := []string{"general", "typescript", "python", "rails"}

	stackForm := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Choose stack packs for this repo").
				Description(stackSelectionDescription(detected)).
				Options(
					huh.NewOption("General — shared/base guidance", "general").Selected(true),
					huh.NewOption("TypeScript", "typescript").Selected(true),
					huh.NewOption("Python", "python").Selected(true),
					huh.NewOption("Rails", "rails").Selected(true),
				).
				Value(&selectedStackPacks),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := stackForm.Run(); err != nil {
		return nil, fmt.Errorf("wizard cancelled: %w", err)
	}

	parsedPacks, err := parseStackPackSelections(selectedStackPacks)
	if err != nil {
		return nil, err
	}
	result.StackPacks = parsedPacks
	result.Stack = detect.PrimaryStackForPacks(parsedPacks, detected.Stack)

	// ── Screen 2: Preset selection ──
	selectedPreset := "starter" // default to starter so all 3 are visible
	presetOptions := make([]huh.Option[string], 0, 3)
	for _, p := range AllPresets() {
		label := fmt.Sprintf("%s %s — %s", p.Emoji, p.Name, p.Description)
		opt := huh.NewOption(label, p.Key)
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

	result.ATVLayers = atvLayers
	result.GstackDirs = gstackDirs
	result.Components = atvLayers
	result.GstackRuntime = preset.ShouldEnableGstackRuntime(gstackDirs, prereqs)
	result.IncludeAgentBrowser = includeAgentBrowser
	result.PresetName = preset.Name

	return result, nil
}

func stackSelectionDescription(detected detect.Environment) string {
	if len(detected.DetectedPacks) == 1 && detected.DetectedPacks[0] == installstate.StackPackGeneral {
		return "All packs start selected. No strong stack signals were detected, so General is your safest base pack."
	}

	labels := make([]string, 0, len(detected.DetectedPacks))
	for _, pack := range detected.DetectedPacks {
		labels = append(labels, stackPackLabel(pack))
	}
	return fmt.Sprintf("All packs start selected. Likely matches in this repo: %s.", strings.Join(labels, ", "))
}

func parseStackPackSelections(selected []string) ([]installstate.StackPack, error) {
	packs := make([]installstate.StackPack, 0, len(selected))
	for _, value := range selected {
		packs = append(packs, installstate.StackPack(value))
	}
	normalized, err := installstate.NormalizeStackPacks(packs)
	if err != nil {
		return nil, fmt.Errorf("invalid stack pack selection: %w", err)
	}
	if err := installstate.ValidateStackPacks(normalized); err != nil {
		return nil, fmt.Errorf("choose at least one stack pack: %w", err)
	}
	return normalized, nil
}

func stackPackLabel(pack installstate.StackPack) string {
	switch pack {
	case installstate.StackPackRails:
		return "Rails"
	case installstate.StackPackPython:
		return "Python"
	case installstate.StackPackTypeScript:
		return "TypeScript"
	default:
		return "General"
	}
}
