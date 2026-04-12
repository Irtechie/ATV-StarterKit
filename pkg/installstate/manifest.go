package installstate

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	ManifestVersion = 1
	manifestDirName = ".atv"
	manifestName    = "install-manifest.json"
)

// ManifestPath returns the canonical repo-local manifest path.
func ManifestPath(root string) string {
	return filepath.Join(root, manifestDirName, manifestName)
}

// WriteManifest atomically writes the current guided installer state to disk.
func WriteManifest(root string, manifest InstallManifest) error {
	if err := ValidateStackPacks(manifest.Requested.StackPacks); err != nil {
		return err
	}

	if manifest.Version == 0 {
		manifest.Version = ManifestVersion
	}
	if manifest.GeneratedAt.IsZero() {
		manifest.GeneratedAt = time.Now().UTC()
	}
	if manifest.RerunPolicy == "" {
		manifest.RerunPolicy = RerunPolicyAdditiveOnly
	}

	path := ManifestPath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create manifest directory: %w", err)
	}

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}
	data = append(data, '\n')

	tempFile, err := os.CreateTemp(filepath.Dir(path), "install-manifest-*.json")
	if err != nil {
		return fmt.Errorf("create temp manifest: %w", err)
	}
	tempPath := tempFile.Name()
	defer func() {
		_ = os.Remove(tempPath)
	}()

	if _, err := tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("write temp manifest: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temp manifest: %w", err)
	}

	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("rename temp manifest: %w", err)
	}

	return nil
}

// ReadManifest loads the repo-local guided installer state from disk.
func ReadManifest(root string) (InstallManifest, error) {
	data, err := os.ReadFile(ManifestPath(root))
	if err != nil {
		return InstallManifest{}, err
	}

	var manifest InstallManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return InstallManifest{}, fmt.Errorf("decode manifest: %w", err)
	}

	return manifest, nil
}

// ComputeFileChecksums walks the given file paths (relative to root) and returns
// a map of path → SHA-256 hex digest for each file that exists.
func ComputeFileChecksums(root string, relPaths []string) map[string]string {
	checksums := make(map[string]string, len(relPaths))
	for _, rel := range relPaths {
		absPath := filepath.Join(root, rel)
		f, err := os.Open(absPath)
		if err != nil {
			continue
		}
		h := sha256.New()
		_, err = io.Copy(h, f)
		_ = f.Close()
		if err != nil {
			continue
		}
		checksums[rel] = hex.EncodeToString(h.Sum(nil))
	}
	return checksums
}
