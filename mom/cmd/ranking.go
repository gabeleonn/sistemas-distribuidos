package cmd

import (
	"mom/core/logger"
	"mom/services/ranking"

	"github.com/spf13/cobra"
)

var rankingCmd = &cobra.Command{
	Use:   "ranking",
	Short: "Inicia o servico de ranking",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ranking.ScoreHandler()
	},
}

func init() {
	logger.Init("ranking")

	rootCmd.AddCommand(rankingCmd)
}
