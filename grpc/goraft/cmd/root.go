package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "goraft",
	Short: "CLI para iniciar os servicos do GoRaft",
	Long:  "Use um tipo de servico para iniciar o processo, como server ou cluster.",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("voce precisa informar o tipo de servico")
	},
}

// Execute executa o comando raiz do CLI.
func Execute() error {
	return rootCommand.Execute()
}
