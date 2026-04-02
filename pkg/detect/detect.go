package detect

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
)

// Stack represents the detected project stack.
type Stack string

const (
	StackRails      Stack = "rails"
	StackPython     Stack = "python"
	StackTypeScript Stack = "typescript"
	StackGeneral    Stack = "general"
)

// Environment holds detection results.
type Environment struct {
	Stack         Stack
	DetectedPacks []installstate.StackPack
	IsGitRepo     bool
	StackHint     string // human-readable detection reason
}

// DetectEnvironment scans the target directory for stack indicators.
func DetectEnvironment(dir string) Environment {
	env := Environment{Stack: StackGeneral}
	var detectedPacks []installstate.StackPack
	var hints []string

	// Check for git repo
	if exists(filepath.Join(dir, ".git")) {
		env.IsGitRepo = true
	}

	// Rails: Gemfile + config/routes.rb
	if exists(filepath.Join(dir, "Gemfile")) && exists(filepath.Join(dir, "config", "routes.rb")) {
		detectedPacks = append(detectedPacks, installstate.StackPackRails)
		hints = append(hints, "Gemfile + config/routes.rb")
	}

	// TypeScript: tsconfig.json
	if exists(filepath.Join(dir, "tsconfig.json")) {
		detectedPacks = append(detectedPacks, installstate.StackPackTypeScript)
		hints = append(hints, "tsconfig.json")
	}

	// Python: pyproject.toml or requirements.txt
	if exists(filepath.Join(dir, "pyproject.toml")) {
		detectedPacks = append(detectedPacks, installstate.StackPackPython)
		hints = append(hints, "pyproject.toml")
	}
	if exists(filepath.Join(dir, "requirements.txt")) {
		if !containsPack(detectedPacks, installstate.StackPackPython) {
			detectedPacks = append(detectedPacks, installstate.StackPackPython)
		}
		hints = append(hints, "requirements.txt")
	}

	normalized, err := installstate.NormalizeStackPacks(detectedPacks)
	if err != nil || len(normalized) == 0 {
		env.DetectedPacks = []installstate.StackPack{installstate.StackPackGeneral}
		env.Stack = StackGeneral
		env.StackHint = "no stack-specific files detected"
		return env
	}

	env.DetectedPacks = normalized
	env.Stack = PrimaryStackForPacks(normalized, StackGeneral)
	env.StackHint = fmtHints(hints)
	return env
}

// StackPackForStack converts a singular stack into the matching stack pack key.
func StackPackForStack(stack Stack) installstate.StackPack {
	switch stack {
	case StackRails:
		return installstate.StackPackRails
	case StackPython:
		return installstate.StackPackPython
	case StackTypeScript:
		return installstate.StackPackTypeScript
	default:
		return installstate.StackPackGeneral
	}
}

// StackForPack converts a stack pack into the matching singular stack identifier.
func StackForPack(pack installstate.StackPack) Stack {
	switch pack {
	case installstate.StackPackRails:
		return StackRails
	case installstate.StackPackPython:
		return StackPython
	case installstate.StackPackTypeScript:
		return StackTypeScript
	default:
		return StackGeneral
	}
}

// PrimaryStackForPacks chooses the root/template stack for a selected set of packs.
func PrimaryStackForPacks(packs []installstate.StackPack, preferred Stack) Stack {
	normalized, err := installstate.NormalizeStackPacks(packs)
	if err != nil || len(normalized) == 0 {
		return StackGeneral
	}

	preferredPack := StackPackForStack(preferred)
	if containsPack(normalized, preferredPack) {
		return preferred
	}

	priority := []installstate.StackPack{
		installstate.StackPackRails,
		installstate.StackPackTypeScript,
		installstate.StackPackPython,
		installstate.StackPackGeneral,
	}
	for _, pack := range priority {
		if containsPack(normalized, pack) {
			return StackForPack(pack)
		}
	}

	return StackGeneral
}

func containsPack(packs []installstate.StackPack, want installstate.StackPack) bool {
	for _, pack := range packs {
		if pack == want {
			return true
		}
	}
	return false
}

func fmtHints(hints []string) string {
	if len(hints) == 0 {
		return "no stack-specific files detected"
	}
	return strings.Join(hints, " + ") + " found"
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
