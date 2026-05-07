package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clusterCommand = &cobra.Command{
	Use:   "cluster",
	Short: "Comando para iniciar um cluster de nós Raft",
	Long:  "Use este comando para iniciar um cluster de nós Raft, especificando o número de nós e outras configurações.",
	Run: func(cmd *cobra.Command, args []string) {
		nodes, _ := cmd.Flags().GetInt("nodes")
		fmt.Printf("starting cluster with %d nodes...\n", nodes)
	},
}

func init() {
	clusterCommand.Flags().IntP("nodes", "n", 3, "Número de nós no cluster Raft")
	rootCommand.AddCommand(clusterCommand)
}
