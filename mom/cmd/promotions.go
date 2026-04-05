package cmd

import (
	"mom/services/promotions"

	"github.com/spf13/cobra"
)

var promoCmd = &cobra.Command{
	Use:   "promotions",
	Short: "Inicia o servico de promocoes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return promotions.PromotionsHandler()
	},
}

func init() {
	rootCmd.AddCommand(promoCmd)
}
