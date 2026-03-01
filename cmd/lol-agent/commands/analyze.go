package commands

import (
	"fmt"
	"log"

	"github.com/huykhuong/lol/internal/config"
	"github.com/huykhuong/lol/internal/gemini"
	"github.com/huykhuong/lol/internal/opgg"
	"github.com/huykhuong/lol/internal/prompt"
	"github.com/spf13/cobra"
)

type AnalyzePromptData struct {
	Champion string
	Region   string
}

var (
	champion string
	region   string
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Perform an AI-powered strategic analysis",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validate(); err != nil {
			return err
		}

		opggClient, geminiClient := initializeClients(cmd)

		fmt.Printf("🚀 Analyzing %s (%s) using OP.GG Insights...\n", champion, region)

		data := AnalyzePromptData{
			Champion: champion,
			Region: region,
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

func validate() error {
	if champion == "" {
		return fmt.Errorf("champion name is required (use -c)")
	}
	
	if region == "" {
		return fmt.Errorf("region is required (use -r)")
	}

	return nil
}

func initializeClients(cmd *cobra.Command) (*opgg.Client, *gemini.Client) {
	opggClient := opgg.NewClient(config.AppConfig.OPGGURL)
	aiClient, err := gemini.NewClient(cmd.Context(), config.AppConfig.GeminiAPIKey)
	if err != nil {
		log.Fatalf("failed to initialize Gemini Client: %w", err)
	}

	return opggClient, aiClient
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVarP(&champion, "champion", "c", "", "Champion to analyze")
	analyzeCmd.Flags().StringVarP(&region, "region", "r", "NA", "Game region")
}