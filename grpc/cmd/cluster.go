package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var clusterCommand = &cobra.Command{
	Use:   "cluster",
	Short: "Comando para iniciar um cluster de nós Raft",
	Long:  "Use este comando para iniciar um cluster de nós Raft, especificando o número de nós e outras configurações.",
	RunE: func(cmd *cobra.Command, args []string) error {
		nodesLength, _ := cmd.Flags().GetInt("nodes")
		nodes := buildNodes(nodesLength)

		g, ctx := errgroup.WithContext(context.Background())

		for i := range nodesLength {
			g.Go(func() error {
				cmd := exec.CommandContext(
					ctx,
					"go", "run", "main.go",
					"node",
					"--id", fmt.Sprintf("%d", i),
					"--addr", fmt.Sprintf("localhost:%d", 50050+i),
					"--peers", buildPeers(nodes, i),
				)

				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				return cmd.Run()
			})
		}

		return g.Wait()
	},
}

func init() {
	clusterCommand.Flags().IntP("nodes", "n", 3, "Número de nós no cluster Raft")
	rootCommand.AddCommand(clusterCommand)
}

func buildNodes(nodesLength int) []string {
	nodes := make([]string, nodesLength)
	for i := range nodesLength {
		nodes[i] = fmt.Sprintf("%d=localhost:%d", i, 50050+i)
	}

	return nodes
}

func buildPeers(nodes []string, self int) string {
	if len(nodes) <= 1 {
		return ""
	}
	peers := ""
	for i, node := range nodes {
		if i != self {
			peers += node + ","
		}
	}

	return peers[:len(peers)-1]
}
