package commands

import (
	"fmt"

	"github.com/huykhuong/lol/internal/prompt"
	"github.com/spf13/cobra"
)

type ProfilePromptData struct {
	PlayerName string
	PlayerTag  string
	Region     string
	Question   string
}

var (
	playerName string
	playerTag  string
	question   string
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Ask questions about a player's profile/recent matches",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateParams(
			requiredParams{name: "player-name", value: playerName, flag: "p"},
			requiredParams{name: "player-tag", value: playerTag, flag: "t"},
			requiredParams{name: "region", value: region, flag: "r"},
			requiredParams{name: "question", value: question, flag: "q"},
		); err != nil {
			return err
		}

		opggClient, geminiClient := initializeClients(cmd)

		fmt.Printf("🚀 Analyzing %s (%s) using OP.GG Insights...\n", playerName, playerTag)

		data := ProfilePromptData{
			PlayerName: playerName,
			PlayerTag:  playerTag,
			Region:     region,
			Question:   question,
		}

		prompt, err := prompt.ReadPrompt("profile", data)
		if err != nil {
			return fmt.Errorf("failed to read prompt: %w", err)
		}

		result, err := geminiClient.Execute(prompt, opggClient)
		if err != nil {
			return fmt.Errorf("profile failed: %w", err)
		}

		fmt.Printf("\n--- Profile of Summoner %s (%s) in %s ---\n", playerName, playerTag, region)
		fmt.Println(result)
		fmt.Println("--------------------------")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.Flags().StringVarP(&playerName, "player-name", "p", "", "Player name")
	profileCmd.Flags().StringVarP(&playerTag, "player-tag", "t", "", "Player tag")
	profileCmd.Flags().StringVarP(&region, "region", "r", "NA", "Game region")
	profileCmd.Flags().StringVarP(&question, "question", "q", "", "Question to answer about the player's profile/recent matches")
}
