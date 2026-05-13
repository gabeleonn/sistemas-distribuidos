package service

import (
	"context"
	"encoding/json"
	"goraft/proto"
	"goraft/raft"
	"log/slog"
)

type NodeService struct {
	proto.UnimplementedNodeServer

	node        *raft.Node
	logger      *slog.Logger
	onHeartbeat func()
	onCommand   func(ctx context.Context, entry raft.LogEntry) error
	onStop      func()
	onStart     func()
	iddle       bool
}

func NewNodeService(
	node *raft.Node,
	logger *slog.Logger,
	onHeartbeat func(),
	onCommand func(ctx context.Context, entry raft.LogEntry) error,
	onStop func(),
	onStart func(),
) *NodeService {
	return &NodeService{
		node:        node,
		logger:      logger,
		onHeartbeat: onHeartbeat,
		onCommand:   onCommand,
		onStop:      onStop,
		onStart:     onStart,
		iddle:       false,
	}
}

func (s *NodeService) RequestVote(
	ctx context.Context,
	req *proto.RequestVoteRequest,
) (*proto.RequestVoteResponse, error) {
	if s.iddle {
		return &proto.RequestVoteResponse{
			VoteGranted: false,
		}, nil
	}

	request := raft.RequestVoteRequestFromProto(req)
	vote := s.node.CandidateResponse(*request)

	if vote.VoteGranted {
		s.onHeartbeat()
	}

	return vote.ToProto(), nil
}

func (s *NodeService) AppendEntries(
	ctx context.Context,
	req *proto.AppendEntriesRequest,
) (*proto.AppendEntriesResponse, error) {
	if s.iddle {
		return &proto.AppendEntriesResponse{
			Success: false,
		}, nil
	}

	request := raft.AppendEntriesRequestFromProto(req)
	response := s.node.HeartbeatResponse(*request)

	if response.Success {
		s.onHeartbeat()
	}

	return response.ToProto(), nil
}

func (s *NodeService) ExecuteCommand(
	ctx context.Context,
	req *proto.CommandRequest,
) (*proto.CommandResponse, error) {
	command, err := raft.CommandFromProto(req.Command)
	if err != nil {
		s.logger.Info("Received command", "command", req.Command, "error", err.Error())

		return &proto.CommandResponse{
			Success: false,
			Message: "invalid command format",
		}, nil
	}

	if command.Type == "NODE" {
		switch command.Key {
		case "STATE":
			state := s.node.GetState()
			message := state.String()

			s.logger.Info("Received command", "command", req.Command, "success", true)

			return &proto.CommandResponse{
				Success: true,
				Message: message,
			}, nil
		case "STORE":
			data := s.node.GetStoreState()
			bytes, _ := json.MarshalIndent(data, "", " ")
			message := string(bytes)

			s.logger.Info("Received command", "command", req.Command, "success", true)

			return &proto.CommandResponse{
				Success: true,
				Message: message,
			}, nil

		case "STOP":
			s.logger.Info("Received command", "command", req.Command, "success", true)
			s.iddle = true
			s.onStop()

			return &proto.CommandResponse{
				Success: true,
				Message: "stopping node",
			}, nil
		case "START":
			s.logger.Info("Received command", "command", req.Command, "success", true)
			s.iddle = false
			s.onStart()

			return &proto.CommandResponse{
				Success: true,
				Message: "starting node",
			}, nil

		default:
			s.logger.Info("Received command", "command", req.Command, "error", "unknown node command")

			return &proto.CommandResponse{
				Success: false,
				Message: "unknown node command",
			}, nil
		}
	}

	if _, err := s.node.EnsureLeader(); err != nil {
		return &proto.CommandResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	if command.Type == "GET" {
		response := s.node.ExecuteGetCommandResponse(command)

		s.logger.Info("Received command", "command", req.Command, "success", response.Success)

		return response.ToProto(), nil
	}

	logentry, err := s.node.AppendLog(command)
	if err != nil {
		s.logger.Info("Received command", "command", req.Command, "error", err.Error())

		return &proto.CommandResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	if err := s.onCommand(ctx, logentry); err != nil {
		s.logger.Info("Received command", "command", req.Command, "error", err.Error())

		return &proto.CommandResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	response := s.node.ExecuteCommandResponse(logentry)
	s.logger.Info("Received command", "command", req.Command, "success", response.Success)

	return response.ToProto(), nil
}
