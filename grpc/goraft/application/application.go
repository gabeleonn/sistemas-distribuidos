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

func (c *Config) majority() int {
	return (len(c.Peers)+1)/2 + 1
}

type Application struct {
	config   Config
	logger   *slog.Logger
	server   *grpc.Server
	listener net.Listener
	node     *raft.Node

	heartbeat chan struct{}
}

/*
==========================================================================================================
========================================== Core Functionalities ==========================================
==========================================================================================================
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

	go app.runRaftLifecycle(ctx)

	select {
	case <-ctx.Done():
		app.shutdown()
		return nil

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

func (app *Application) shutdown() {
	app.logger.Debug("shutting down gRPC server")

	stopped := make(chan struct{})

	go func() {
		app.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		app.logger.Info("node stopped gracefully")

	case <-time.After(5 * time.Second):
		app.logger.Warn("forcing gRPC server shutdown")
		app.server.Stop()
	}
}

/*
===============================================================================================
========================================== Lifecycle ==========================================
===============================================================================================
*/

func (app *Application) runRaftLifecycle(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		role := app.node.Role()

		switch role {
		case raft.Follower:
			app.runFollower(ctx)

		case raft.Candidate:
			app.runCandidate(ctx)

		case raft.Leader:
			app.runLeader(ctx)
		}
	}
}

func (app *Application) runFollower(ctx context.Context) {
	for {
		timeout := raft.RandomElectionTimeout()
		timer := time.NewTimer(timeout)

		select {
		case <-ctx.Done():
			timer.Stop()
			return

		case <-app.heartbeat:
			timer.Stop()
			continue // reset election timer on heartbeat
		case <-timer.C:
			app.node.BecomeCandidate()
			return
		}
	}
}

func (app *Application) runCandidate(ctx context.Context) {
	req := app.node.CandidateRequest()
	votes := app.requestVotes(ctx, req)

	if votes >= app.config.majority() {
		term := app.node.BecomeLeader()

		app.logger.Info(
			"Becoming leader",
			"term", term,
			"votes", votes,
		)

		return
	}

	timeout := raft.RandomElectionTimeout()

	select {
	case <-ctx.Done():
		return

	case <-time.After(timeout):
		app.node.BecomeCandidate()

		return
	}
}

func (app *Application) runLeader(ctx context.Context) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return

	case <-ticker.C:
		req := app.node.AppendEntriesRequest()
		app.sendHeartBeats(ctx, req)
	}
}

/*
===================================================================================================
========================================== Request Votes ==========================================
===================================================================================================
*/

func (app *Application) requestVotes(ctx context.Context, req raft.RequestVoteRequest) int {
	votes := 1

	for i := range app.config.Peers {
		p := &app.config.Peers[i]

		callContext, cancel := context.WithTimeout(ctx, 1*time.Millisecond)

		resp, err := p.Client.RequestVote(callContext, req.ToProto())

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
		if resp.VoteGranted {
			votes++
		}
	}

	return votes
}

func (app *Application) sendHeartBeats(ctx context.Context, req raft.AppendEntriesRequest) {
	state := app.node.GetState()

	for i := range app.config.Peers {
		p := &app.config.Peers[i]

		callContext, cancel := context.WithTimeout(ctx, 500*time.Millisecond)

		resp, err := p.Client.AppendEntries(callContext, req.ToProto())

		cancel()

		if err != nil {
			app.logger.Error(
				"Error sending heartbeat to peer",
				"peer_id", p.ID,
				"peer_addr", p.Addr,
				"error", err.Error(),
			)
			continue
		}

		if !resp.Success {
			if resp.Term > state.CurrentTerm {
				app.node.BecomeFollower(resp.Term)
				return
			}

			continue
		}
	}
}

/*
=============================================================================================
========================================== Helpers ==========================================
=============================================================================================
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

		heartbeat: make(chan struct{}, 1),
	}

	proto.RegisterNodeServer(
		app.server,
		service.NewNodeService(
			app.node,
			logger,
			func() {
				select {
				case app.heartbeat <- struct{}{}:
				default:
				}
			}),
	)

	return &app, nil
}
