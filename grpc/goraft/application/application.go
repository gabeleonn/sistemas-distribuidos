package application

import (
	"context"
	"fmt"
	"goraft/peer"
	"goraft/proto"
	"goraft/raft"
	"goraft/service"
	"goraft/store"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
)

// Config holds the configuration for the Application, including the node's ID, address, and a list of peers in the cluster.
type Config struct {
	ID    int64
	Addr  string
	Peers []peer.Peer
}

// Application represents the main application that manages the gRPC server, peer connections, and logging for a Raft node.
type Application struct {
	config   Config
	logger   *slog.Logger
	server   *grpc.Server
	listener net.Listener
	node     *raft.Node
}

// Run starts the application by opening connections to peers, starting the gRPC server, and handling shutdown gracefully when the context is canceled or an error occurs.
func (app *Application) Run(ctx context.Context) error {
	if err := app.openPeerConnections(); err != nil {
		return err
	}
	defer app.closePeerConnections()
	defer app.listener.Close()

	errChannel := make(chan error, 1)

	go func() {
		app.logger.Info(
			"Starting gRPC server",
			"addr", app.config.Addr,
			"id", app.config.ID,
		)

		if err := app.server.Serve(app.listener); err != nil {
			errChannel <- err
		}
	}()

	// SafeZone
	go app.pingPeers(ctx)

	select {
	case <-ctx.Done():
		app.logger.Debug("shutting down gRPC server")

		stopped := make(chan struct{})

		go func() {
			app.server.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
			app.logger.Debug("gRPC server stopped gracefully")
			return nil

		case <-time.After(5 * time.Second):
			app.logger.Warn("forcing gRPC server shutdown")
			app.server.Stop()
			return nil
		}

	case err := <-errChannel:
		return fmt.Errorf("gRPC server error: %w", err)
	}
}

func (app *Application) openPeerConnections() error {
	for i := range app.config.Peers {
		p := &app.config.Peers[i]
		if err := p.Open(); err != nil {
			return err
		}

		app.logger.Debug(
			"Connected to peer",
			"peer_id", p.ID,
			"peer_addr", p.Addr,
		)
	}
	return nil
}

func (app *Application) closePeerConnections() {
	for i := range app.config.Peers {
		p := &app.config.Peers[i]
		if err := p.Close(); err != nil {
			app.logger.Error(
				"Error closing connection to peer",
				"peer_id", p.ID,
				"peer_addr", p.Addr,
				"error", err,
			)
		}
	}
}

func (app *Application) pingPeers(ctx context.Context) {
	for i := range app.config.Peers {
		p := &app.config.Peers[i]

		if p.Client == nil {
			app.logger.Warn(
				"Skipping ping to peer with no client",
				"peer_id", p.ID,
				"peer_addr", p.Addr,
			)
			continue
		}

		callCtx, cancel := context.WithTimeout(ctx, 3*time.Second)

		resp, err := p.Client.Ping(callCtx, &proto.PingRequest{
			Message: "ping",
			FromId:  app.config.ID,
		})

		cancel()

		if err != nil {
			app.logger.Error(
				"Error pinging peer",
				"peer_id", p.ID,
				"peer_addr", p.Addr,
				"error", err,
			)
			continue
		}

		app.logger.Info(
			"Received ping response",
			"peer_id", p.ID,
			"peer_addr", p.Addr,
			"message", resp.Message,
			"from_id", resp.FromId,
		)
	}
}

// NewApplication creates a new Application instance with the given configuration and logger, sets up the gRPC server, and returns the application ready to run.
func NewApplication(config Config, logger *slog.Logger) (*Application, error) {
	listener, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return nil, err
	}

	storage := store.NewStore()

	app := Application{
		config:   config,
		logger:   logger,
		server:   grpc.NewServer(),
		listener: listener,
		node: raft.NewNode(raft.Config{
			ID:           config.ID,
			StateMachine: storage,
		}),
	}

	proto.RegisterNodeServer(
		app.server,
		service.NewPingService(app.node, logger),
	)

	return &app, nil
}
