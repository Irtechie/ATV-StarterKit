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

// copyGeneratedSkills copies gstack-*/SKILL.md from .agents/skills/ to .github/skills/.
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
			// Only copy gstack-prefixed skill dirs
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
		if len(dirs) > 0 {
			return true, dirs
		}
	}

	// Fallback: copy raw skill dirs directly from gstack root
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
	return len(dirs) > 0, dirs
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
