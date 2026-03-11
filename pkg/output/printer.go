package output

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/github/atv-installer/pkg/detect"
	"github.com/github/atv-installer/pkg/scaffold"
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

	bannerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")). // bright yellow
			Bold(true)

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

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("136")) // dim gold for decorative lines
)

// Printer handles terminal output with colored status indicators.
type Printer struct{}

// NewPrinter creates a new Printer.
func NewPrinter() *Printer {
	return &Printer{}
}

// PrintBanner shows the ATV ASCII art logo with gradient yellow effect.
func (p *Printer) PrintBanner() {
	art := strings.TrimRight(bannerText, "\n\r ")
	lines := strings.Split(art, "\n")

	// Top sparkle border
	border := "  ✦ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ ✦"
	fmt.Println(dimStyle.Render(border))
	fmt.Println()

	// Render each line with a gradient color (cycling through the palette)
	for i, line := range lines {
		style := bannerGradient[i%len(bannerGradient)]
		fmt.Println(style.Render("  " + line))
	}

	fmt.Println()
	// Bottom sparkle border
	fmt.Println(dimStyle.Render(border))
	fmt.Println()
	fmt.Println(accentStyle.Render("        ⚡") + titleStyle.Render(" Agentic Tool & Vibes ") + accentStyle.Render("⚡"))
	fmt.Println(dimStyle.Render("        One command. Instant agentic coding."))
	fmt.Println()
}

// PrintDetection shows the detected environment.
func (p *Printer) PrintDetection(env detect.Environment) {
	repoType := "empty directory"
	if env.IsGitRepo {
		repoType = "existing git repo"
	}
	fmt.Printf("  Auto-detected: %s project (%s, %s)\n\n",
		titleStyle.Render(string(env.Stack)), env.StackHint, repoType)
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
		}
	}

	fmt.Println()
	if skipped > 0 {
		fmt.Printf("  %s\n", skipStyle.Render(fmt.Sprintf("⏭️  Skipped %d existing files", skipped)))
	}
	if merged > 0 {
		fmt.Printf("  %s\n", mergeStyle.Render(fmt.Sprintf("🔀 Merged %d JSON configs", merged)))
	}
	fmt.Printf("  %s\n", successStyle.Render(fmt.Sprintf("✅ Created %d files, %d directories", created, dirs)))
}

// PrintNextSteps shows post-install guidance.
func (p *Printer) PrintNextSteps(stack detect.Stack) {
	fmt.Println()
	fmt.Println(successStyle.Render("  🎉 ATV Starter Kit ready!"))
	fmt.Println()
	fmt.Println("  Next steps:")
	fmt.Println(titleStyle.Render("    1.") + " Open this folder in VS Code")
	fmt.Println(titleStyle.Render("    2.") + " Install recommended extensions when prompted")
	fmt.Println(titleStyle.Render("    3.") + ` Try: /ce-brainstorm "your first feature idea"`)
	fmt.Println()
}

// Info prints an informational message.
func (p *Printer) Info(msg string) {
	fmt.Printf("  %s\n", msg)
}
