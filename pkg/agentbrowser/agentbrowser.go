package agentbrowser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	// SkillRepo is the GitHub raw URL for the agent-browser SKILL.md.
	SkillRepo = "https://raw.githubusercontent.com/vercel-labs/agent-browser/main/skills/agent-browser/SKILL.md"
	// NpmPackage is the npm package name.
	NpmPackage = "agent-browser"
)

// InstallResult reports what happened during agent-browser setup.
type InstallResult struct {
	Installed   bool   // npm package was installed
	SkillCopied bool   // SKILL.md was copied to .github/skills/
	Warning     string // non-fatal warning
	Error       error
}

// IsInstalled checks if agent-browser CLI is available on PATH.
func IsInstalled() bool {
	_, err := exec.LookPath("agent-browser")
	return err == nil
}

// Install installs agent-browser globally via npm, downloads Chrome, and copies the SKILL.md for Copilot discovery.
func Install(projectDir string) *InstallResult {
	result := &InstallResult{}

	// Step 1: Install via npm if not already installed
	if !IsInstalled() {
		if err := npmInstall(); err != nil {
			result.Warning = fmt.Sprintf("npm install failed (%v), skipping agent-browser binary. SKILL.md still usable", err)
		} else {
			result.Installed = true
		}
	} else {
		result.Installed = true
	}

	// Step 2: Download Chrome for Testing (agent-browser install)
	if result.Installed {
		if err := installChrome(); err != nil {
			result.Warning = fmt.Sprintf("Chrome download failed (%v), run 'agent-browser install' manually", err)
		}
	}

	// Step 3: Fetch and copy SKILL.md to .github/skills/agent-browser/
	if err := copySkill(projectDir); err != nil {
		if result.Warning != "" {
			result.Error = fmt.Errorf("agent-browser setup failed: %w", err)
		} else {
			result.Warning = fmt.Sprintf("failed to copy SKILL.md: %v", err)
		}
	} else {
		result.SkillCopied = true
	}

	return result
}

// npmInstall runs npm install -g agent-browser.
func npmInstall() error {
	if _, err := exec.LookPath("npm"); err != nil {
		return fmt.Errorf("npm not found; install Node.js (https://nodejs.org)")
	}

	cmd := exec.Command("npm", "install", "-g", NpmPackage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("npm install -g agent-browser failed: %w", err)
	}
	return nil
}

// installChrome runs "agent-browser install" to download Chrome for Testing.
func installChrome() error {
	cmd := exec.Command("agent-browser", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("agent-browser install failed: %w", err)
	}
	return nil
}

// copySkill fetches the SKILL.md and writes it to .github/skills/agent-browser/SKILL.md.
// Uses npx skills add if available, otherwise fetches directly.
func copySkill(projectDir string) error {
	destDir := filepath.Join(projectDir, ".github", "skills", "agent-browser")
	destFile := filepath.Join(destDir, "SKILL.md")

	// Skip if already exists
	if _, err := os.Stat(destFile); err == nil {
		return nil
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Try fetching via curl/wget
	if err := fetchSkillFile(destFile); err != nil {
		// Fallback: write a minimal SKILL.md that references the CLI
		return writeMinimalSkill(destFile)
	}

	return nil
}

// fetchSkillFile downloads the SKILL.md from GitHub.
func fetchSkillFile(destFile string) error {
	// Try curl first (available on all modern OS)
	if curlPath, err := exec.LookPath("curl"); err == nil {
		cmd := exec.Command(curlPath, "-fsSL", "-o", destFile, SkillRepo)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Try PowerShell on Windows
	if pwshPath, err := exec.LookPath("powershell"); err == nil {
		cmd := exec.Command(pwshPath, "-Command",
			fmt.Sprintf("Invoke-WebRequest -Uri '%s' -OutFile '%s'", SkillRepo, destFile))
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("could not download SKILL.md (no curl or powershell)")
}

// writeMinimalSkill writes a fallback SKILL.md when download fails.
func writeMinimalSkill(destFile string) error {
	content := `---
name: agent-browser
description: Browser automation CLI for AI agents. Use when asked to open websites, fill forms, click buttons, take screenshots, scrape data, test web apps, or automate browser tasks.
---

# Browser Automation with agent-browser

Install: ` + "`npm install -g agent-browser && agent-browser install`" + `

Core workflow:
1. ` + "`agent-browser open <url>`" + ` — Navigate to page
2. ` + "`agent-browser snapshot -i`" + ` — Get interactive elements with refs (@e1, @e2)
3. ` + "`agent-browser click @e1`" + ` / ` + "`fill @e2 \"text\"`" + ` — Interact using refs
4. Re-snapshot after page changes

Run ` + "`agent-browser --help`" + ` for full command reference.
`
	return os.WriteFile(destFile, []byte(content), 0644)
}
