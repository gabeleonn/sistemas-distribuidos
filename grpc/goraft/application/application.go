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

type Config struct {
	ID    int64
	Addr  string
	Peers []peer.Peer
}

type Application struct {
	config   Config
	logger   *slog.Logger
	server   *grpc.Server
	listener net.Listener
	node     *raft.Node
}

/*
========================================== Core Functionalities ==========================================
*/
func (app *Application) Run(ctx context.Context) error {
	if err := app.openPeerConnections(); err != nil {
		return err
	}
	defer app.closePeerConnections()
	defer app.listener.Close()

	errChannel := make(chan error, 1)

	go func() {
		if err := app.server.Serve(app.listener); err != nil {
			errChannel <- err
		}
	}()

	if err := app.waitForPeers(ctx); err != nil {
		return err
	}

	// SweatSpot
	if app.config.ID == 0 {
		app.RequestVotes(ctx)
	}

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
			app.logger.Info("node stopped gracefully")
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

func (app *Application) waitForPeers(ctx context.Context) error {
	deadline := time.After(10 * time.Second)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	ready := make(map[int64]bool)

	for {
		if len(ready) == len(app.config.Peers) {
			app.logger.Info(fmt.Sprintf("Server ready at %s", app.config.Addr))
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-deadline:
			return fmt.Errorf("timeout waiting for peers")

		case <-ticker.C:
			for i := range app.config.Peers {
				p := &app.config.Peers[i]

				if ready[p.ID] {
					continue
				}

				callCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)

				resp, err := p.Client.Ping(callCtx, &proto.PingRequest{
					Message: "ping",
					FromId:  app.config.ID,
				})

				cancel()

				if err != nil {
					app.logger.Debug(
						"peer not ready yet",
						"peer_id", p.ID,
						"peer_addr", p.Addr,
						"error", err.Error(),
					)
					continue
				}

				app.logger.Debug(
					"peer is ready",
					"peer_id", p.ID,
					"peer_addr", p.Addr,
					"message", resp.Message,
					"from_id", resp.FromId,
				)

				ready[p.ID] = true
			}
		}
	}
}

/*
========================================== Request Votes ==========================================
*/
func (app *Application) RequestVotes(ctx context.Context) {
	for i := range app.config.Peers {
		p := &app.config.Peers[i]

		callContext, cancel := context.WithTimeout(ctx, 1*time.Millisecond)

		resp, err := p.Client.RequestVote(callContext, &proto.RequestVoteRequest{
			Term:        1,
			CandidateId: app.config.ID,
		})

		cancel()

		if err != nil {
			app.logger.Error(
				"Error requesting vote from peer",
				"peer_id", p.ID,
				"peer_addr", p.Addr,
				"error", err.Error(),
			)
			continue
		}

		app.logger.Info(
			"Received vote response from peer",
			"peer_id", p.ID,
			"peer_addr", p.Addr,
			"term", resp.Term,
			"vote_granted", resp.VoteGranted,
		)
	}
}

/*
========================================== Helpers ==========================================
*/
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
		service.NewNodeService(app.node, logger),
	)

	return &app, nil
}
