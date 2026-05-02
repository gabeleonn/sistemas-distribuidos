package messages

import (
	pb "raft/autogen"
)

// StreamNodeStateArguments represents the arguments for a StreamNodeState RPC call.
type StreamNodeStateArguments struct {
	ID int64
}

// ToProto converts StreamNodeStateArguments to its protobuf representation.
func (r StreamNodeStateArguments) ToProto() *pb.StreamNodeStateArguments {
	return &pb.StreamNodeStateArguments{
		Id: r.ID,
	}
}

// StreamNodeStateArgumentsFromProto converts a protobuf StreamNodeStateArguments to its model representation.
func StreamNodeStateArgumentsFromProto(p *pb.StreamNodeStateArguments) StreamNodeStateArguments {
	return StreamNodeStateArguments{
		ID: p.GetId(),
	}
}

// NodeInfoArguments represents the arguments for AddNode and RemoveNode RPC calls.
type NodeInfoArguments struct {
	ID   int64
	Addr *string // optional
}

// ToProto converts NodeInfoArguments to its protobuf representation.
func (r NodeInfoArguments) ToProto() *pb.NodeInfoArguments {
	return &pb.NodeInfoArguments{
		Id:   r.ID,
		Addr: r.Addr,
	}
}

// NodeInfoArgumentsFromProto converts a protobuf NodeInfoArguments to its model representation.
func NodeInfoArgumentsFromProto(p *pb.NodeInfoArguments) NodeInfoArguments {
	var addr *string
	if p.Addr != nil {
		v := p.GetAddr()
		addr = &v
	}

	return NodeInfoArguments{
		ID:   p.GetId(),
		Addr: addr,
	}
}
