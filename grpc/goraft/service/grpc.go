package service

import (
	"context"
	"goraft/proto"
	"goraft/raft"
	"log/slog"
)

// NodeService implements the gRPC service for handling requests related to a Raft node,
// such as pinging the node to check if it's alive and responding with its ID.
type NodeService struct {
	proto.UnimplementedNodeServer

	node   *raft.Node
	logger *slog.Logger
}

// NewNodeService creates a new instance of the NodeService, which implements the NodeServer
// interface and can be registered with the gRPC server to handle incoming requests.
func NewNodeService(node *raft.Node, logger *slog.Logger) *NodeService {
	return &NodeService{
		node:   node,
		logger: logger,
	}
}

// Ping responds to a ping request with a pong message and the ID of the node,
// allowing clients to check if the node is alive and get its ID.
func (s *NodeService) Ping(ctx context.Context, req *proto.PingRequest) (*proto.PingResponse, error) {
	return &proto.PingResponse{
		Message: "pong",
		FromId:  s.node.ID(),
	}, nil
}

// RequestVote handles incoming RequestVote RPCs from other nodes in the cluster,
// allowing the node to participate in leader elections by granting or denying votes
// based on its current state and the information provided in the request.
func (s *NodeService) RequestVote(
	ctx context.Context,
	req *proto.RequestVoteRequest,
) (*proto.RequestVoteResponse, error) {
	request := raft.RequestVoteRequestFromProto(req)
	vote := s.node.RequestVote(*request)

	return vote.ToProto(), nil
}
