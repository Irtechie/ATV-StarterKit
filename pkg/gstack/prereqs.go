package gstack

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Prerequisites holds the detected state of required tools.
type Prerequisites struct {
	HasGit      bool
	HasBun      bool
	HasNode     bool
	GitVersion  string
	BunVersion  string
	NodeVersion string
}

// DetectPrerequisites checks for git, bun, and node availability.
func DetectPrerequisites() Prerequisites {
	p := Prerequisites{}

	if path, err := exec.LookPath("git"); err == nil && path != "" {
		p.HasGit = true
		p.GitVersion = getVersion("git", "--version")
	}

	bunCmd := "bun"
	if runtime.GOOS == "windows" {
		// On Windows, bun may be installed via npm or scoop
		if _, err := exec.LookPath("bun"); err != nil {
			if _, err := exec.LookPath("bun.exe"); err == nil {
				bunCmd = "bun.exe"
			}
		}
	}
	if path, err := exec.LookPath(bunCmd); err == nil && path != "" {
		p.HasBun = true
		p.BunVersion = getVersion(bunCmd, "--version")
	}

	if path, err := exec.LookPath("node"); err == nil && path != "" {
		p.HasNode = true
		p.NodeVersion = getVersion("node", "--version")
	}

	return p
}

// RuntimeAvailable returns true if either Bun or Node is available for building gstack.
func (p Prerequisites) RuntimeAvailable() bool {
	return p.HasBun || p.HasNode
}

// Summary returns a human-readable string of detected prerequisites.
func (p Prerequisites) Summary() string {
	var parts []string
	if p.HasGit {
		parts = append(parts, fmt.Sprintf("git %s", p.GitVersion))
	} else {
		parts = append(parts, "git (missing)")
	}
	if p.HasBun {
		parts = append(parts, fmt.Sprintf("bun %s", p.BunVersion))
	} else {
		parts = append(parts, "bun (missing)")
	}
	if p.HasNode {
		parts = append(parts, fmt.Sprintf("node %s", p.NodeVersion))
	}
	return strings.Join(parts, ", ")
}

func getVersion(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return "unknown"
	}
	version := strings.TrimSpace(string(out))
	// git --version returns "git version 2.x.x", extract just the version
	if strings.HasPrefix(version, "git version ") {
		version = strings.TrimPrefix(version, "git version ")
	}
	return version
}
