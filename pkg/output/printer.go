package output

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/detect"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/gstack"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/scaffold"
	"github.com/charmbracelet/lipgloss"
)

//go:embed banner.txt
var bannerText string

var (
	// Gradient yellow palette for banner lines (top to bottom: gold → bright yellow → white-yellow)
	bannerGradient = []lipgloss.Style{
		lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true), // deep gold
		lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true), // bright yellow
		lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true), // yellow
		lipgloss.NewStyle().Foreground(lipgloss.Color("228")).Bold(true), // light yellow
		lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Bold(true), // pale yellow
		lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Bold(true), // near-white yellow
	}

	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")). // orange accent
			Bold(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("178")). // golden yellow
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")). // yellow-green
			Bold(true)

	skipStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")) // dim gray

	mergeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")) // orange-yellow

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	failureStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("136")) // dim gold for decorative lines

	cloneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")). // blue
			Bold(true)

	buildStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")). // orange
			Bold(true)
)

// Printer handles terminal output with colored status indicators.
type Printer struct{}

// NewPrinter creates a new Printer.
func NewPrinter() *Printer {
	return &Printer{}
}

// PrintBanner shows the retro terminal-style ATV 2.0 banner.
func (p *Printer) PrintBanner() {
	art := strings.TrimRight(bannerText, "\n\r ")
	lines := strings.Split(art, "\n")

	border := "  ✦ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ ✦"

	// Boot messages (retro terminal style)
	fmt.Println(dimStyle.Render(border))
	fmt.Println()
	fmt.Println(dimStyle.Render("  > booting all the vibes 2.0 starter kit..."))
	fmt.Println(dimStyle.Render("  > compound-engineering + gstack + memory + agent browser"))
	fmt.Println()

	// Render ASCII art with gradient
	for i, line := range lines {
		style := bannerGradient[i%len(bannerGradient)]
		fmt.Println(style.Render("  " + line))
	}

	fmt.Println()
	fmt.Println(dimStyle.Render(border))
	fmt.Println()
	fmt.Println(accentStyle.Render("  ⚡") + titleStyle.Render(" All The Vibes 2.0 ") + accentStyle.Render("⚡"))
	fmt.Println(dimStyle.Render("  One command. Full agentic coding setup."))
	fmt.Println()
}

// PrintDetection shows the detected environment.
func (p *Printer) PrintDetection(env detect.Environment) {
	repoType := "empty directory"
	if env.IsGitRepo {
		repoType = "existing git repo"
	}
	fmt.Printf("  Auto-detected primary: %s project (%s, %s)\n",
		titleStyle.Render(string(env.Stack)), env.StackHint, repoType)
	if len(env.DetectedPacks) > 0 {
		labels := make([]string, 0, len(env.DetectedPacks))
		for _, pack := range env.DetectedPacks {
			labels = append(labels, stackPackTitle(pack))
		}
		fmt.Printf("  Likely stack packs: %s\n", titleStyle.Render(strings.Join(labels, ", ")))
	}
	fmt.Println()
}

func stackPackTitle(pack installstate.StackPack) string {
	switch pack {
	case installstate.StackPackRails:
		return "Rails"
	case installstate.StackPackPython:
		return "Python"
	case installstate.StackPackTypeScript:
		return "TypeScript"
	default:
		return "General"
	}
}

// PrintResults shows what was created, skipped, or merged.
func (p *Printer) PrintResults(results []scaffold.WriteResult) {
	created := 0
	skipped := 0
	merged := 0
	dirs := 0

	for _, r := range results {
		switch r.Status {
		case scaffold.StatusCreated:
			fmt.Printf("  %s %s\n", successStyle.Render("✅"), r.Path)
			created++
		case scaffold.StatusSkipped:
			fmt.Printf("  %s %s\n", skipStyle.Render("⏭️  "+r.Path+" (exists)"), "")
			skipped++
		case scaffold.StatusMerged:
			fmt.Printf("  %s %s\n", mergeStyle.Render("🔀 "+r.Path+" (merged)"), "")
			merged++
		case scaffold.StatusDirCreated:
			fmt.Printf("  %s %s\n", successStyle.Render("📁"), r.Path)
			dirs++
		case scaffold.StatusFailed:
			fmt.Printf("  %s %s — %s\n", failureStyle.Render("❌"), r.Path, r.Error)
		}
	}

	fmt.Println()
	if skipped > 0 {
		fmt.Printf("  %s\n", skipStyle.Render(fmt.Sprintf("⏭️  Skipped %d existing files", skipped)))
	}
	if merged > 0 {
		fmt.Printf("  %s\n", mergeStyle.Render(fmt.Sprintf("🔀 Merged %d JSON configs", merged)))
	}
	failed := scaffold.SummarizeResults(results).Failed
	if failed > 0 {
		fmt.Printf("  %s\n", failureStyle.Render(fmt.Sprintf("❌ %d writes failed", failed)))
	}
	fmt.Printf("  %s\n", successStyle.Render(fmt.Sprintf("✅ Created %d files, %d directories", created, dirs)))
}

// PrintGuidedSummary shows the structured result of a guided install.
func (p *Printer) PrintGuidedSummary(outcomes []installstate.InstallOutcome, manifestPath string) {
	fmt.Print(guidedSummaryText(outcomes, manifestPath))
}

// PrintRecommendations shows the deterministic next-step recommendations derived from local state.
func (p *Printer) PrintRecommendations(recommendations []installstate.Recommendation) {
	if len(recommendations) == 0 {
		return
	}

	fmt.Println(titleStyle.Render("  Recommended next moves"))
	fmt.Println()
	for i, recommendation := range recommendations {
		fmt.Println(titleStyle.Render(fmt.Sprintf("    %d.", i+1)) + " " + recommendation.Title)
		fmt.Println(dimStyle.Render("       " + recommendation.Reason))
	}
	fmt.Println()
}

// PrintLaunchpad renders the reopenable terminal launchpad for the current repository.
func (p *Printer) PrintLaunchpad(snapshot installstate.LaunchpadSnapshot) {
	fmt.Print(launchpadText(snapshot))
}

// PrintNextSteps shows post-install guidance.
func (p *Printer) PrintNextSteps(hasGstack bool, hasAgentBrowser bool, manifestPath string) {
	fmt.Println()
	fmt.Println(successStyle.Render("  🎉 ATV Starter Kit ready!"))
	fmt.Println()
	fmt.Println("  Next steps:")
	fmt.Println(titleStyle.Render("    1.") + " Open this folder in VS Code")
	fmt.Println(titleStyle.Render("    2.") + " Install recommended extensions when prompted")
	fmt.Println(titleStyle.Render("    3.") + ` Try: /ce-brainstorm "your first feature idea"`)
	step := 4
	if hasGstack {
		fmt.Println(titleStyle.Render(fmt.Sprintf("    %d.", step)) + ` Try: /gstack-office-hours to start a gstack sprint`)
		step++
	}
	if hasAgentBrowser {
		fmt.Println(titleStyle.Render(fmt.Sprintf("    %d.", step)) + ` Try: agent-browser open https://yourapp.com`)
	}
	fmt.Println()
	if hasGstack {
		fmt.Println(dimStyle.Render("  Note: gstack creates ~/.gstack/ for session tracking and config."))
	}
	if manifestPath != "" {
		fmt.Println(dimStyle.Render("  Install state saved to " + manifestPath + " for future reopen/launchpad work."))
		fmt.Println(dimStyle.Render("  Reopen later with: atv-installer launchpad"))
	}
	fmt.Println()
}

func launchpadText(snapshot installstate.LaunchpadSnapshot) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(titleStyle.Render("  ATV Launchpad"))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  Local memory + install intelligence for this repo"))
	b.WriteString("\n\n")

	b.WriteString(titleStyle.Render("  Installed intelligence"))
	b.WriteString("\n")
	if snapshot.HasManifest {
		b.WriteString(fmt.Sprintf("  %s Guided manifest found at %s\n", successStyle.Render("✅"), snapshot.ManifestPath))
		if !snapshot.GeneratedAt.IsZero() {
			b.WriteString(fmt.Sprintf("  %s Last guided run: %s\n", dimStyle.Render("•"), snapshot.GeneratedAt.Format("2006-01-02 15:04 MST")))
		}
		if snapshot.Requested.PresetName != "" {
			b.WriteString(fmt.Sprintf("  %s Preset: %s\n", dimStyle.Render("•"), snapshot.Requested.PresetName))
		}
		if labels := snapshot.StackPackLabels(); len(labels) > 0 {
			b.WriteString(fmt.Sprintf("  %s Stack packs: %s\n", dimStyle.Render("•"), strings.Join(labels, ", ")))
		}
		if snapshot.HasGstack() {
			mode := "markdown-only"
			if snapshot.Requested.GstackRuntime {
				mode = "runtime requested"
			}
			b.WriteString(fmt.Sprintf("  %s gstack: %d skill dirs requested (%s)\n", dimStyle.Render("•"), len(snapshot.Requested.GstackDirs), mode))
		} else {
			b.WriteString(fmt.Sprintf("  %s gstack: not requested in the last guided run\n", dimStyle.Render("•")))
		}
		if snapshot.HasAgentBrowser() {
			b.WriteString(fmt.Sprintf("  %s agent-browser: requested in the last guided run\n", dimStyle.Render("•")))
		} else {
			b.WriteString(fmt.Sprintf("  %s agent-browser: not requested in the last guided run\n", dimStyle.Render("•")))
		}
		b.WriteString(fmt.Sprintf("  %s Outcomes: %d done, %d warnings, %d failed, %d skipped\n", dimStyle.Render("•"), snapshot.OutcomeSummary.Done, snapshot.OutcomeSummary.Warning, snapshot.OutcomeSummary.Failed, snapshot.OutcomeSummary.Skipped))
	} else {
		b.WriteString(fmt.Sprintf("  %s No guided manifest found yet. Run %q to capture install state.\n", warningStyle.Render("⚠️"), "atv-installer init --guided"))
	}

	b.WriteString("\n")
	b.WriteString(titleStyle.Render("  Repo memory snapshot"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s brainstorms: %d\n", dimStyle.Render("•"), snapshot.RepoState.BrainstormCount))
	b.WriteString(fmt.Sprintf("  %s plans: %d\n", dimStyle.Render("•"), snapshot.RepoState.PlanCount))
	b.WriteString(fmt.Sprintf("  %s solutions: %d\n", dimStyle.Render("•"), snapshot.RepoState.SolutionCount))
	if snapshot.RepoState.HasUncheckedPlan {
		b.WriteString(fmt.Sprintf("  %s active plan state: unchecked work remains\n", dimStyle.Render("•")))
	} else if snapshot.RepoState.HasCompletedPlan {
		b.WriteString(fmt.Sprintf("  %s active plan state: at least one completed plan found\n", dimStyle.Render("•")))
	}
	if snapshot.RepoState.InstalledAgents > 0 || snapshot.RepoState.InstalledSkills > 0 || snapshot.RepoState.HasCopilotInstructions {
		b.WriteString(fmt.Sprintf("  %s agents: %d, skills: %d\n", dimStyle.Render("•"), snapshot.RepoState.InstalledAgents, snapshot.RepoState.InstalledSkills))
		if snapshot.RepoState.HasCopilotInstructions {
			b.WriteString(fmt.Sprintf("  %s copilot-instructions.md: present\n", dimStyle.Render("•")))
		}
	}
	if snapshot.RepoState.HasGstackStaging {
		b.WriteString(fmt.Sprintf("  %s gstack staging: present (.gstack/)\n", dimStyle.Render("•")))
	}
	if snapshot.RepoState.HasAgentBrowserSkill {
		b.WriteString(fmt.Sprintf("  %s agent-browser skill: installed\n", dimStyle.Render("•")))
	}

	b.WriteString("\n")
	b.WriteString(titleStyle.Render("  Recommended next moves"))
	b.WriteString("\n")
	if len(snapshot.Recommendations) == 0 {
		b.WriteString(fmt.Sprintf("  %s No recommendations yet. Add repo memory or rerun guided install.\n", dimStyle.Render("•")))
	} else {
		for i, recommendation := range snapshot.CloneRecommendations() {
			b.WriteString(fmt.Sprintf("  %s %s\n", titleStyle.Render(fmt.Sprintf("%d.", i+1)), recommendation.Title))
			b.WriteString(fmt.Sprintf("     %s\n", dimStyle.Render(recommendation.Reason)))
		}
	}

	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  Reopen this dashboard any time with: atv-installer launchpad"))
	if snapshot.HasManifest {
		b.WriteString("\n")
		b.WriteString(dimStyle.Render("  Manifest path: "))
		b.WriteString(filepath.ToSlash(snapshot.ManifestPath))
	}
	b.WriteString("\n\n")
	return b.String()
}

func guidedSummaryText(outcomes []installstate.InstallOutcome, manifestPath string) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(titleStyle.Render("  Guided install summary"))
	b.WriteString("\n\n")
	for _, outcome := range outcomes {
		b.WriteString(fmt.Sprintf("  %s %s", guidedOutcomeIcon(outcome.Status), outcome.Step))
		detailParts := make([]string, 0, 2)
		if outcome.Detail != "" {
			detailParts = append(detailParts, outcome.Detail)
		}
		if outcome.Duration != "" {
			detailParts = append(detailParts, outcome.Duration)
		}
		if len(detailParts) > 0 {
			b.WriteString(" (")
			b.WriteString(strings.Join(detailParts, " · "))
			b.WriteString(")")
		}
		if outcome.Reason != "" {
			b.WriteString(" — ")
			b.WriteString(outcome.Reason)
		}
		b.WriteString("\n")
	}
	if manifestPath != "" {
		b.WriteString("\n")
		b.WriteString(dimStyle.Render("  Install manifest: "))
		b.WriteString(manifestPath)
		b.WriteString("\n")
	}
	b.WriteString("\n")
	return b.String()
}

func guidedOutcomeIcon(status installstate.InstallStepStatus) string {
	switch status {
	case installstate.InstallStepWarning:
		return warningStyle.Render("⚠️ ")
	case installstate.InstallStepFailed:
		return failureStyle.Render("❌")
	case installstate.InstallStepSkipped:
		return skipStyle.Render("⏭️ ")
	default:
		return successStyle.Render("✅")
	}
}

// PrintGstackStart shows the beginning of gstack installation.
func (p *Printer) PrintGstackStart(mode gstack.InstallMode) {
	fmt.Println()
	modeLabel := "markdown-only"
	if mode == gstack.ModeFullRuntime {
		modeLabel = "full runtime"
	}
	fmt.Printf("  %s Installing gstack (%s)...\n", cloneStyle.Render("🔗"), modeLabel)
}

// PrintGstackResult shows the result of gstack installation.
func (p *Printer) PrintGstackResult(result *gstack.InstallResult) {
	if result.Warning != "" {
		fmt.Printf("  %s %s\n", skipStyle.Render("⚠️"), result.Warning)
		return
	}
	if result.Cloned {
		fmt.Printf("  %s gstack cloned (%d skills)\n", cloneStyle.Render("🔗"), len(result.SkillDirs))
	}
	if result.Built {
		fmt.Printf("  %s gstack binary built\n", buildStyle.Render("🔨"))
	}
}

// GstackError shows a gstack installation error.
func (p *Printer) GstackError(err error) {
	fmt.Printf("  %s gstack install failed: %v\n", skipStyle.Render("⚠️"), err)
	fmt.Println(dimStyle.Render("    ATV files were installed successfully. gstack can be added later."))
}

// Info prints an informational message.
func (p *Printer) Info(msg string) {
	fmt.Printf("  %s\n", msg)
}
