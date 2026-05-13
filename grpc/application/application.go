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
	"sync/atomic"
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

	electionReset chan struct{}
	stopped       atomic.Bool
	startCh       chan struct{}
}

/*
==========================================================================================================
========================================== Core Functionalities ==========================================
==========================================================================================================
*/
func (app *Application) Run(ctx context.Context) error {
	defer app.closePeerConnections()
	defer app.listener.Close()

	errChannel := make(chan error, 1)

	go func() {
		app.logger.Info(
			"Starting gRPC server",
			"addr", app.config.Addr,
		)
		if err := app.server.Serve(app.listener); err != nil {
			errChannel <- err
		}
	}()

	go app.runRaftLifecycle(ctx)

	select {
	case <-ctx.Done():
		app.shutdown()
		return nil

	case err := <-errChannel:
		return fmt.Errorf("gRPC server error: %w", err)
	}
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

		if app.stopped.Load() {
			app.logger.Info("node paused, waiting for START")
			select {
			case <-app.startCh:
				app.logger.Info("node resumed")
			case <-ctx.Done():
				return
			}
			continue
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

		case <-app.electionReset:
			timer.Stop()
			if app.stopped.Load() {
				return
			}
			continue
		case <-timer.C:
			app.node.BecomeCandidate()
			return
		}
	}
}

func (app *Application) runCandidate(ctx context.Context) {
	req := app.node.CandidateRequest()
	votes := app.requestVotes(ctx, req)

	if app.node.Role() != raft.Candidate {
		return
	}

	if votes >= app.config.majority() {
		term := app.node.BecomeLeader(app.peerIDs())
		app.sendHeartBeats(ctx)

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

	case <-app.electionReset:
		return

	case <-time.After(timeout):
		app.node.BecomeCandidate()
		return
	}
}

func (app *Application) runLeader(ctx context.Context) {
	if app.stopped.Load() {
		return
	}

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return

	case <-ticker.C:
		app.sendHeartBeats(ctx)
	}
}

/*
===================================================================================================
========================================== Request Votes ==========================================
===================================================================================================
*/

func (app *Application) requestVotes(ctx context.Context, req raft.RequestVoteRequest) int {
	votes := 1
	majority := app.config.majority()

	voteCh := make(chan bool, len(app.config.Peers))
	termCh := make(chan int64, len(app.config.Peers))

	for i := range app.config.Peers {
		p := &app.config.Peers[i]

		go func(peer *peer.Peer) {
			if err := peer.EnsureConnected(); err != nil {
				voteCh <- false
				return
			}

			callContext, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
			defer cancel()

			resp, err := peer.Client.RequestVote(callContext, req.ToProto())
			if err != nil {
				voteCh <- false
				return
			}

			if resp.Term > req.Term {
				termCh <- resp.Term
				return
			}

			voteCh <- resp.VoteGranted
		}(p)
	}

	for responses := 0; responses < len(app.config.Peers); responses++ {
		select {
		case <-ctx.Done():
			return votes

		case term := <-termCh:
			app.node.BecomeFollower(term)
			return votes

		case granted := <-voteCh:
			if granted {
				votes++
			}

			if votes >= majority {
				return votes
			}
		}
	}

	return votes
}

/*
===================================================================================================
========================================== Send Heartbeats ========================================
===================================================================================================
*/
func (app *Application) sendHeartBeats(ctx context.Context) {
	state := app.node.GetState()

	for i := range app.config.Peers {
		p := &app.config.Peers[i]

		if err := p.EnsureConnected(); err != nil {
			continue
		}

		req := app.node.HeartbeatRequestForPeer(p.ID)

		callContext, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		resp, err := p.Client.AppendEntries(callContext, req.ToProto())
		cancel()

		if err != nil {
			continue
		}

		if !resp.Success {
			if resp.Term > state.CurrentTerm {
				app.node.BecomeFollower(resp.Term)
				return
			}

			app.node.HandleAppendFailure(p.ID)
			continue
		}

		if len(req.Entries) > 0 {
			lastEntry := req.Entries[len(req.Entries)-1]
			app.node.HandleAppendSuccess(p.ID, lastEntry.Index)
		}
	}
}

/*
===================================================================================================
======================================== Send Append Entries ======================================
===================================================================================================
*/
func (app *Application) sendAppendEntries(ctx context.Context, entry raft.LogEntry) error {
	prevLogIndex, prevLogTerm, err := app.node.PreviousLogInfo(entry.Index)
	if err != nil {
		return err
	}

	state := app.node.GetState()

	request := raft.NewAppendEntriesRequest(
		state.CurrentTerm,
		state.ID,
		prevLogIndex,
		prevLogTerm,
		[]raft.LogEntry{entry},
		state.CommitIndex,
	)

	majority := app.config.majority()
	successCh := make(chan bool, len(app.config.Peers))
	doneCh := make(chan struct{}, 1)

	for i := range app.config.Peers {
		p := &app.config.Peers[i]

		go func(peer *peer.Peer) {
			if err := peer.EnsureConnected(); err != nil {
				successCh <- false
				return
			}

			callContext, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
			defer cancel()

			resp, err := peer.Client.AppendEntries(callContext, request.ToProto())
			if err != nil {
				successCh <- false
				return
			}

			if resp.Term > state.CurrentTerm {
				app.node.BecomeFollower(resp.Term)
				select {
				case doneCh <- struct{}{}:
				default:
				}
				return
			}

			successCh <- resp.Success
		}(p)
	}

	successCount := 1

	for responses := 0; responses < len(app.config.Peers); responses++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled")

		case <-doneCh:
			return fmt.Errorf("stepping down to follower")

		case success := <-successCh:
			if success {
				successCount++
			}

			if successCount >= majority {
				return nil
			}
		}
	}

	return fmt.Errorf("failed to replicate log entry to majority of peers")
}

/*
=============================================================================================
========================================== Helpers ==========================================
=============================================================================================
*/

func (app *Application) peerIDs() []int64 {
	ids := make([]int64, len(app.config.Peers))
	for i, p := range app.config.Peers {
		ids[i] = p.ID
	}
	return ids
}

func (app *Application) Pause() {
	app.stopped.Store(true)
	select {
	case app.electionReset <- struct{}{}:
	default:
	}
}

func (app *Application) Resume() {
	app.stopped.Store(false)
	select {
	case app.startCh <- struct{}{}:
	default:
	}
}

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

		electionReset: make(chan struct{}, 1),
		startCh:       make(chan struct{}, 1),
	}

	proto.RegisterNodeServer(
		app.server,
		service.NewNodeService(
			app.node,
			logger,
			func() {
				select {
				case app.electionReset <- struct{}{}:
				default:
				}
			},
			func(ctx context.Context, entry raft.LogEntry) error {
				return app.sendAppendEntries(ctx, entry)
			},
			app.Pause,
			app.Resume,
		),
	)

	return &app, nil
}
