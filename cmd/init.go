package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/agentbrowser"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/detect"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/gstack"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/output"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/scaffold"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/tui"
	"github.com/spf13/cobra"
)

const installerModulePath = "module github.com/All-The-Vibes/ATV-StarterKit"

const (
	scaffoldStepName     = "Scaffolding ATV files"
	gstackStepName       = "Syncing gstack skills"
	agentBrowserStepName = "Installing agent-browser"
)

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
		catalog = scaffold.BuildFilteredCatalogForPacks(result.StackPacks, result.Stack, result.Components)
		gstackDirs = result.GstackDirs
		gstackRuntime = result.GstackRuntime
		installAgentBrowser = result.IncludeAgentBrowser
		presetName = result.PresetName

		// Build install steps for progress display
		steps := buildInstallSteps(targetDir, catalog, gstackDirs, gstackRuntime, installAgentBrowser)

		// Run with animated progress
		outcomes, err := tui.RunProgress(steps, presetName, string(result.Stack))
		if err != nil {
			return fmt.Errorf("install failed: %w", err)
		}

		manifestPath := ""

		// Collect file paths from catalog for checksum computation
		var catalogPaths []string
		for _, comp := range catalog {
			if !comp.IsDir && len(comp.Content) > 0 {
				catalogPaths = append(catalogPaths, comp.Path)
			}
		}

		manifest := installstate.InstallManifest{
			Requested: installstate.RequestedState{
				StackPacks:          result.StackPacks,
				ATVLayers:           result.ATVLayers,
				GstackDirs:          result.GstackDirs,
				GstackRuntime:       result.GstackRuntime,
				IncludeAgentBrowser: result.IncludeAgentBrowser,
				PresetName:          result.PresetName,
			},
			Outcomes:       outcomes,
			FileChecksums:  installstate.ComputeFileChecksums(targetDir, catalogPaths),
		}
		manifest.Recommendations = installstate.BuildRecommendations(targetDir, manifest)
		if err := installstate.WriteManifest(targetDir, manifest); err != nil {
			printer.Info(fmt.Sprintf("⚠️  Failed to write guided install manifest: %v", err))
		} else {
			manifestPath = filepath.ToSlash(filepath.Join(".atv", "install-manifest.json"))
		}

		// Print summary after progress completes
		printer.PrintGuidedSummary(outcomes, manifestPath)
		printer.PrintRecommendations(manifest.Recommendations)
		printer.PrintNextSteps(hasUsableOutcome(outcomes, gstackStepName), hasUsableOutcome(outcomes, agentBrowserStepName), manifestPath)
	} else {
		// One-click mode — install everything for detected stack (ATV only, no gstack)
		catalog = scaffold.BuildCatalog(env.Stack)

		// Phase 3: Write ATV files
		results := scaffold.WriteAll(targetDir, catalog)

		// Phase 4: Print summary
		printer.PrintResults(results)
		printer.PrintNextSteps(false, false, "")
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
		Name: scaffoldStepName,
		Action: func() tui.InstallStepResult {
			results := scaffold.WriteAll(targetDir, catalog)
			summary := scaffold.SummarizeResults(results)
			substeps := scaffold.ResultsToSubsteps(results)
			status := installstate.InstallStepDone
			if !summary.Successful() {
				status = installstate.InstallStepWarning
				if summary.Created == 0 && summary.Merged == 0 && summary.Directories == 0 {
					status = installstate.InstallStepFailed
				}
			}
			return tui.InstallStepResult{
				Status:   status,
				Detail:   summary.Detail(),
				Reason:   summary.FailureReason(),
				Substeps: substeps,
			}
		},
	})

	// Step 2: gstack clone (if selected)
	if len(gstackDirs) > 0 {
		mode := gstack.ModeMarkdownOnly
		if gstackRuntime {
			mode = gstack.ModeFullRuntime
		}
		steps = append(steps, tui.InstallStep{
			Name: gstackStepName,
			Action: func() tui.InstallStepResult {
				result := gstack.Install(targetDir, mode)
				if result.Error != nil {
					return tui.InstallStepResult{
						Status: installstate.InstallStepFailed,
						Reason: result.Error.Error(),
						Error:  result.Error,
					}
				}
				detail := fmt.Sprintf("%d skills synced", len(result.SkillDirs))
				if result.Built {
					detail = detail + ", runtime ready"
				} else if mode == gstack.ModeMarkdownOnly {
					detail = detail + ", markdown-only"
				}
				status := installstate.InstallStepDone
				reason := ""
				var skipReason installstate.SkipReason
				if result.Warning != "" {
					status = installstate.InstallStepWarning
					reason = result.Warning
					if result.Copied && strings.Contains(result.Warning, "already cloned") {
						status = installstate.InstallStepSkipped
						skipReason = installstate.SkipReasonAlreadyInstalled
					}
				}

				// Build substeps for clone, build/doc-gen, and skill copy
				var substeps []installstate.InstallOutcome
				if result.Cloned {
					substeps = append(substeps, installstate.InstallOutcome{
						Step: "git clone", Status: installstate.InstallStepDone, Detail: "shallow clone completed",
					})
				} else if result.Copied {
					substeps = append(substeps, installstate.InstallOutcome{
						Step: "git clone", Status: installstate.InstallStepSkipped,
						Detail: "already cloned", SkipReason: installstate.SkipReasonAlreadyInstalled,
					})
				}
				if result.Built {
					substeps = append(substeps, installstate.InstallOutcome{
						Step: "runtime build", Status: installstate.InstallStepDone, Detail: "setup completed",
					})
				} else if mode == gstack.ModeFullRuntime {
					substeps = append(substeps, installstate.InstallOutcome{
						Step: "runtime build", Status: installstate.InstallStepWarning, Detail: "build failed, fell back to docs",
						Reason: result.Warning,
					})
				} else {
					substeps = append(substeps, installstate.InstallOutcome{
						Step: "runtime build", Status: installstate.InstallStepSkipped,
						Detail: "markdown-only mode", SkipReason: installstate.SkipReasonUserSkip,
					})
				}
				if result.Copied {
					substeps = append(substeps, installstate.InstallOutcome{
						Step:   "copy skills",
						Status: installstate.InstallStepDone,
						Detail: fmt.Sprintf("%d skill dirs", len(result.SkillDirs)),
					})
				}

				return tui.InstallStepResult{
					Status:     status,
					Detail:     detail,
					Reason:     reason,
					SkipReason: skipReason,
					Substeps:   substeps,
				}
			},
		})
	}

	// Step 3: agent-browser (if selected)
	if installAgentBrowser {
		steps = append(steps, tui.InstallStep{
			Name: agentBrowserStepName,
			Action: func() tui.InstallStepResult {
				result := agentbrowser.Install(targetDir)
				if result.Error != nil {
					return tui.InstallStepResult{
						Status: installstate.InstallStepFailed,
						Reason: result.Error.Error(),
						Error:  result.Error,
					}
				}

				detailParts := make([]string, 0, 2)
				if result.Installed {
					detailParts = append(detailParts, "CLI ready")
				}
				if result.SkillCopied {
					detailParts = append(detailParts, "skill copied")
				}
				detail := strings.Join(detailParts, ", ")
				if detail == "" {
					detail = "no local changes"
				}

				status := installstate.InstallStepDone
				if result.Warning != "" {
					status = installstate.InstallStepWarning
				}

				// Build substeps for npm install, chrome download, and skill copy
				var substeps []installstate.InstallOutcome
				if result.Installed {
					substeps = append(substeps, installstate.InstallOutcome{
						Step: "npm install", Status: installstate.InstallStepDone, Detail: "agent-browser CLI available",
					})
				} else {
					substeps = append(substeps, installstate.InstallOutcome{
						Step: "npm install", Status: installstate.InstallStepWarning,
						Detail: "install skipped or failed", Reason: result.Warning,
						SkipReason: installstate.SkipReasonPrereqMissing,
					})
				}
				if result.SkillCopied {
					substeps = append(substeps, installstate.InstallOutcome{
						Step: "copy SKILL.md", Status: installstate.InstallStepDone, Detail: "skill registered",
					})
				}

				return tui.InstallStepResult{
					Status:   status,
					Detail:   detail,
					Reason:   result.Warning,
					Substeps: substeps,
				}
			},
		})
	}

	return steps
}

func hasUsableOutcome(outcomes []installstate.InstallOutcome, stepName string) bool {
	for _, outcome := range outcomes {
		if outcome.Step != stepName {
			continue
		}
		return outcome.Status != installstate.InstallStepFailed
	}
	return false
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
