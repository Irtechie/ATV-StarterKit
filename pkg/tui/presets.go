package tui

import "github.com/All-The-Vibes/ATV-StarterKit/pkg/gstack"

// Preset defines a curated skill selection level.
type Preset struct {
	Key                 string
	Name                string
	Emoji               string
	Description         string
	Detail              string
	ATVLayers           []string
	GstackDirs          []string
	EnableGstackRuntime bool
	IncludeAgentBrowser bool
	NeedsBun            bool
}

// StarterPreset includes core ATV skills only — no network calls, instant install.
var StarterPreset = Preset{
	Key:         "starter",
	Name:        "Starter",
	Emoji:       "⚡",
	Description: "Core workflow (13 skills, instant install)",
	Detail:      "Plan, build, review, compound. No browser tools or gstack.",
	ATVLayers: []string{
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
	},
	EnableGstackRuntime: false,
	GstackDirs:          nil,
	IncludeAgentBrowser: false,
}

// ProPreset adds gstack text-only sprint skills — no browser QA.
var ProPreset = Preset{
	Key:         "pro",
	Name:        "Pro",
	Emoji:       "🚀",
	Description: "Full sprint process (35+ skills)",
	Detail:      "+ gstack: review, ship, safety guardrails, security, debugging, retro",
	ATVLayers: []string{
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
	},
	GstackDirs: []string{
		"office-hours", "plan-ceo-review", "plan-eng-review", "plan-design-review",
		"design-consultation", "autoplan",
		"review", "design-review", "design-shotgun", "codex",
		"cso",
		"ship", "land-and-deploy", "canary", "document-release",
		"careful", "freeze", "guard", "unfreeze",
		"investigate",
		"retro",
		"learn",
	},
	EnableGstackRuntime: false,
	IncludeAgentBrowser: false,
}

// FullPreset installs everything — all gstack skills including browser QA.
var FullPreset = Preset{
	Key:         "full",
	Name:        "Full",
	Emoji:       "🔥",
	Description: "Complete engineering team (45+ skills)",
	Detail:      "+ browser QA, benchmarks, agent-browser + Chrome. Requires Bun, ~2min",
	ATVLayers: []string{
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
	},
	EnableGstackRuntime: true,
	GstackDirs:          allGstackDirs(),
	IncludeAgentBrowser: true,
	NeedsBun:            true,
}

// ShouldEnableGstackRuntime returns whether the guided flow should enable gstack runtime work.
func (p Preset) ShouldEnableGstackRuntime(selectedDirs []string, prereqs gstack.Prerequisites) bool {
	return p.EnableGstackRuntime && prereqs.HasBun && len(selectedDirs) > 0
}

// AllPresets returns the ordered list of presets.
func AllPresets() []Preset {
	return []Preset{StarterPreset, ProPreset, FullPreset}
}

func allGstackDirs() []string {
	var dirs []string
	for _, s := range gstack.AllSkills() {
		dirs = append(dirs, s.Dir)
	}
	return dirs
}
