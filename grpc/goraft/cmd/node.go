package cmd

import (
	"context"
	"fmt"
	"goraft/application"
	"goraft/peer"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var serverCommand = &cobra.Command{
	Use:    "node",
	Short:  "Comando para iniciar um nó de servidor Raft",
	Long:   "Use este comando para iniciar um nó de servidor Raft, especificando as configurações necessárias.",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		nodeID, err := cmd.Flags().GetInt("id")
		if err != nil {
			return err
		}

		nodeAddr, err := cmd.Flags().GetString("addr")
		if err != nil {
			return err
		}

		nodePeers, err := cmd.Flags().GetStringSlice("peers")
		if err != nil {
			return err
		}

		peers, err := parsePeers(nodePeers)
		if err != nil {
			return err
		}

		level := slog.LevelInfo
		if os.Getenv("DEBUG") == "*" || os.Getenv("DEBUG") == "goraft" {
			level = slog.LevelDebug
		}

		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})

		logger := slog.New(handler).With("node", nodeID)

		app, err := application.NewApplication(application.Config{
			ID:    int64(nodeID),
			Addr:  nodeAddr,
			Peers: peers,
		}, logger)

		if err != nil {
			return err
		}

		return app.Run(ctx)
	},
}

func init() {
	serverCommand.Flags().IntP("id", "i", 0, "ID do nó Raft")
	serverCommand.Flags().StringP("addr", "a", "localhost:50050", "Endereço do nó Raft")
	serverCommand.Flags().
		StringSliceP("peers", "r", []string{}, "Peers do nó Raft no formato id:addr (ex: 1=localhost:50051,2=localhost:50052)")
	rootCommand.AddCommand(serverCommand)
}

func parsePeers(peers []string) ([]peer.Peer, error) {
	peerList := make([]peer.Peer, len(peers))
	for i, p := range peers {
		var id int64
		var addr string
		_, err := fmt.Sscanf(p, "%d=%s", &id, &addr)
		if err != nil {
			return nil, fmt.Errorf("invalid peer format: %s", p)
		}
		peerList[i] = peer.Peer{ID: id, Addr: addr}
	}
	return peerList, nil
}
