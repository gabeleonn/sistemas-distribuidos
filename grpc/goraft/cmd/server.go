package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serverCommand = &cobra.Command{
	Use:    "node",
	Short:  "Comando para iniciar um nó de servidor Raft",
	Long:   "Use este comando para iniciar um nó de servidor Raft, especificando as configurações necessárias.",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("starting server...")
	},
}

func init() {
	rootCommand.AddCommand(serverCommand)
}
