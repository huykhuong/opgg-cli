package commands

import (
	"fmt"
	"log"

	"github.com/huykhuong/lol/internal/config"
	"github.com/huykhuong/lol/internal/gemini"
	"github.com/huykhuong/lol/internal/opgg"
	"github.com/spf13/cobra"
)

var region string // shared across commands

type requiredParams struct {
	name  string
	value string
	flag  string
}

func validateParams(params ...requiredParams) error {
	for _, param := range params {
		if param.value == "" {
			return fmt.Errorf("the %s parameter is required", param.flag)
		}
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
