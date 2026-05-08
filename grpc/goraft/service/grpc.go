package service

import (
	"context"
	"goraft/proto"
	"goraft/raft"
	"log/slog"
)

type pingService struct {
	proto.UnimplementedNodeServer

	node   *raft.Node
	logger *slog.Logger
}

// NewPingService creates a new instance of the ping service, which implements the NodeServer interface and can be registered with the gRPC server to handle incoming ping requests.
func NewPingService(node *raft.Node, logger *slog.Logger) *pingService {
	return &pingService{
		node:   node,
		logger: logger,
	}
}

func (s *pingService) Ping(ctx context.Context, req *proto.PingRequest) (*proto.PingResponse, error) {
	return &proto.PingResponse{
		Message: "pong",
		FromId:  s.node.ID(),
	}, nil
}
