package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/concierge"
	"github.com/spf13/cobra"
)

var conciergeCmd = &cobra.Command{
	Use:   "concierge",
	Short: "Typed tools for the optional Copilot SDK assistant",
	Long: `Expose the local memory index and deterministic recommendations as
structured JSON for a GitHub Copilot SDK assistant.

Each subcommand returns a JSON object to stdout. The assistant should
present results to the user, never silently override deterministic ranking.`,
}

var memorySummaryCmd = &cobra.Command{
	Use:   "memory-summary",
	Short: "Return a structured overview of repo memory and install intelligence",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := os.Getwd()
		if err != nil {
			return err
		}
		return printJSON(concierge.GetMemorySummary(root))
	},
}

var listRecommendationsCmd = &cobra.Command{
	Use:   "list-recommendations",
	Short: "Return the deterministic next-step recommendations",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := os.Getwd()
		if err != nil {
			return err
		}
		return printJSON(concierge.ListRecommendations(root))
	},
}

var explainRecommendationCmd = &cobra.Command{
	Use:   "explain-recommendation [id]",
	Short: "Return a detailed explanation for a single recommendation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := os.Getwd()
		if err != nil {
			return err
		}
		return printJSON(concierge.ExplainRecommendation(root, args[0]))
	},
}

var openArtifactCmd = &cobra.Command{
	Use:   "open-artifact [name]",
	Short: "Resolve a logical artifact name to a filesystem path",
	Long:  "Known artifacts: manifest, instructions, brainstorms, plans, solutions, agents, skills",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := os.Getwd()
		if err != nil {
			return err
		}
		return printJSON(concierge.OpenArtifact(root, args[0]))
	},
}

var runSuggestedActionCmd = &cobra.Command{
	Use:   "run-suggested-action [action-id]",
	Short: "Describe a suggested action for user confirmation (never executes)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("action-id is required")
		}
		root, err := os.Getwd()
		if err != nil {
			return err
		}
		return printJSON(concierge.RunSuggestedAction(root, args[0]))
	},
}

func init() {
	conciergeCmd.AddCommand(memorySummaryCmd)
	conciergeCmd.AddCommand(listRecommendationsCmd)
	conciergeCmd.AddCommand(explainRecommendationCmd)
	conciergeCmd.AddCommand(openArtifactCmd)
	conciergeCmd.AddCommand(runSuggestedActionCmd)
	rootCmd.AddCommand(conciergeCmd)
}

func printJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
