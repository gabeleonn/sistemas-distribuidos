package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Inicia o servico de gateway",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello world")
	},
}

func init() {
	rootCmd.AddCommand(gatewayCmd)
}
