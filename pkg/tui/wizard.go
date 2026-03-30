package tui

import (
	"fmt"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/detect"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/gstack"
	"github.com/charmbracelet/huh"
)

// WizardResult holds the user's selections from the guided wizard.
type WizardResult struct {
	Stack         detect.Stack
	Components    []string // legacy ATV layer keys (for backward compat)
	ATVLayers     []string // parsed ATV layer keys
	GstackDirs    []string // selected gstack skill directories
	GstackRuntime bool     // whether to build the gstack TS binary
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

// RunWizard runs the interactive guided mode TUI.
func RunWizard(detected detect.Environment) (*WizardResult, error) {
	result := &WizardResult{
		Stack: detected.Stack,
	}

	// Detect prerequisites for gstack runtime status
	prereqs := gstack.DetectPrerequisites()

	// Step 1: Confirm or override detected stack
	stackOptions := []huh.Option[string]{
		huh.NewOption("TypeScript", "typescript"),
		huh.NewOption("Python", "python"),
		huh.NewOption("Rails", "rails"),
		huh.NewOption("General", "general"),
	}

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

	// Step 2: Build category-based skill options
	groups := BuildCategoryGroups(prereqs)

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
			// Default: select all ATV skills, select gstack skills that don't require unavailable runtime
			if !skill.RequiresBun || prereqs.HasBun {
				opt = opt.Selected(true)
			}
			skillOptions = append(skillOptions, opt)
		}
	}

	// Step 3: Infrastructure layers (unchanged from original)
	var infraOptions []huh.Option[string]
	for _, infra := range InfraLayers {
		infraOptions = append(infraOptions, huh.NewOption(infra.Label, infra.Key).Selected(true))
	}

	var selectedSkills []string
	var selectedInfra []string

	// Runtime status display
	runtimeDesc := fmt.Sprintf("Prerequisites: %s", prereqs.Summary())

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What's your primary stack?").
				Description(fmt.Sprintf("Auto-detected: %s (%s)", detected.Stack, detected.StackHint)).
				Options(stackOptions...).
				Value(&selectedStack),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Workflow skills (ATV + gstack)").
				Description(runtimeDesc).
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

	err := form.Run()
	if err != nil {
		return nil, fmt.Errorf("wizard cancelled: %w", err)
	}

	// Map selected stack string back to Stack type
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

	// Parse selections into ATV layers and gstack dirs
	atvLayers, gstackDirs := ParseSelections(selectedSkills)
	atvLayers = append(atvLayers, selectedInfra...)
	result.ATVLayers = atvLayers
	result.GstackDirs = gstackDirs
	result.Components = atvLayers // backward compat
	result.GstackRuntime = prereqs.HasBun && len(gstackDirs) > 0

	return result, nil
}
