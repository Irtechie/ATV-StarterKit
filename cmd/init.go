package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/agentbrowser"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/detect"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/gstack"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/output"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/scaffold"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/tui"
	"github.com/spf13/cobra"
)

const installerModulePath = "module github.com/All-The-Vibes/ATV-StarterKit"

var guided bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ATV Starter Kit in current directory",
	Long: `Scaffold a complete agentic coding environment with all 6 Copilot lifecycle hooks.

Default: auto-detects your stack and installs everything (zero questions).
Use --guided for interactive mode with component selection.`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVar(&guided, "guided", false, "Interactive mode with component selection")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	printer := output.NewPrinter()
	targetDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Show banner
	printer.PrintBanner()

	// Phase 1: Detect environment
	env := detect.DetectEnvironment(targetDir)
	printer.PrintDetection(env)

	// Phase 2: Determine components
	var catalog []scaffold.Component
	var gstackDirs []string
	var gstackRuntime bool
	var installAgentBrowser bool
	var presetName string

	if guided {
		// Interactive TUI wizard
		result, err := tui.RunWizard(env)
		if err != nil {
			return err
		}
		catalog = scaffold.BuildFilteredCatalog(result.Stack, result.Components)
		gstackDirs = result.GstackDirs
		gstackRuntime = result.GstackRuntime
		installAgentBrowser = result.IncludeAgentBrowser
		presetName = result.PresetName

		// Build install steps for progress display
		steps := buildInstallSteps(targetDir, catalog, gstackDirs, gstackRuntime, installAgentBrowser)

		// Run with animated progress
		if err := tui.RunProgress(steps, presetName, string(result.Stack)); err != nil {
			return fmt.Errorf("install failed: %w", err)
		}

		// Print summary after progress completes
		printer.PrintNextSteps(len(gstackDirs) > 0, installAgentBrowser)
	} else {
		// One-click mode — install everything for detected stack (ATV only, no gstack)
		catalog = scaffold.BuildCatalog(env.Stack)

		// Phase 3: Write ATV files
		results := scaffold.WriteAll(targetDir, catalog)

		// Phase 4: Print summary
		printer.PrintResults(results)
		printer.PrintNextSteps(false, false)
	}

	// Update plan checkboxes only when running inside the installer repository.
	if isInstallerRepo(targetDir) {
		printer.Info("Plan directory detected. Update plan checkboxes manually.")
	}

	return nil
}

// buildInstallSteps creates the ordered list of install steps for the progress display.
func buildInstallSteps(targetDir string, catalog []scaffold.Component, gstackDirs []string, gstackRuntime bool, installAgentBrowser bool) []tui.InstallStep {
	var steps []tui.InstallStep

	// Step 1: ATV scaffold (always)
	steps = append(steps, tui.InstallStep{
		Name: "Scaffolding ATV files",
		Action: func() error {
			scaffold.WriteAll(targetDir, catalog)
			return nil
		},
	})

	// Step 2: gstack clone (if selected)
	if len(gstackDirs) > 0 {
		mode := gstack.ModeMarkdownOnly
		if gstackRuntime {
			mode = gstack.ModeFullRuntime
		}
		steps = append(steps, tui.InstallStep{
			Name: "Cloning gstack",
			Action: func() error {
				result := gstack.Install(targetDir, mode)
				if result.Error != nil {
					return result.Error
				}
				return nil
			},
		})
	}

	// Step 3: agent-browser (if selected)
	if installAgentBrowser {
		steps = append(steps, tui.InstallStep{
			Name: "Installing agent-browser + Chrome",
			Action: func() error {
				result := agentbrowser.Install(targetDir)
				if result.Error != nil {
					return result.Error
				}
				return nil
			},
		})
	}

	return steps
}

func isInstallerRepo(dir string) bool {
	goModPath := filepath.Join(dir, "go.mod")
	goMod, err := os.ReadFile(goModPath)
	if err != nil {
		return false
	}

	templatesPath := filepath.Join(dir, "pkg", "scaffold", "templates")
	if info, err := os.Stat(templatesPath); err != nil || !info.IsDir() {
		return false
	}

	return bytes.Contains(goMod, []byte(installerModulePath))
}
