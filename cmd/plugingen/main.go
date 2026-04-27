// Command plugingen regenerates the ATV Copilot CLI plugin marketplace
// from the scaffold templates.
//
// Usage:
//
//	go run ./cmd/plugingen          # regenerate plugins/ and .github/plugin/marketplace.json
//	go run ./cmd/plugingen -check   # CI mode: exit 1 if regeneration would change anything
//
// The generator must be run from any directory inside the repository
// (it walks up from the working directory to find go.mod). The kit
// version is read from the VERSION file at the repo root.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/plugingen"
)

func main() {
	check := flag.Bool("check", false, "Verify the committed plugins/ tree matches what the generator would produce. Exit 1 if it doesn't. Used by CI.")
	flag.Parse()

	repoRoot, err := findRepoRoot()
	if err != nil {
		exit("locate repo root: %v", err)
	}
	versionPath := filepath.Join(repoRoot, "VERSION")
	versionBytes, err := os.ReadFile(versionPath)
	if err != nil {
		exit("read %s: %v", versionPath, err)
	}
	version := strings.TrimSpace(string(versionBytes))
	if version == "" {
		exit("VERSION file is empty")
	}

	cfg := plugingen.Config{RepoRoot: repoRoot, KitVersion: version}

	if *check {
		if err := plugingen.CheckClean(cfg); err != nil {
			exit("%v", err)
		}
		fmt.Printf("plugingen: plugin tree is in sync with templates (kit version %s)\n", version)
		return
	}

	if err := plugingen.Generate(cfg); err != nil {
		exit("%v", err)
	}
	fmt.Printf("plugingen: regenerated plugins/ and .github/plugin/marketplace.json (kit version %s)\n", version)
}

// findRepoRoot walks up from the working directory looking for go.mod.
// We don't rely on git in case the generator is invoked from a tarball
// or a CI checkout without .git.
func findRepoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no go.mod found at or above %s", wd)
		}
		dir = parent
	}
}

func exit(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "plugingen: "+format+"\n", args...)
	os.Exit(1)
}
