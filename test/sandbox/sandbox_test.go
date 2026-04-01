package sandbox

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/concierge"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/detect"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/gstack"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/scaffold"
)

func TestAutoModeNoGstack(t *testing.T) {
	// Auto mode should install ATV files only, no gstack
	sandboxDir := t.TempDir()

	// Create a fake TypeScript project
	if err := os.WriteFile(filepath.Join(sandboxDir, "tsconfig.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	env := detect.DetectEnvironment(sandboxDir)
	if env.Stack != detect.StackTypeScript {
		t.Fatalf("expected TypeScript stack, got %s", env.Stack)
	}

	catalog := scaffold.BuildCatalog(env.Stack)
	results := scaffold.WriteAll(sandboxDir, catalog)

	if len(results) == 0 {
		t.Fatal("auto mode should produce results")
	}

	// ATV files should exist
	assertFileExists(t, filepath.Join(sandboxDir, ".github", "copilot-instructions.md"))
	assertFileExists(t, filepath.Join(sandboxDir, ".github", "copilot-setup-steps.yml"))

	// gstack should NOT exist
	gstackDir := filepath.Join(sandboxDir, ".github", "skills", "gstack")
	if _, err := os.Stat(gstackDir); !os.IsNotExist(err) {
		t.Error("auto mode should NOT install gstack")
	}

	manifestPath := filepath.Join(sandboxDir, ".atv", "install-manifest.json")
	if _, err := os.Stat(manifestPath); !os.IsNotExist(err) {
		t.Error("auto mode should NOT write a guided install manifest")
	}
}

func TestATVFilesCreatedCorrectly(t *testing.T) {
	sandboxDir := t.TempDir()

	// Create a fake Python project
	if err := os.WriteFile(filepath.Join(sandboxDir, "pyproject.toml"), []byte("[project]"), 0644); err != nil {
		t.Fatal(err)
	}

	env := detect.DetectEnvironment(sandboxDir)
	catalog := scaffold.BuildCatalog(env.Stack)
	results := scaffold.WriteAll(sandboxDir, catalog)

	// Count created files
	created := 0
	for _, r := range results {
		if r.Status == scaffold.StatusCreated || r.Status == scaffold.StatusDirCreated {
			created++
		}
	}

	if created == 0 {
		t.Error("should have created files")
	}

	// Check key ATV files
	assertFileExists(t, filepath.Join(sandboxDir, ".github", "copilot-instructions.md"))
	assertFileExists(t, filepath.Join(sandboxDir, ".github", "copilot-setup-steps.yml"))
	assertFileExists(t, filepath.Join(sandboxDir, ".github", "copilot-mcp-config.json"))
	assertDirExists(t, filepath.Join(sandboxDir, ".github", "skills"))
	assertDirExists(t, filepath.Join(sandboxDir, ".github", "agents"))
	assertDirExists(t, filepath.Join(sandboxDir, "docs", "plans"))
	assertDirExists(t, filepath.Join(sandboxDir, "docs", "brainstorms"))
}

func TestGstackInstallMarkdownOnly(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network-dependent test in short mode")
	}

	sandboxDir := t.TempDir()

	// Create .github/skills/ so copy has a target
	if err := os.MkdirAll(filepath.Join(sandboxDir, ".github", "skills"), 0755); err != nil {
		t.Fatal(err)
	}

	result := gstack.Install(sandboxDir, gstack.ModeMarkdownOnly)

	if result.Error != nil {
		t.Fatalf("gstack install failed: %v", result.Error)
	}

	if !result.Cloned {
		t.Error("expected gstack to be cloned")
	}

	// .gstack/ staging dir should exist
	gstackDir := filepath.Join(sandboxDir, ".gstack")
	if _, err := os.Stat(gstackDir); os.IsNotExist(err) {
		t.Error(".gstack/ staging directory should exist")
	}

	// Should not build in markdown-only mode
	if result.Built {
		t.Error("markdown-only mode should not build")
	}

	// Should have copied skill dirs to .github/skills/gstack-*
	if len(result.SkillDirs) == 0 {
		t.Error("expected at least some skill directories")
	}

	// Verify skills are at the right level for Copilot discovery
	for _, dir := range result.SkillDirs {
		skillPath := filepath.Join(sandboxDir, ".github", "skills", dir, "SKILL.md")
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			t.Errorf("expected skill file at %s", skillPath)
		}
	}
}

func TestGstackIdempotent(t *testing.T) {
	sandboxDir := t.TempDir()
	gstackDir := filepath.Join(sandboxDir, ".gstack")

	// Create a fake existing install
	if err := os.MkdirAll(gstackDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gstackDir, "SKILL.md"), []byte("# gstack"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(sandboxDir, ".github", "skills"), 0755); err != nil {
		t.Fatal(err)
	}

	result := gstack.Install(sandboxDir, gstack.ModeMarkdownOnly)
	if result.Error != nil {
		t.Fatalf("idempotent install should not error: %v", result.Error)
	}
	if result.Cloned {
		t.Error("should not clone when already installed")
	}
}

func TestInstructionsContainGstackSection(t *testing.T) {
	sandboxDir := t.TempDir()

	env := detect.DetectEnvironment(sandboxDir) // general stack
	catalog := scaffold.BuildCatalog(env.Stack)
	scaffold.WriteAll(sandboxDir, catalog)

	// Read copilot-instructions.md
	content, err := os.ReadFile(filepath.Join(sandboxDir, ".github", "copilot-instructions.md"))
	if err != nil {
		t.Fatalf("failed to read copilot-instructions.md: %v", err)
	}

	contentStr := string(content)

	// Should contain gstack section
	if !containsStr(contentStr, "gstack Skills") {
		t.Error("copilot-instructions.md should contain gstack Skills section")
	}

	// Should contain ATV override rules
	if !containsStr(contentStr, "ATV Override Rules") {
		t.Error("copilot-instructions.md should contain ATV Override Rules")
	}

	// Should contain protected artifacts
	if !containsStr(contentStr, "Protected artifacts") {
		t.Error("copilot-instructions.md should contain protected artifacts rule")
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", path)
	}
}

func assertDirExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Errorf("expected directory to exist: %s", path)
		return
	}
	if !info.IsDir() {
		t.Errorf("expected %s to be a directory", path)
	}
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// --- Integration test scenarios from the plan ---

func TestPrereqSkipProducesCorrectOutcomeAndManifest(t *testing.T) {
	// "Missing prerequisites downgrade or skip runtime-dependent capabilities
	//  with clear telemetry and manifest output"
	root := t.TempDir()

	// Simulate guided install with gstack runtime requested but build failed (prereq missing)
	manifest := installstate.InstallManifest{
		Requested: installstate.RequestedState{
			StackPacks:          []installstate.StackPack{installstate.StackPackGeneral},
			GstackDirs:          []string{"review"},
			GstackRuntime:       true,
			IncludeAgentBrowser: true,
			PresetName:          "Full",
		},
		Outcomes: []installstate.InstallOutcome{
			{Step: "Scaffolding ATV files", Status: installstate.InstallStepDone, Detail: "5 files created"},
			{
				Step: "Syncing gstack skills", Status: installstate.InstallStepWarning,
				Detail: "3 skills synced, markdown-only", Reason: "./setup failed, falling back to docs",
				Substeps: []installstate.InstallOutcome{
					{Step: "git clone", Status: installstate.InstallStepDone},
					{Step: "runtime build", Status: installstate.InstallStepWarning, Reason: "bun not found",
						SkipReason: installstate.SkipReasonPrereqMissing},
					{Step: "copy skills", Status: installstate.InstallStepDone, Detail: "3 skill dirs"},
				},
			},
			{
				Step: "Installing agent-browser", Status: installstate.InstallStepWarning,
				Reason: "npm not found",
				Substeps: []installstate.InstallOutcome{
					{Step: "npm install", Status: installstate.InstallStepWarning, Reason: "npm not found",
						SkipReason: installstate.SkipReasonPrereqMissing},
				},
			},
		},
	}

	if err := installstate.WriteManifest(root, manifest); err != nil {
		t.Fatalf("WriteManifest() error = %v", err)
	}

	read, err := installstate.ReadManifest(root)
	if err != nil {
		t.Fatalf("ReadManifest() error = %v", err)
	}

	// Verify substeps survive round-trip
	if len(read.Outcomes) != 3 {
		t.Fatalf("expected 3 outcomes, got %d", len(read.Outcomes))
	}
	gstackOutcome := read.Outcomes[1]
	if len(gstackOutcome.Substeps) != 3 {
		t.Fatalf("expected 3 gstack substeps, got %d", len(gstackOutcome.Substeps))
	}
	if gstackOutcome.Substeps[1].SkipReason != installstate.SkipReasonPrereqMissing {
		t.Fatalf("expected prereq-missing skip reason, got %q", gstackOutcome.Substeps[1].SkipReason)
	}

	// Verify recommendations prioritize the install issue
	recs := installstate.BuildRecommendations(root, read)
	if len(recs) == 0 || recs[0].ID != "fix-install-issues" {
		t.Fatalf("expected fix-install-issues as first recommendation, got %+v", recs)
	}
}

func TestPartialInstallWritesValidManifest(t *testing.T) {
	// "Partial install success still writes a valid manifest and useful completion summary"
	root := t.TempDir()

	manifest := installstate.InstallManifest{
		Requested: installstate.RequestedState{
			StackPacks:          []installstate.StackPack{installstate.StackPackGeneral, installstate.StackPackTypeScript},
			GstackDirs:          []string{"review"},
			IncludeAgentBrowser: true,
			PresetName:          "Full",
		},
		Outcomes: []installstate.InstallOutcome{
			{Step: "Scaffolding ATV files", Status: installstate.InstallStepDone, Detail: "created"},
			{Step: "Syncing gstack skills", Status: installstate.InstallStepFailed, Reason: "git clone failed"},
			{Step: "Installing agent-browser", Status: installstate.InstallStepFailed, Reason: "npm not found"},
		},
	}

	if err := installstate.WriteManifest(root, manifest); err != nil {
		t.Fatalf("WriteManifest() error = %v", err)
	}

	read, err := installstate.ReadManifest(root)
	if err != nil {
		t.Fatalf("ReadManifest() error = %v", err)
	}

	if read.Version != installstate.ManifestVersion {
		t.Fatalf("expected manifest version %d, got %d", installstate.ManifestVersion, read.Version)
	}
	if read.RerunPolicy != installstate.RerunPolicyAdditiveOnly {
		t.Fatalf("expected additive-only rerun policy, got %q", read.RerunPolicy)
	}

	summary := installstate.SummarizeOutcomes(read.Outcomes)
	if summary.Done != 1 || summary.Failed != 2 {
		t.Fatalf("expected 1 done + 2 failed, got %+v", summary)
	}

	// Even with failures, recommendations should be generated
	recs := installstate.BuildRecommendations(root, read)
	if len(recs) == 0 {
		t.Fatal("expected recommendations even with partial install")
	}
	if recs[0].ID != "fix-install-issues" {
		t.Fatalf("expected fix-install-issues first, got %s", recs[0].ID)
	}
}

func TestRerunPreservesManifestAdditively(t *testing.T) {
	// "Re-running guided mode after a prior install preserves or updates state"
	root := t.TempDir()

	// First run
	first := installstate.InstallManifest{
		Requested: installstate.RequestedState{
			StackPacks: []installstate.StackPack{installstate.StackPackGeneral},
			PresetName: "Starter",
		},
		Outcomes: []installstate.InstallOutcome{
			{Step: "Scaffolding ATV files", Status: installstate.InstallStepDone, Detail: "10 files created"},
		},
	}
	if err := installstate.WriteManifest(root, first); err != nil {
		t.Fatalf("first WriteManifest() error = %v", err)
	}

	read1, err := installstate.ReadManifest(root)
	if err != nil {
		t.Fatalf("first ReadManifest() error = %v", err)
	}
	gen1 := read1.GeneratedAt

	// Second run (simulates rerun with more packs)
	second := installstate.InstallManifest{
		Requested: installstate.RequestedState{
			StackPacks:    []installstate.StackPack{installstate.StackPackGeneral, installstate.StackPackTypeScript},
			GstackDirs:    []string{"review"},
			GstackRuntime: false,
			PresetName:    "Pro",
		},
		Outcomes: []installstate.InstallOutcome{
			{Step: "Scaffolding ATV files", Status: installstate.InstallStepDone,
				Detail: "2 files created",
				Substeps: []installstate.InstallOutcome{
					{Step: "copilot-instructions.md", Status: installstate.InstallStepSkipped,
						SkipReason: installstate.SkipReasonAlreadyInstalled},
					{Step: "tsconfig-agent.md", Status: installstate.InstallStepDone, Detail: "created"},
				}},
			{Step: "Syncing gstack skills", Status: installstate.InstallStepDone, Detail: "3 skills synced"},
		},
	}
	if err := installstate.WriteManifest(root, second); err != nil {
		t.Fatalf("second WriteManifest() error = %v", err)
	}

	read2, err := installstate.ReadManifest(root)
	if err != nil {
		t.Fatalf("second ReadManifest() error = %v", err)
	}

	// Manifest should reflect the latest run
	if len(read2.Requested.StackPacks) != 2 {
		t.Fatalf("expected 2 stack packs after rerun, got %d", len(read2.Requested.StackPacks))
	}
	if read2.Requested.PresetName != "Pro" {
		t.Fatalf("expected Pro preset after rerun, got %q", read2.Requested.PresetName)
	}
	if !read2.GeneratedAt.After(gen1) {
		t.Fatal("second manifest should have a later timestamp")
	}
	if read2.RerunPolicy != installstate.RerunPolicyAdditiveOnly {
		t.Fatalf("rerun policy should remain additive-only, got %q", read2.RerunPolicy)
	}
}

func TestLaunchpadMatchesManifestAndRepoMemory(t *testing.T) {
	// "Launchpad recommendation output matches manifest + repo memory for empty, partial, and mature repos"
	root := t.TempDir()

	// --- Empty repo (no manifest, no docs) ---
	snapshot, err := installstate.BuildLaunchpadSnapshot(root)
	if err != nil {
		t.Fatalf("empty snapshot error = %v", err)
	}
	if snapshot.HasManifest {
		t.Fatal("empty repo should not have manifest")
	}
	if len(snapshot.Recommendations) == 0 || snapshot.Recommendations[0].ID != "start-brainstorm" {
		t.Fatalf("empty repo should recommend start-brainstorm, got %+v", snapshot.Recommendations)
	}

	// --- Partial repo (has brainstorm and manifest with warnings) ---
	mustWriteFile(t, filepath.Join(root, "docs", "brainstorms", "idea.md"), "# Idea\n")
	manifest := installstate.InstallManifest{
		Requested: installstate.RequestedState{
			StackPacks:          []installstate.StackPack{installstate.StackPackGeneral, installstate.StackPackRails},
			GstackDirs:          []string{"review"},
			IncludeAgentBrowser: true,
			PresetName:          "Full",
		},
		Outcomes: []installstate.InstallOutcome{
			{Step: "Scaffolding ATV files", Status: installstate.InstallStepDone},
			{Step: "Syncing gstack skills", Status: installstate.InstallStepWarning, Reason: "setup failed"},
			{Step: "Installing agent-browser", Status: installstate.InstallStepFailed, Reason: "npm not found"},
		},
	}
	if err := installstate.WriteManifest(root, manifest); err != nil {
		t.Fatalf("WriteManifest() error = %v", err)
	}

	snapshot, err = installstate.BuildLaunchpadSnapshot(root)
	if err != nil {
		t.Fatalf("partial snapshot error = %v", err)
	}
	if !snapshot.HasManifest {
		t.Fatal("should find manifest")
	}
	if snapshot.OutcomeSummary.Done != 1 || snapshot.OutcomeSummary.Warning != 1 || snapshot.OutcomeSummary.Failed != 1 {
		t.Fatalf("unexpected outcome summary: %+v", snapshot.OutcomeSummary)
	}
	if len(snapshot.Recommendations) == 0 || snapshot.Recommendations[0].ID != "fix-install-issues" {
		t.Fatalf("partial repo should prioritize fix-install-issues, got %+v", snapshot.Recommendations)
	}

	// --- Mature repo (manifest clean, brainstorms, plans with unchecked items) ---
	mustWriteFile(t, filepath.Join(root, "docs", "plans", "work.md"), "status: active\n- [ ] finish feature\n")
	matureManifest := installstate.InstallManifest{
		Requested: installstate.RequestedState{
			StackPacks: []installstate.StackPack{installstate.StackPackGeneral, installstate.StackPackRails},
			GstackDirs: []string{"review"},
			PresetName: "Pro",
		},
		Outcomes: []installstate.InstallOutcome{
			{Step: "Scaffolding ATV files", Status: installstate.InstallStepDone},
			{Step: "Syncing gstack skills", Status: installstate.InstallStepDone},
		},
	}
	if err := installstate.WriteManifest(root, matureManifest); err != nil {
		t.Fatalf("WriteManifest() error = %v", err)
	}

	snapshot, err = installstate.BuildLaunchpadSnapshot(root)
	if err != nil {
		t.Fatalf("mature snapshot error = %v", err)
	}
	if snapshot.OutcomeSummary.Failed != 0 && snapshot.OutcomeSummary.Warning != 0 {
		t.Fatalf("mature repo should have clean outcomes: %+v", snapshot.OutcomeSummary)
	}
	// Should recommend executing the active plan, not fixing install issues
	foundPlan := false
	for _, rec := range snapshot.Recommendations {
		if rec.ID == "execute-active-plan" {
			foundPlan = true
		}
		if rec.ID == "fix-install-issues" {
			t.Fatal("mature repo with clean install should not recommend fix-install-issues")
		}
	}
	if !foundPlan {
		t.Fatalf("mature repo with active plan should recommend execute-active-plan, got %+v", snapshot.Recommendations)
	}

	// Verify determinism
	snap1, _ := installstate.BuildLaunchpadSnapshot(root)
	snap2, _ := installstate.BuildLaunchpadSnapshot(root)
	if len(snap1.Recommendations) != len(snap2.Recommendations) {
		t.Fatal("launchpad should be deterministic")
	}
	for i := range snap1.Recommendations {
		if snap1.Recommendations[i].ID != snap2.Recommendations[i].ID {
			t.Fatalf("recommendation order differs: %+v vs %+v", snap1.Recommendations, snap2.Recommendations)
		}
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// =============================================================================
// Comprehensive E2E: Full guided install lifecycle
// =============================================================================

// TestE2EFullGuidedInstallLifecycle simulates a complete guided install:
//   - detect stack → build filtered catalog → scaffold files → write manifest
//   - verify every expected file is installed in the right location
//   - verify memory index picks up installed intelligence
//   - verify launchpad dashboard reads the correct state
//   - verify all 5 concierge tools return correct structured data
//   - verify concierge degrades gracefully when state is missing
//   - verify determinism across repeated calls
func TestE2EFullGuidedInstallLifecycle(t *testing.T) {
	root := t.TempDir()

	// --- Step 1: Set up a Rails project ---
	mustWriteFile(t, filepath.Join(root, "Gemfile"), `gem "rails"`)
	mustWriteFile(t, filepath.Join(root, "config", "routes.rb"), `Rails.application.routes.draw { }`)

	// --- Step 2: Detect environment ---
	env := detect.DetectEnvironment(root)
	if env.Stack != detect.StackRails {
		t.Fatalf("expected Rails stack, got %s", env.Stack)
	}
	if len(env.DetectedPacks) == 0 {
		t.Fatal("should detect at least one pack")
	}

	// --- Step 3: Build multi-pack filtered catalog (simulating guided "Pro" preset) ---
	selectedPacks := []installstate.StackPack{installstate.StackPackGeneral, installstate.StackPackRails}
	selectedLayers := []string{
		"copilot-instructions", "setup-steps", "mcp-servers",
		"core-skills", "orchestrators", "universal-agents", "stack-agents",
		"file-instructions", "vscode-extensions", "docs-structure",
	}
	catalog := scaffold.BuildFilteredCatalogForPacks(selectedPacks, env.Stack, selectedLayers)

	if len(catalog) == 0 {
		t.Fatal("filtered catalog should produce components")
	}

	// --- Step 4: Write all scaffold files ---
	results := scaffold.WriteAll(root, catalog)
	summary := scaffold.SummarizeResults(results)

	if summary.Created == 0 {
		t.Fatal("scaffold should create files")
	}
	if summary.Failed > 0 {
		t.Fatalf("scaffold should not fail: %d failures", summary.Failed)
	}

	// --- Step 5: Verify core Copilot lifecycle hooks installed ---
	t.Run("hook1_system_instructions", func(t *testing.T) {
		path := filepath.Join(root, ".github", "copilot-instructions.md")
		assertFileExists(t, path)
		content, _ := os.ReadFile(path)
		if !strings.Contains(string(content), "ATV Override Rules") {
			t.Error("instructions should contain ATV Override Rules")
		}
	})

	t.Run("hook2_setup_steps", func(t *testing.T) {
		assertFileExists(t, filepath.Join(root, ".github", "copilot-setup-steps.yml"))
	})

	t.Run("hook3_mcp_config", func(t *testing.T) {
		path := filepath.Join(root, ".github", "copilot-mcp-config.json")
		assertFileExists(t, path)
		content, _ := os.ReadFile(path)
		var parsed map[string]interface{}
		if err := json.Unmarshal(content, &parsed); err != nil {
			t.Fatalf("MCP config should be valid JSON: %v", err)
		}
	})

	t.Run("hook4_skills_installed", func(t *testing.T) {
		skillsDir := filepath.Join(root, ".github", "skills")
		assertDirExists(t, skillsDir)
		// Verify at least core skills exist
		for _, skill := range []string{"brainstorming", "ce-brainstorm", "ce-plan", "ce-review", "ce-work"} {
			skillPath := filepath.Join(skillsDir, skill, "SKILL.md")
			assertFileExists(t, skillPath)
		}
	})

	t.Run("hook5_agents_installed", func(t *testing.T) {
		agentsDir := filepath.Join(root, ".github", "agents")
		assertDirExists(t, agentsDir)
		// Universal agents
		for _, agent := range []string{"code-simplicity-reviewer", "security-sentinel", "architecture-strategist"} {
			assertFileExists(t, filepath.Join(agentsDir, agent+".agent.md"))
		}
		// Stack-specific agents for Rails
		for _, agent := range []string{"kieran-rails-reviewer", "dhh-rails-reviewer", "data-integrity-guardian"} {
			assertFileExists(t, filepath.Join(agentsDir, agent+".agent.md"))
		}
		// Concierge agent should be scaffolded
		assertFileExists(t, filepath.Join(agentsDir, "atv-concierge.agent.md"))
	})

	t.Run("hook6_file_instructions", func(t *testing.T) {
		assertFileExists(t, filepath.Join(root, ".github", "rails.instructions.md"))
	})

	t.Run("vscode_extensions", func(t *testing.T) {
		path := filepath.Join(root, ".vscode", "extensions.json")
		assertFileExists(t, path)
		content, _ := os.ReadFile(path)
		var parsed map[string]interface{}
		if err := json.Unmarshal(content, &parsed); err != nil {
			t.Fatalf("extensions.json should be valid JSON: %v", err)
		}
	})

	t.Run("docs_structure", func(t *testing.T) {
		for _, dir := range []string{"docs/plans", "docs/brainstorms", "docs/solutions"} {
			assertDirExists(t, filepath.Join(root, dir))
		}
	})

	// --- Step 6: Write install manifest (simulating guided mode completion) ---
	substeps := scaffold.ResultsToSubsteps(results)
	manifest := installstate.InstallManifest{
		Requested: installstate.RequestedState{
			StackPacks:          selectedPacks,
			ATVLayers:           selectedLayers,
			GstackDirs:          []string{"review", "plan"},
			GstackRuntime:       false,
			IncludeAgentBrowser: true,
			PresetName:          "Pro",
		},
		Outcomes: []installstate.InstallOutcome{
			{
				Step:     "Scaffolding ATV files",
				Status:   installstate.InstallStepDone,
				Detail:   summary.Detail(),
				Duration: "340ms",
				Substeps: substeps,
			},
			{
				Step:   "Syncing gstack skills",
				Status: installstate.InstallStepWarning,
				Detail: "2 skills synced, markdown-only",
				Reason: "setup script failed, fell back to doc generation",
				Substeps: []installstate.InstallOutcome{
					{Step: "git clone", Status: installstate.InstallStepDone, Detail: "shallow clone"},
					{Step: "runtime build", Status: installstate.InstallStepWarning,
						Reason: "bun not found", SkipReason: installstate.SkipReasonPrereqMissing},
					{Step: "copy skills", Status: installstate.InstallStepDone, Detail: "2 skill dirs"},
				},
			},
			{
				Step:   "Installing agent-browser",
				Status: installstate.InstallStepDone,
				Detail: "CLI ready, skill copied",
				Substeps: []installstate.InstallOutcome{
					{Step: "npm install", Status: installstate.InstallStepDone, Detail: "agent-browser CLI available"},
					{Step: "copy SKILL.md", Status: installstate.InstallStepDone, Detail: "skill registered"},
				},
			},
		},
	}
	manifest.Recommendations = installstate.BuildRecommendations(root, manifest)

	if err := installstate.WriteManifest(root, manifest); err != nil {
		t.Fatalf("WriteManifest() error = %v", err)
	}

	// Simulate gstack staging and agent-browser skill presence (as if install ran)
	os.MkdirAll(filepath.Join(root, ".gstack"), 0o755)
	os.MkdirAll(filepath.Join(root, ".github", "skills", "agent-browser"), 0o755)
	mustWriteFile(t, filepath.Join(root, ".github", "skills", "agent-browser", "SKILL.md"), "# agent-browser")

	// --- Step 7: Verify manifest round-trip ---
	t.Run("manifest_roundtrip", func(t *testing.T) {
		read, err := installstate.ReadManifest(root)
		if err != nil {
			t.Fatalf("ReadManifest() error = %v", err)
		}
		if read.Version != installstate.ManifestVersion {
			t.Fatalf("version = %d, want %d", read.Version, installstate.ManifestVersion)
		}
		if read.RerunPolicy != installstate.RerunPolicyAdditiveOnly {
			t.Fatalf("rerun policy = %q, want additive-only", read.RerunPolicy)
		}
		if read.Requested.PresetName != "Pro" {
			t.Fatalf("preset = %q, want Pro", read.Requested.PresetName)
		}
		if len(read.Outcomes) != 3 {
			t.Fatalf("expected 3 outcomes, got %d", len(read.Outcomes))
		}
		// Substeps survive serialization
		if len(read.Outcomes[0].Substeps) == 0 {
			t.Fatal("scaffold substeps should survive JSON round-trip")
		}
		if len(read.Outcomes[1].Substeps) != 3 {
			t.Fatalf("gstack substeps = %d, want 3", len(read.Outcomes[1].Substeps))
		}
		// Skip reason survives
		if read.Outcomes[1].Substeps[1].SkipReason != installstate.SkipReasonPrereqMissing {
			t.Fatalf("skip reason = %q, want prereq-missing", read.Outcomes[1].Substeps[1].SkipReason)
		}
	})

	// --- Step 8: Verify memory index ---
	t.Run("memory_index_classification", func(t *testing.T) {
		state := installstate.ScanRepoState(root)

		if !state.HasCopilotInstructions {
			t.Error("should detect copilot-instructions.md")
		}
		if state.InstalledAgents == 0 {
			t.Error("should detect installed agents")
		}
		if state.InstalledSkills == 0 {
			t.Error("should detect installed skills")
		}
		if !state.HasGstackStaging {
			t.Error("should detect .gstack/ staging directory")
		}
		if !state.HasAgentBrowserSkill {
			t.Error("should detect agent-browser skill")
		}
	})

	// --- Step 9: Verify launchpad snapshot ---
	t.Run("launchpad_snapshot", func(t *testing.T) {
		snapshot, err := installstate.BuildLaunchpadSnapshot(root)
		if err != nil {
			t.Fatalf("BuildLaunchpadSnapshot() error = %v", err)
		}
		if !snapshot.HasManifest {
			t.Fatal("should find manifest")
		}
		if snapshot.Requested.PresetName != "Pro" {
			t.Fatalf("preset = %q, want Pro", snapshot.Requested.PresetName)
		}
		labels := snapshot.StackPackLabels()
		if len(labels) != 2 || labels[1] != "Rails" {
			t.Fatalf("stack pack labels = %v, want [General, Rails]", labels)
		}
		if !snapshot.HasGstack() {
			t.Error("should report gstack requested")
		}
		if !snapshot.HasAgentBrowser() {
			t.Error("should report agent-browser requested")
		}
		if snapshot.OutcomeSummary.Done != 2 || snapshot.OutcomeSummary.Warning != 1 {
			t.Fatalf("outcomes = %+v, want 2 done + 1 warning", snapshot.OutcomeSummary)
		}
		if len(snapshot.Recommendations) == 0 {
			t.Fatal("should have recommendations")
		}
	})

	// --- Step 10: Verify concierge tools ---
	t.Run("concierge_memory_summary", func(t *testing.T) {
		summary := concierge.GetMemorySummary(root)
		if summary.Status != "ok" {
			t.Fatalf("status = %q, want ok", summary.Status)
		}
		if summary.Manifest == nil {
			t.Fatal("should have manifest info")
		}
		if summary.Manifest.PresetName != "Pro" {
			t.Fatalf("preset = %q, want Pro", summary.Manifest.PresetName)
		}
		if len(summary.Manifest.StackPacks) != 2 {
			t.Fatalf("stack packs = %d, want 2", len(summary.Manifest.StackPacks))
		}
		if summary.Manifest.OutcomeSummary.Done != 2 {
			t.Fatalf("done = %d, want 2", summary.Manifest.OutcomeSummary.Done)
		}
		if !summary.RepoState.HasCopilotInstructions {
			t.Error("memory summary should reflect copilot-instructions presence")
		}
		if summary.RepoState.InstalledAgents == 0 {
			t.Error("memory summary should reflect installed agents")
		}
		if summary.RepoState.InstalledSkills == 0 {
			t.Error("memory summary should reflect installed skills")
		}
		if !summary.RepoState.HasGstackStaging {
			t.Error("memory summary should reflect gstack staging")
		}
		if !summary.RepoState.HasAgentBrowserSkill {
			t.Error("memory summary should reflect agent-browser skill")
		}

		// JSON serialization check
		data, err := json.Marshal(summary)
		if err != nil {
			t.Fatalf("should serialize to JSON: %v", err)
		}
		if !strings.Contains(string(data), `"status":"ok"`) {
			t.Error("JSON output should contain status:ok")
		}
	})

	t.Run("concierge_list_recommendations", func(t *testing.T) {
		list := concierge.ListRecommendations(root)
		if list.Status != "ok" {
			t.Fatalf("status = %q, want ok", list.Status)
		}
		if list.Source != "local-deterministic" {
			t.Fatalf("source = %q, want local-deterministic", list.Source)
		}
		if len(list.Recommendations) == 0 {
			t.Fatal("should have recommendations")
		}
		// First rec should be fix-install-issues (we have a warning)
		if list.Recommendations[0].ID != "fix-install-issues" {
			t.Fatalf("first rec = %q, want fix-install-issues", list.Recommendations[0].ID)
		}

		// Verify determinism: call 10 times, always same order
		for i := 0; i < 10; i++ {
			check := concierge.ListRecommendations(root)
			if len(check.Recommendations) != len(list.Recommendations) {
				t.Fatalf("iteration %d: length differs", i)
			}
			for j := range check.Recommendations {
				if check.Recommendations[j].ID != list.Recommendations[j].ID {
					t.Fatalf("iteration %d, index %d: order differs", i, j)
				}
			}
		}
	})

	t.Run("concierge_explain_recommendation", func(t *testing.T) {
		list := concierge.ListRecommendations(root)
		for _, rec := range list.Recommendations {
			detail := concierge.ExplainRecommendation(root, rec.ID)
			if detail.Status != "ok" {
				t.Fatalf("explain %q status = %q, want ok", rec.ID, detail.Status)
			}
			if detail.Title == "" {
				t.Fatalf("explain %q should have a title", rec.ID)
			}
			if detail.Reason == "" {
				t.Fatalf("explain %q should have a reason", rec.ID)
			}
			if detail.SuggestedCmd == "" {
				t.Fatalf("explain %q should have a suggested command", rec.ID)
			}
			if detail.Priority == 0 {
				t.Fatalf("explain %q should have priority > 0", rec.ID)
			}
		}

		// Non-existent recommendation
		bad := concierge.ExplainRecommendation(root, "invented-by-assistant")
		if bad.Status != "not-found" {
			t.Fatalf("invented rec status = %q, want not-found", bad.Status)
		}
	})

	t.Run("concierge_open_artifact", func(t *testing.T) {
		// Existing artifacts
		for _, tc := range []struct {
			name   string
			exists bool
			typ    string
		}{
			{"manifest", true, "file"},
			{"instructions", true, "file"},
			{"brainstorms", true, "directory"},
			{"plans", true, "directory"},
			{"solutions", true, "directory"},
			{"agents", true, "directory"},
			{"skills", true, "directory"},
		} {
			info := concierge.OpenArtifact(root, tc.name)
			if info.Status != "ok" {
				t.Fatalf("artifact %q status = %q, want ok", tc.name, info.Status)
			}
			if info.Exists != tc.exists {
				t.Fatalf("artifact %q exists = %v, want %v", tc.name, info.Exists, tc.exists)
			}
			if info.Type != tc.typ {
				t.Fatalf("artifact %q type = %q, want %q", tc.name, info.Type, tc.typ)
			}
			if info.Path == "" {
				t.Fatalf("artifact %q should have a path", tc.name)
			}
		}

		// Unknown artifact
		unknown := concierge.OpenArtifact(root, "electron-dashboard")
		if unknown.Status != "unknown" {
			t.Fatalf("unknown artifact status = %q, want unknown", unknown.Status)
		}
	})

	t.Run("concierge_run_suggested_action", func(t *testing.T) {
		// All known recommendation IDs should have a suggested action
		knownActions := []string{
			"fix-install-issues", "start-brainstorm", "turn-brainstorm-into-plan",
			"execute-active-plan", "compound-learnings", "start-gstack-sprint", "browser-check",
		}
		for _, id := range knownActions {
			result := concierge.RunSuggestedAction(root, id)
			if result.Status != "ready" {
				t.Fatalf("action %q status = %q, want ready", id, result.Status)
			}
			if result.Message == "" {
				t.Fatalf("action %q should have a message", id)
			}
			// Should contain "present this to the user" — never auto-execute
			if !strings.Contains(result.Message, "present this to the user") {
				t.Fatalf("action %q message should prevent silent execution: %q", id, result.Message)
			}
		}

		// Invented action should be rejected
		invented := concierge.RunSuggestedAction(root, "delete-all-files")
		if invented.Status != "unknown" {
			t.Fatalf("invented action status = %q, want unknown", invented.Status)
		}
		if !strings.Contains(invented.Message, "cannot invent actions") {
			t.Fatalf("invented action should explain refusal: %q", invented.Message)
		}
	})

	// --- Step 11: Verify concierge parity with raw launchpad ---
	t.Run("concierge_parity_with_launchpad", func(t *testing.T) {
		snapshot, _ := installstate.BuildLaunchpadSnapshot(root)
		list := concierge.ListRecommendations(root)

		if len(list.Recommendations) != len(snapshot.Recommendations) {
			t.Fatalf("concierge recs (%d) != launchpad recs (%d)",
				len(list.Recommendations), len(snapshot.Recommendations))
		}
		for i := range list.Recommendations {
			if list.Recommendations[i].ID != snapshot.Recommendations[i].ID {
				t.Fatalf("rec %d: concierge %q != launchpad %q",
					i, list.Recommendations[i].ID, snapshot.Recommendations[i].ID)
			}
			if list.Recommendations[i].Priority != snapshot.Recommendations[i].Priority {
				t.Fatalf("rec %d priority: concierge %d != launchpad %d",
					i, list.Recommendations[i].Priority, snapshot.Recommendations[i].Priority)
			}
		}
	})
}

// TestE2EConciergeDegradedBehavior tests all concierge tools in degraded states
func TestE2EConciergeDegradedBehavior(t *testing.T) {
	t.Run("empty_repo_no_manifest", func(t *testing.T) {
		root := t.TempDir()

		// Memory summary should degrade gracefully
		summary := concierge.GetMemorySummary(root)
		if summary.Status != "no-manifest" {
			t.Fatalf("status = %q, want no-manifest", summary.Status)
		}
		if summary.Manifest != nil {
			t.Fatal("should not have manifest info")
		}
		if summary.Message == "" {
			t.Fatal("should have a helpful degraded message")
		}
		// RepoState should still work (all zeros)
		if summary.RepoState.BrainstormCount != 0 {
			t.Fatalf("brainstorms = %d, want 0", summary.RepoState.BrainstormCount)
		}

		// Recommendations still work without manifest
		list := concierge.ListRecommendations(root)
		if list.Status != "ok" {
			t.Fatalf("status = %q, want ok", list.Status)
		}
		if len(list.Recommendations) == 0 {
			t.Fatal("should still recommend start-brainstorm")
		}
		if list.Recommendations[0].ID != "start-brainstorm" {
			t.Fatalf("first rec = %q, want start-brainstorm", list.Recommendations[0].ID)
		}

		// Explain still works for active recommendations
		detail := concierge.ExplainRecommendation(root, "start-brainstorm")
		if detail.Status != "ok" {
			t.Fatalf("explain status = %q, want ok", detail.Status)
		}

		// Artifacts resolve even if they don't exist yet
		for _, name := range []string{"manifest", "instructions", "brainstorms", "plans"} {
			info := concierge.OpenArtifact(root, name)
			if info.Status != "ok" {
				t.Fatalf("artifact %q status = %q, want ok", name, info.Status)
			}
			// manifest and instructions don't exist, dirs may or may not
		}
	})

	t.Run("repo_with_docs_but_no_manifest", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "docs", "brainstorms", "idea.md"), "# Idea\n")
		mustWriteFile(t, filepath.Join(root, "docs", "plans", "work.md"), "status: active\n- [ ] todo\n")

		summary := concierge.GetMemorySummary(root)
		if summary.Status != "no-manifest" {
			t.Fatalf("status = %q, want no-manifest", summary.Status)
		}
		// But repo state should still be scanned
		if summary.RepoState.BrainstormCount != 1 {
			t.Fatalf("brainstorms = %d, want 1", summary.RepoState.BrainstormCount)
		}
		if summary.RepoState.PlanCount != 1 {
			t.Fatalf("plans = %d, want 1", summary.RepoState.PlanCount)
		}
		if !summary.RepoState.HasUncheckedPlan {
			t.Fatal("should detect unchecked plan")
		}

		// Recommendations should reflect repo state even without manifest
		list := concierge.ListRecommendations(root)
		if list.Recommendations[0].ID != "execute-active-plan" {
			t.Fatalf("first rec = %q, want execute-active-plan", list.Recommendations[0].ID)
		}
	})

	t.Run("corrupt_manifest_file", func(t *testing.T) {
		root := t.TempDir()
		os.MkdirAll(filepath.Join(root, ".atv"), 0o755)
		os.WriteFile(filepath.Join(root, ".atv", "install-manifest.json"), []byte("not json"), 0o644)

		summary := concierge.GetMemorySummary(root)
		if summary.Status != "error" {
			t.Fatalf("status = %q, want error", summary.Status)
		}
		if summary.Message == "" {
			t.Fatal("should explain the error")
		}
	})
}

// TestE2EMultiStackDeterminism verifies that multi-stack selection is deterministic
// regardless of pack selection order
func TestE2EMultiStackDeterminism(t *testing.T) {
	// Order A: General, TypeScript, Rails
	catalogA := scaffold.BuildFilteredCatalogForPacks(
		[]installstate.StackPack{installstate.StackPackGeneral, installstate.StackPackTypeScript, installstate.StackPackRails},
		detect.StackRails,
		[]string{"core-skills", "universal-agents", "stack-agents", "file-instructions"},
	)

	// Order B: Rails, General, TypeScript (reversed)
	catalogB := scaffold.BuildFilteredCatalogForPacks(
		[]installstate.StackPack{installstate.StackPackRails, installstate.StackPackGeneral, installstate.StackPackTypeScript},
		detect.StackRails,
		[]string{"core-skills", "universal-agents", "stack-agents", "file-instructions"},
	)

	if len(catalogA) != len(catalogB) {
		t.Fatalf("catalog lengths differ: %d vs %d", len(catalogA), len(catalogB))
	}

	pathsA := make(map[string]bool, len(catalogA))
	for _, c := range catalogA {
		pathsA[c.Path] = true
	}

	for _, c := range catalogB {
		if !pathsA[c.Path] {
			t.Fatalf("catalog B has path %q not in catalog A", c.Path)
		}
	}

	// Write both catalogs and verify identical output
	rootA := t.TempDir()
	rootB := t.TempDir()

	scaffold.WriteAll(rootA, catalogA)
	scaffold.WriteAll(rootB, catalogB)

	// Verify same files exist
	var filesA, filesB []string
	filepath.Walk(rootA, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			rel, _ := filepath.Rel(rootA, path)
			filesA = append(filesA, filepath.ToSlash(rel))
		}
		return nil
	})
	filepath.Walk(rootB, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			rel, _ := filepath.Rel(rootB, path)
			filesB = append(filesB, filepath.ToSlash(rel))
		}
		return nil
	})

	if len(filesA) != len(filesB) {
		t.Fatalf("file count differs: %d vs %d", len(filesA), len(filesB))
	}

	setA := make(map[string]bool, len(filesA))
	for _, f := range filesA {
		setA[f] = true
	}
	for _, f := range filesB {
		if !setA[f] {
			t.Fatalf("file %q in B but not A", f)
		}
	}
}

// TestE2EConciergeAgentTemplateContent verifies the atv-concierge agent template
// is correct for Copilot discovery
func TestE2EConciergeAgentTemplateContent(t *testing.T) {
	root := t.TempDir()
	catalog := scaffold.BuildCatalog(detect.StackGeneral)
	scaffold.WriteAll(root, catalog)

	agentPath := filepath.Join(root, ".github", "agents", "atv-concierge.agent.md")
	assertFileExists(t, agentPath)

	content, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("read agent template: %v", err)
	}
	text := string(content)

	// Verify frontmatter
	if !strings.Contains(text, "description:") {
		t.Error("agent template should have description frontmatter")
	}

	// Verify all 5 concierge tool commands are documented
	for _, cmd := range []string{
		"memory-summary",
		"list-recommendations",
		"explain-recommendation",
		"open-artifact",
		"run-suggested-action",
	} {
		if !strings.Contains(text, cmd) {
			t.Errorf("agent template should document %q command", cmd)
		}
	}

	// Verify key design principles
	if !strings.Contains(text, "deterministic") {
		t.Error("agent template should mention deterministic recommendations")
	}
	if !strings.Contains(text, "never execute without asking") || !strings.Contains(text, "never own") {
		t.Error("agent template should enforce assistant-as-secondary principle")
	}

	// Verify degraded behavior section
	if !strings.Contains(text, "No manifest") || !strings.Contains(text, "Offline") {
		t.Error("agent template should document graceful degradation")
	}
}
