package gstack

import (
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
	// ModeFullRuntime clones and runs bun run build for browser skills.
	ModeFullRuntime
)

// InstallResult reports what happened during installation.
type InstallResult struct {
	Cloned    bool
	Built     bool
	SkillDirs []string
	Mode      InstallMode
	Error     error
	Warning   string // non-fatal warning (e.g., "bun not found, using markdown-only")
}

// Install orchestrates the full gstack installation: clone → strip .git → optionally build.
func Install(targetDir string, mode InstallMode) *InstallResult {
	result := &InstallResult{Mode: mode}

	// Idempotent: skip if already installed
	if _, err := os.Stat(targetDir); err == nil {
		// Check if it looks like a valid gstack install
		if _, err := os.Stat(filepath.Join(targetDir, "SKILL.md")); err == nil {
			result.Warning = "gstack already installed, skipping"
			return result
		}
	}

	// Clone
	if err := Clone(targetDir); err != nil {
		result.Error = fmt.Errorf("failed to clone gstack: %w", err)
		return result
	}
	result.Cloned = true

	// Strip .git
	if err := StripGit(targetDir); err != nil {
		result.Error = fmt.Errorf("failed to strip .git: %w", err)
		cleanup(targetDir)
		return result
	}

	// Build if full runtime requested
	if mode == ModeFullRuntime {
		if err := Build(targetDir); err != nil {
			result.Warning = fmt.Sprintf("build failed (%v), SKILL.md files still usable", err)
		} else {
			result.Built = true
		}
	}

	// Collect installed skill directories
	result.SkillDirs = listSkillDirs(targetDir)

	return result
}

// Clone performs a shallow git clone of gstack to the target directory.
func Clone(targetDir string) error {
	// Ensure parent directory exists
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

// StripGit removes the .git directory from a cloned repo.
func StripGit(targetDir string) error {
	gitDir := filepath.Join(targetDir, ".git")
	if err := os.RemoveAll(gitDir); err != nil {
		return fmt.Errorf("failed to remove .git directory: %w", err)
	}
	return nil
}

// Build runs bun run build in the gstack directory. Falls back to node on Windows if bun unavailable.
func Build(targetDir string) error {
	var cmd *exec.Cmd

	if _, err := exec.LookPath("bun"); err == nil {
		cmd = exec.Command("bun", "run", "build")
	} else if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("node"); err == nil {
			cmd = exec.Command("node", "node_modules/.bin/build")
		} else {
			return fmt.Errorf("neither bun nor node found; install bun (https://bun.sh) for full runtime")
		}
	} else {
		return fmt.Errorf("bun not found; install bun (https://bun.sh) for full runtime")
	}

	cmd.Dir = targetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return nil
}

// listSkillDirs returns the list of skill directory names in the gstack install.
func listSkillDirs(targetDir string) []string {
	var dirs []string
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return dirs
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		// Check if directory contains a SKILL.md
		skillPath := filepath.Join(targetDir, e.Name(), "SKILL.md")
		if _, err := os.Stat(skillPath); err == nil {
			dirs = append(dirs, e.Name())
		}
	}
	return dirs
}

// cleanup removes a partially-installed gstack directory.
func cleanup(targetDir string) {
	_ = os.RemoveAll(targetDir)
}
