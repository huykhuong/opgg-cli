package commands

import (
	"fmt"
	"os"

	"github.com/huykhuong/lol/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lol-agent",
	Short: "AI League Analyst",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.Load()
	},
}

func Execute() {
	if len(os.Args) == 1 {
		printHelp()
		os.Exit(0)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  lolcli analyze -c <champion> -r <region>")
	fmt.Println("  lolcli help")
}