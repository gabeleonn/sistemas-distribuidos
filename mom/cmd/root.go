package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mom",
	Short: "CLI para iniciar os servicos do MOM",
	Long:  "Use um tipo de servico para iniciar o processo, como client ou gateway.",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("voce precisa informar o tipo de servico")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
