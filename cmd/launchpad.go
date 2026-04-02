package cmd

import (
	"fmt"
	"os"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/monitor"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/output"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/tui"
	"github.com/spf13/cobra"
)

var staticMode bool

var launchpadCmd = &cobra.Command{
	Use:   "launchpad",
	Short: "Open the ATV launchpad dashboard",
	Long: `Live terminal dashboard showing install intelligence, repo memory,
installed agents/skills, and deterministic next-step recommendations.

Use --static for a one-shot printable view (auto-detected when output
is piped or running inside VS Code Copilot Chat).`,
	RunE: runLaunchpad,
}

func init() {
	launchpadCmd.Flags().BoolVar(&staticMode, "static", false, "Print static dashboard instead of live TUI (auto-detected for non-interactive terminals)")
	rootCmd.AddCommand(launchpadCmd)
}

func runLaunchpad(cmd *cobra.Command, args []string) error {
	targetDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Use static mode when explicitly requested or when not interactive
	if staticMode || !isInteractiveTerminal() {
		snapshot, err := installstate.BuildLaunchpadSnapshot(targetDir)
		if err != nil {
			return fmt.Errorf("failed to build launchpad snapshot: %w", err)
		}
		printer := output.NewPrinter()
		printer.PrintLaunchpad(snapshot)
		return nil
	}

	// Initialize filesystem watcher
	w, err := monitor.NewWatcher(targetDir, monitor.WatcherOptions{})
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	if err := w.Start(cmd.Context()); err != nil {
		return fmt.Errorf("failed to start watcher: %w", err)
	}
	defer w.Stop()

	return tui.RunLaunchpad(targetDir, w)
}

func isInteractiveTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
