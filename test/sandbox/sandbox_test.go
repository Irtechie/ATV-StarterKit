package sandbox

import (
	"os"
	"path/filepath"
	"testing"

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
