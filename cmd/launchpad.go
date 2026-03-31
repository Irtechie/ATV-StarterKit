package cmd

import (
	"fmt"
	"os"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/output"
	"github.com/spf13/cobra"
)

var launchpadCmd = &cobra.Command{
	Use:   "launchpad",
	Short: "Open the local ATV launchpad for the current repo",
	Long:  "Read the local install manifest and repo memory state, then render deterministic next moves for the current repository.",
	RunE:  runLaunchpad,
}

func init() {
	rootCmd.AddCommand(launchpadCmd)
}

func runLaunchpad(cmd *cobra.Command, args []string) error {
	targetDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	snapshot, err := installstate.BuildLaunchpadSnapshot(targetDir)
	if err != nil {
		return fmt.Errorf("failed to build launchpad snapshot: %w", err)
	}

	printer := output.NewPrinter()
	printer.PrintLaunchpad(snapshot)
	return nil
}
