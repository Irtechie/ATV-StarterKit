package gstack

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	// GstackRepo is the GitHub repository URL for gstack.
	GstackRepo = "https://github.com/garrytan/gstack.git"
)

// InstallMode controls how much of gstack to install.
type InstallMode int

const (
	// ModeMarkdownOnly copies SKILL.md files only, no binary build.
	ModeMarkdownOnly InstallMode = iota
	// ModeFullRuntime clones and runs ./setup for browser skills + binary.
	ModeFullRuntime
)

// InstallResult reports what happened during installation.
type InstallResult struct {
	Cloned    bool
	Built     bool
	Copied    bool
	SkillDirs []string
	Mode      InstallMode
	Error     error
	Warning   string
}

// Install clones gstack, runs its own ./setup, then copies generated skills to .github/skills/.
//
// Flow:
//  1. Clone gstack to <projectDir>/.gstack/ (gitignored staging area)
//  2. Run gstack's own ./setup script (handles deps, build, skill generation)
//  3. Copy generated skills from .gstack/.agents/skills/gstack-* → .github/skills/gstack-*
func Install(projectDir string, mode InstallMode) *InstallResult {
	result := &InstallResult{Mode: mode}
	stagingDir := filepath.Join(projectDir, ".gstack")
	skillsTargetDir := filepath.Join(projectDir, ".github", "skills")

	// Idempotent: skip if .gstack/ already exists with SKILL.md
	if _, err := os.Stat(filepath.Join(stagingDir, "SKILL.md")); err == nil {
		// Already cloned — just re-copy skills
		copied, dirs := copyGeneratedSkills(stagingDir, skillsTargetDir)
		result.Copied = copied
		result.SkillDirs = dirs
		result.Warning = "gstack already cloned, re-synced skills"
		return result
	}

	// Step 1: Clone to .gstack/
	if err := clone(stagingDir); err != nil {
		result.Error = fmt.Errorf("failed to clone gstack: %w", err)
		return result
	}
	result.Cloned = true

	// Step 2: Run gstack's ./setup (full runtime) or just generate docs (markdown-only)
	if mode == ModeFullRuntime {
		if err := runSetup(stagingDir); err != nil {
			result.Warning = fmt.Sprintf("./setup failed (%v), falling back to doc generation only", err)
			// Fallback: try generating docs without full build
			_ = generateDocs(stagingDir)
		} else {
			result.Built = true
		}
	} else {
		// Markdown-only: just install deps and generate skill docs
		if err := generateDocs(stagingDir); err != nil {
			result.Warning = fmt.Sprintf("doc generation failed (%v), copying raw SKILL.md files", err)
		}
	}

	// Step 3: Copy generated skills to .github/skills/
	copied, dirs := copyGeneratedSkills(stagingDir, skillsTargetDir)
	result.Copied = copied
	result.SkillDirs = dirs

	return result
}

// clone performs a shallow git clone of gstack.
func clone(targetDir string) error {
	if err := os.MkdirAll(filepath.Dir(targetDir), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	cmd := exec.Command("git", "clone", "--single-branch", "--depth", "1", GstackRepo, targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed (check network connection): %w", err)
	}
	return nil
}

// runSetup runs gstack's own ./setup script which handles deps, build, and skill generation.
func runSetup(gstackDir string) error {
	// Fix CRLF in scripts on Windows before running
	if runtime.GOOS == "windows" {
		fixCRLFInScripts(gstackDir)
	}

	setupPath := filepath.Join(gstackDir, "setup")

	// Find bash — on Windows it comes with Git
	bashPath := findBash()
	if bashPath == "" {
		return fmt.Errorf("bash not found; gstack's setup requires bash (included with Git for Windows)")
	}

	cmd := exec.Command(bashPath, setupPath, "--host", "codex")
	cmd.Dir = gstackDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "NONINTERACTIVE=1")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gstack setup failed: %w", err)
	}
	return nil
}

// generateDocs runs bun install + bun run gen:skill-docs without the full build.
func generateDocs(gstackDir string) error {
	if _, err := exec.LookPath("bun"); err != nil {
		return fmt.Errorf("bun not found; install bun (https://bun.sh)")
	}

	installCmd := exec.Command("bun", "install")
	installCmd.Dir = gstackDir
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("bun install failed: %w", err)
	}

	genCmd := exec.Command("bun", "run", "gen:skill-docs", "--host", "codex")
	genCmd.Dir = gstackDir
	genCmd.Stdout = os.Stdout
	genCmd.Stderr = os.Stderr
	if err := genCmd.Run(); err != nil {
		return fmt.Errorf("gen:skill-docs failed: %w", err)
	}
	return nil
}

// copyGeneratedSkills copies gstack-*/SKILL.md from .agents/skills/ to .github/skills/,
// and creates a runtime sidecar at .github/skills/gstack/ with bin/, browse/dist/, ETHOS.md.
// Falls back to copying raw skill dirs if .agents/ doesn't exist.
func copyGeneratedSkills(gstackDir, targetSkillsDir string) (bool, []string) {
	var dirs []string

	// Try generated .agents/skills/ first (produced by ./setup or gen:skill-docs)
	agentsDir := filepath.Join(gstackDir, ".agents", "skills")
	if entries, err := os.ReadDir(agentsDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			name := e.Name()
			// Only copy gstack-prefixed skill dirs (not the root "gstack" sidecar)
			if len(name) < 7 || name[:7] != "gstack-" {
				continue
			}
			srcSkill := filepath.Join(agentsDir, name, "SKILL.md")
			if _, err := os.Stat(srcSkill); err != nil {
				continue
			}
			destDir := filepath.Join(targetSkillsDir, name)
			if err := os.MkdirAll(destDir, 0755); err != nil {
				continue
			}
			data, err := os.ReadFile(srcSkill)
			if err != nil {
				continue
			}
			if err := os.WriteFile(filepath.Join(destDir, "SKILL.md"), data, 0644); err != nil {
				continue
			}
			dirs = append(dirs, name)
		}
	}

	// Fallback: copy raw skill dirs directly from gstack root
	if len(dirs) == 0 {
		entries, err := os.ReadDir(gstackDir)
		if err != nil {
			return false, nil
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			name := e.Name()
			srcSkill := filepath.Join(gstackDir, name, "SKILL.md")
			if _, err := os.Stat(srcSkill); err != nil {
				continue
			}
			destName := "gstack-" + name
			destDir := filepath.Join(targetSkillsDir, destName)
			if err := os.MkdirAll(destDir, 0755); err != nil {
				continue
			}
			data, err := os.ReadFile(srcSkill)
			if err != nil {
				continue
			}
			if err := os.WriteFile(filepath.Join(destDir, "SKILL.md"), data, 0644); err != nil {
				continue
			}
			dirs = append(dirs, destName)
		}
	}

	// Create runtime sidecar at .github/skills/gstack/
	// Skills reference $GSTACK_ROOT for binaries, ETHOS.md, and review assets
	createSidecar(gstackDir, targetSkillsDir)

	return len(dirs) > 0, dirs
}

// createSidecar creates .github/skills/gstack/ with the runtime assets that skills reference.
// This mirrors what gstack's ./setup creates for Codex at .agents/skills/gstack/.
func createSidecar(gstackDir, targetSkillsDir string) {
	sidecar := filepath.Join(targetSkillsDir, "gstack")
	_ = os.MkdirAll(sidecar, 0755)

	// Copy root SKILL.md (the meta-skill for browse)
	copyFileIfExists(filepath.Join(gstackDir, "SKILL.md"), filepath.Join(sidecar, "SKILL.md"))

	// Copy ETHOS.md (referenced by "Search Before Building" in all skill preambles)
	copyFileIfExists(filepath.Join(gstackDir, "ETHOS.md"), filepath.Join(sidecar, "ETHOS.md"))

	// Copy AGENTS.md
	copyFileIfExists(filepath.Join(gstackDir, "AGENTS.md"), filepath.Join(sidecar, "AGENTS.md"))

	// Symlink or copy bin/ (gstack-config, gstack-update-check, etc.)
	linkOrCopyDir(filepath.Join(gstackDir, "bin"), filepath.Join(sidecar, "bin"))

	// Symlink or copy browse/dist/ (compiled browse binary)
	browseDist := filepath.Join(gstackDir, "browse", "dist")
	if _, err := os.Stat(browseDist); err == nil {
		_ = os.MkdirAll(filepath.Join(sidecar, "browse"), 0755)
		linkOrCopyDir(browseDist, filepath.Join(sidecar, "browse", "dist"))
	}

	// Symlink or copy browse/bin/ (helper scripts)
	browseBin := filepath.Join(gstackDir, "browse", "bin")
	if _, err := os.Stat(browseBin); err == nil {
		_ = os.MkdirAll(filepath.Join(sidecar, "browse"), 0755)
		linkOrCopyDir(browseBin, filepath.Join(sidecar, "browse", "bin"))
	}

	// Symlink or copy design/dist/ (design binary)
	designDist := filepath.Join(gstackDir, "design", "dist")
	if _, err := os.Stat(designDist); err == nil {
		_ = os.MkdirAll(filepath.Join(sidecar, "design"), 0755)
		linkOrCopyDir(designDist, filepath.Join(sidecar, "design", "dist"))
	}

	// Copy review runtime assets (checklists, not the SKILL.md)
	reviewAssets := []string{"checklist.md", "design-checklist.md", "greptile-triage.md", "TODOS-format.md"}
	reviewDir := filepath.Join(sidecar, "review")
	_ = os.MkdirAll(reviewDir, 0755)
	for _, f := range reviewAssets {
		copyFileIfExists(filepath.Join(gstackDir, "review", f), filepath.Join(reviewDir, f))
	}
}

func copyFileIfExists(src, dst string) {
	data, err := os.ReadFile(src)
	if err != nil {
		return
	}
	_ = os.WriteFile(dst, data, 0644)
}

// linkOrCopyDir creates a symlink from src to dst. Falls back to copy on Windows if symlinks fail.
func linkOrCopyDir(src, dst string) {
	// Remove existing target
	os.RemoveAll(dst)

	// Try symlink first (works on most systems, needs dev mode on Windows)
	if err := os.Symlink(src, dst); err == nil {
		return
	}

	// Fallback: copy the directory
	copyDir(src, dst)
}

func copyDir(src, dst string) {
	os.MkdirAll(dst, 0755)
	entries, err := os.ReadDir(src)
	if err != nil {
		return
	}
	for _, e := range entries {
		srcPath := filepath.Join(src, e.Name())
		dstPath := filepath.Join(dst, e.Name())
		if e.IsDir() {
			copyDir(srcPath, dstPath)
		} else {
			copyFileIfExists(srcPath, dstPath)
		}
	}
}

// findBash locates bash executable — on Windows it ships with Git.
func findBash() string {
	if path, err := exec.LookPath("bash"); err == nil {
		return path
	}
	// Common Git for Windows locations
	if runtime.GOOS == "windows" {
		candidates := []string{
			filepath.Join(os.Getenv("ProgramFiles"), "Git", "bin", "bash.exe"),
			filepath.Join(os.Getenv("ProgramFiles(x86)"), "Git", "bin", "bash.exe"),
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Git", "bin", "bash.exe"),
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				return c
			}
		}
	}
	return ""
}

// fixCRLFInScripts converts \r\n to \n in shell scripts to fix Windows git clone issues.
func fixCRLFInScripts(gstackDir string) {
	// Fix the main setup script
	setupPath := filepath.Join(gstackDir, "setup")
	fixCRLFInFile(setupPath)

	// Fix scripts in browse/scripts/
	scriptsDir := filepath.Join(gstackDir, "browse", "scripts")
	entries, err := os.ReadDir(scriptsDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() {
			fixCRLFInFile(filepath.Join(scriptsDir, e.Name()))
		}
	}

	// Fix scripts in bin/
	binDir := filepath.Join(gstackDir, "bin")
	entries, err = os.ReadDir(binDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() {
			fixCRLFInFile(filepath.Join(binDir, e.Name()))
		}
	}
}

func fixCRLFInFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	fixed := bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	if !bytes.Equal(data, fixed) {
		_ = os.WriteFile(path, fixed, 0644)
	}
}
