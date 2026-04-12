package monitor

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
)

// ComputeDrift compares the install-time checksums stored in the manifest
// against the current state of files on disk.
//
// Statuses:
//   - Missing:      in manifest but not on disk
//   - UserModified: on disk but checksum differs from install time
//   - Extra:        on disk in .github/ or .copilot/ but not in manifest
//
// Files matching .atv/drift-ignore patterns are excluded.
func ComputeDrift(root string, manifest installstate.InstallManifest) []DriftEntry {
	if len(manifest.FileChecksums) == 0 {
		return nil
	}

	ignorePatterns := loadDriftIgnore(root)
	var entries []DriftEntry

	for relPath, installHash := range manifest.FileChecksums {
		if isIgnoredByDrift(relPath, ignorePatterns) {
			continue
		}

		absPath := filepath.Join(root, relPath)
		diskHash := hashFile(absPath)

		switch {
		case diskHash == "":
			entries = append(entries, DriftEntry{
				Path:        relPath,
				Status:      DriftMissing,
				InstallHash: installHash,
			})
		case diskHash != installHash:
			entries = append(entries, DriftEntry{
				Path:        relPath,
				Status:      DriftUserModified,
				DiskHash:    diskHash,
				InstallHash: installHash,
			})
		// diskHash == installHash → no drift, skip
		}
	}

	// Sort for deterministic output
	sortDriftEntries(entries)
	return entries
}

// hashFile computes the SHA-256 of a file on disk. Returns "" if file doesn't exist.
func hashFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer func() { _ = f.Close() }()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

// hashString computes the SHA-256 of a string.
func hashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// loadDriftIgnore loads glob patterns from .atv/drift-ignore.
func loadDriftIgnore(root string) []string {
	path := filepath.Join(root, ".atv", "drift-ignore")
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer func() { _ = f.Close() }()

	var patterns []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns
}

// isIgnoredByDrift checks if a path matches any drift-ignore glob pattern.
func isIgnoredByDrift(relPath string, patterns []string) bool {
	slashPath := filepath.ToSlash(relPath)
	for _, pat := range patterns {
		if matched, _ := filepath.Match(pat, slashPath); matched {
			return true
		}
		// Also match against the basename
		if matched, _ := filepath.Match(pat, filepath.Base(slashPath)); matched {
			return true
		}
	}
	return false
}

// sortDriftEntries sorts entries by path for deterministic output.
func sortDriftEntries(entries []DriftEntry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].Path < entries[j-1].Path; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
}
