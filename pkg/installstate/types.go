package installstate

import (
	"errors"
	"slices"
	"time"
)

// StackPack represents an additive guided-install stack selection.
type StackPack string

const (
	StackPackGeneral    StackPack = "general"
	StackPackTypeScript StackPack = "typescript"
	StackPackPython     StackPack = "python"
	StackPackRails      StackPack = "rails"
)

var orderedStackPacks = []StackPack{
	StackPackGeneral,
	StackPackTypeScript,
	StackPackPython,
	StackPackRails,
}

// AllStackPacks returns the supported stack packs in deterministic display order.
func AllStackPacks() []StackPack {
	return slices.Clone(orderedStackPacks)
}

// NormalizeStackPacks deduplicates and orders stack packs using the canonical pack order.
func NormalizeStackPacks(packs []StackPack) ([]StackPack, error) {
	seen := make(map[StackPack]bool, len(packs))
	for _, pack := range packs {
		if !IsValidStackPack(pack) {
			return nil, ErrInvalidStackPack
		}
		seen[pack] = true
	}

	normalized := make([]StackPack, 0, len(seen))
	for _, pack := range orderedStackPacks {
		if seen[pack] {
			normalized = append(normalized, pack)
		}
	}

	return normalized, nil
}

// ValidateStackPacks enforces the current Phase 0 contract: at least one stack pack must be selected.
func ValidateStackPacks(packs []StackPack) error {
	normalized, err := NormalizeStackPacks(packs)
	if err != nil {
		return err
	}
	if len(normalized) == 0 {
		return ErrNoStackPacks
	}
	return nil
}

// IsValidStackPack reports whether the value is a recognized guided-install stack pack.
func IsValidStackPack(pack StackPack) bool {
	for _, candidate := range orderedStackPacks {
		if pack == candidate {
			return true
		}
	}
	return false
}

// RerunPolicy describes how guided installs behave when run repeatedly.
type RerunPolicy string

const (
	RerunPolicyAdditiveOnly RerunPolicy = "additive-only"
)

// InstallStepStatus captures structured installer outcomes for manifesting and telemetry.
type InstallStepStatus string

const (
	InstallStepPending InstallStepStatus = "pending"
	InstallStepRunning InstallStepStatus = "running"
	InstallStepDone    InstallStepStatus = "done"
	InstallStepWarning InstallStepStatus = "warning"
	InstallStepFailed  InstallStepStatus = "failed"
	InstallStepSkipped InstallStepStatus = "skipped"
)

// InstallOutcome captures the machine-readable result of one installer action.
type InstallOutcome struct {
	Step     string            `json:"step"`
	Status   InstallStepStatus `json:"status"`
	Reason   string            `json:"reason,omitempty"`
	Duration string            `json:"duration,omitempty"`
}

// RequestedState records what the guided installer attempted to install.
type RequestedState struct {
	StackPacks          []StackPack `json:"stackPacks"`
	ATVLayers           []string    `json:"atvLayers,omitempty"`
	GstackDirs          []string    `json:"gstackDirs,omitempty"`
	GstackRuntime       bool        `json:"gstackRuntime"`
	IncludeAgentBrowser bool        `json:"includeAgentBrowser"`
	PresetName          string      `json:"presetName,omitempty"`
}

// Recommendation is the deterministic, local suggestion shape that later launchpad work will consume.
type Recommendation struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Reason   string `json:"reason"`
	Priority int    `json:"priority"`
}

// InstallManifest is the canonical installer-state snapshot written after guided installs.
type InstallManifest struct {
	Version         int              `json:"version"`
	GeneratedAt     time.Time        `json:"generatedAt"`
	RerunPolicy     RerunPolicy      `json:"rerunPolicy"`
	Requested       RequestedState   `json:"requested"`
	Outcomes        []InstallOutcome `json:"outcomes,omitempty"`
	Recommendations []Recommendation `json:"recommendations,omitempty"`
}

var (
	ErrInvalidStackPack = errors.New("invalid stack pack")
	ErrNoStackPacks     = errors.New("at least one stack pack must be selected")
)
