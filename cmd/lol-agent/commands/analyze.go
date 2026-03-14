package commands

import (
	"fmt"

	"github.com/huykhuong/lol/internal/prompt"
	"github.com/spf13/cobra"
)

type AnalyzePromptData struct {
	Champion string
	Region   string
}

var (
	champion string
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Perform an AI-powered strategic analysis",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateParams(
			requiredParams{name: "champion", value: champion, flag: "c"},
			requiredParams{name: "region", value: region, flag: "r"},
		); err != nil {
			return err
		}

		opggClient, geminiClient := initializeClients(cmd)

		fmt.Printf("🚀 Analyzing %s (%s) using OP.GG Insights...\n", champion, region)

		data := AnalyzePromptData{
			Champion: champion,
			Region:   region,
		}

		prompt, err := prompt.ReadPrompt("analyze", data)
		if err != nil {
			return fmt.Errorf("failed to read prompt: %w", err)
		}

		result, err := geminiClient.Execute(prompt, opggClient)
		if err != nil {
			return fmt.Errorf("analysis failed: %w", err)
		}

		fmt.Println("\n--- AI COACH ANALYSIS ---")
		fmt.Println(result)
		fmt.Println("--------------------------")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVarP(&champion, "champion", "c", "", "Champion to analyze")
	analyzeCmd.Flags().StringVarP(&region, "region", "r", "NA", "Game region")
}
