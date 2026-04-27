package plugingen

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// repoRoot resolves to the top-level ATV-StarterKit checkout. Since
// this test file lives at pkg/plugingen/, the root is two levels up.
func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

// regenerateInto runs Generate against a fresh tempdir mirror of the
// scaffold templates and returns the tempdir root.
func regenerateInto(t *testing.T) string {
	t.Helper()
	src := repoRoot(t)
	tmp := t.TempDir()
	if err := mirrorTree(filepath.Join(src, "pkg", "scaffold", "templates"),
		filepath.Join(tmp, "pkg", "scaffold", "templates")); err != nil {
		t.Fatalf("mirror templates: %v", err)
	}
	if err := Generate(Config{RepoRoot: tmp, KitVersion: "test-1.2.3"}); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	return tmp
}

func TestGenerate_ProducesEveryPerSkillPlugin(t *testing.T) {
	tmp := regenerateInto(t)
	src := filepath.Join(repoRoot(t), "pkg", "scaffold", "templates", "skills")
	skillNames, err := listSkills(src)
	if err != nil {
		t.Fatalf("listSkills: %v", err)
	}
	if len(skillNames) == 0 {
		t.Fatal("no skills discovered — test setup is broken")
	}
	for _, name := range skillNames {
		dir := filepath.Join(tmp, "plugins", "atv-skill-"+name)
		if _, err := os.Stat(filepath.Join(dir, "plugin.json")); err != nil {
			t.Errorf("missing plugin.json for atv-skill-%s: %v", name, err)
		}
		if _, err := os.Stat(filepath.Join(dir, "skills", name, "SKILL.md")); err != nil {
			t.Errorf("missing SKILL.md inside atv-skill-%s: %v", name, err)
		}
	}
}

func TestGenerate_PacksContainExpectedSkills(t *testing.T) {
	tmp := regenerateInto(t)
	for _, p := range Packs() {
		dir := filepath.Join(tmp, "plugins", p.Name)
		for _, sn := range p.SkillNames {
			path := filepath.Join(dir, "skills", sn, "SKILL.md")
			if _, err := os.Stat(path); err != nil {
				t.Errorf("pack %s missing skill %s: %v", p.Name, sn, err)
			}
		}
		// Pack manifest should have category set so users can filter
		// with `copilot plugin marketplace browse`.
		var m PluginManifest
		readJSON(t, filepath.Join(dir, "plugin.json"), &m)
		if m.Category != p.Category {
			t.Errorf("pack %s category mismatch: got %q want %q", p.Name, m.Category, p.Category)
		}
		if len(m.Skills) != 1 || m.Skills[0] != "./skills" {
			t.Errorf("pack %s skills field should be [./skills], got %v", p.Name, m.Skills)
		}
	}
}

func TestGenerate_AtvEverythingBundlesAllSkillsAndAgents(t *testing.T) {
	tmp := regenerateInto(t)
	src := filepath.Join(repoRoot(t), "pkg", "scaffold", "templates")
	skillNames, _ := listSkills(filepath.Join(src, "skills"))
	agentFiles, _ := listAgents(filepath.Join(src, "agents"))

	for _, name := range skillNames {
		path := filepath.Join(tmp, "plugins", "atv-everything", "skills", name, "SKILL.md")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("atv-everything missing skill %s: %v", name, err)
		}
	}
	for _, file := range agentFiles {
		path := filepath.Join(tmp, "plugins", "atv-everything", "agents", file)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("atv-everything missing agent %s: %v", file, err)
		}
	}

	var m PluginManifest
	readJSON(t, filepath.Join(tmp, "plugins", "atv-everything", "plugin.json"), &m)
	if len(m.Skills) == 0 || len(m.Agents) == 0 {
		t.Errorf("atv-everything manifest should declare both skills and agents, got %+v", m)
	}
}

func TestGenerate_AtvAgentsBundlesAllAgents(t *testing.T) {
	tmp := regenerateInto(t)
	src := filepath.Join(repoRoot(t), "pkg", "scaffold", "templates", "agents")
	agentFiles, _ := listAgents(src)

	for _, file := range agentFiles {
		path := filepath.Join(tmp, "plugins", "atv-agents", "agents", file)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("atv-agents missing %s: %v", file, err)
		}
	}
}

func TestGenerate_MarketplaceListsEveryPlugin(t *testing.T) {
	tmp := regenerateInto(t)
	var mp Marketplace
	readJSON(t, filepath.Join(tmp, ".github", "plugin", "marketplace.json"), &mp)

	if mp.Name != "atv-starter-kit" {
		t.Errorf("marketplace name: got %q want atv-starter-kit", mp.Name)
	}
	if mp.Metadata.PluginRoot != "./plugins" {
		t.Errorf("pluginRoot: got %q want ./plugins", mp.Metadata.PluginRoot)
	}

	have := map[string]bool{}
	for _, e := range mp.Plugins {
		have[e.Name] = true
		if e.Source != e.Name {
			t.Errorf("entry %s source should match name, got %q", e.Name, e.Source)
		}
	}

	src := filepath.Join(repoRoot(t), "pkg", "scaffold", "templates", "skills")
	skillNames, _ := listSkills(src)

	wantNames := []string{"atv-everything", "atv-agents"}
	for _, p := range Packs() {
		wantNames = append(wantNames, p.Name)
	}
	for _, n := range skillNames {
		wantNames = append(wantNames, "atv-skill-"+n)
	}
	for _, n := range wantNames {
		if !have[n] {
			t.Errorf("marketplace.json missing entry %q", n)
		}
	}

	// Entries must be sorted alphabetically for deterministic output.
	for i := 1; i < len(mp.Plugins); i++ {
		if mp.Plugins[i-1].Name > mp.Plugins[i].Name {
			t.Errorf("marketplace.json plugins not sorted: %q > %q",
				mp.Plugins[i-1].Name, mp.Plugins[i].Name)
		}
	}
}

func TestGenerate_DeterministicAcrossRuns(t *testing.T) {
	src := repoRoot(t)
	tmpA := t.TempDir()
	tmpB := t.TempDir()
	for _, dst := range []string{tmpA, tmpB} {
		if err := mirrorTree(filepath.Join(src, "pkg", "scaffold", "templates"),
			filepath.Join(dst, "pkg", "scaffold", "templates")); err != nil {
			t.Fatalf("mirror: %v", err)
		}
		if err := Generate(Config{RepoRoot: dst, KitVersion: "det-1.0.0"}); err != nil {
			t.Fatalf("generate: %v", err)
		}
	}
	diffs, err := compareTrees(filepath.Join(tmpA, "plugins"), filepath.Join(tmpB, "plugins"))
	if err != nil {
		t.Fatalf("compareTrees: %v", err)
	}
	if len(diffs) > 0 {
		t.Errorf("non-deterministic generation, diffs: %v", diffs)
	}
}

func TestGenerate_FailsWhenSkillNotInPackOrMisc(t *testing.T) {
	src := repoRoot(t)
	tmp := t.TempDir()
	if err := mirrorTree(filepath.Join(src, "pkg", "scaffold", "templates"),
		filepath.Join(tmp, "pkg", "scaffold", "templates")); err != nil {
		t.Fatalf("mirror: %v", err)
	}
	// Create a fake new skill that's not in any pack or in MiscSkills.
	stray := filepath.Join(tmp, "pkg", "scaffold", "templates", "skills", "zz-stray-skill-for-test")
	if err := os.MkdirAll(stray, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stray, "SKILL.md"),
		[]byte("---\nname: zz-stray\n---\nbody\n"), 0o644); err != nil {
		t.Fatalf("write skill: %v", err)
	}
	err := Generate(Config{RepoRoot: tmp, KitVersion: "stray-test"})
	if err == nil {
		t.Fatal("expected error for unassigned skill, got nil")
	}
	if !strings.Contains(err.Error(), "zz-stray-skill-for-test") {
		t.Errorf("error should mention the stray skill: %v", err)
	}
}

func TestAudit_FlagsNonAllowListedFindings(t *testing.T) {
	contents := map[string]string{
		"skill:test-bad": "this references mcp__plugin_evil__do_thing in a bad way",
	}
	findings := Audit(contents)
	if len(findings) == 0 {
		t.Fatal("expected at least one finding for evil mcp tool")
	}
}

func TestAudit_PassesAllowListedFindings(t *testing.T) {
	// Read the actual ce-plan body which has compound-engineering refs
	// allow-listed. Audit should return nothing for it.
	body, err := os.ReadFile(filepath.Join(repoRoot(t), "pkg", "scaffold", "templates", "skills", "ce-plan", "SKILL.md"))
	if err != nil {
		t.Fatalf("read ce-plan: %v", err)
	}
	contents := map[string]string{"skill:ce-plan": string(body)}
	findings := Audit(contents)
	if len(findings) != 0 {
		t.Errorf("ce-plan should be fully allow-listed, got findings: %v", findings)
	}
}

// readJSON loads a generated JSON file and decodes it into v.
func readJSON(t *testing.T, path string, v interface{}) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(v); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
}
