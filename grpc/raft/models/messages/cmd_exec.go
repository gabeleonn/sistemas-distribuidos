package messages

import (
	pb "raft/autogen"
	"raft/models/node"
)

// CommandExecutionArguments represents the arguments for a CommandExecution RPC call.
type CommandExecutionArguments struct {
	Command node.Command
}

// CommandExecutionReply is the response to a CommandExecution RPC call.
type CommandExecutionReply struct {
	Success  bool
	Value    *string
	LeaderID *int64
}

// ToProto converts CommandExecutionArguments to its protobuf representation.
func (r CommandExecutionArguments) ToProto() *pb.CommandExecutionArguments {
	return &pb.CommandExecutionArguments{
		Command: pb.CommandType(r.Command.Type),
		Key:     r.Command.Key,
		Value:   r.Command.Value, // mantém nil corretamente
	}
}

// CommandExecutionArgumentsFromProto converts a protobuf CommandExecutionArguments to its model representation.
func CommandExecutionArgumentsFromProto(p *pb.CommandExecutionArguments) CommandExecutionArguments {
	var value *string
	if p.Value != nil {
		value = p.Value
	}

	return CommandExecutionArguments{
		Command: node.Command{
			Type:  node.CommandType(p.GetCommand()),
			Key:   p.GetKey(),
			Value: value,
		},
	}
}

// ToProto converts CommandExecutionReply to its protobuf representation.
func (r CommandExecutionReply) ToProto() *pb.CommandExecutionReply {
	return &pb.CommandExecutionReply{
		Success:  r.Success,
		Value:    r.Value,
		LeaderId: r.LeaderID,
	}
}

// CommandExecutionReplyFromProto converts a protobuf CommandExecutionReply to its model representation.
func CommandExecutionReplyFromProto(p *pb.CommandExecutionReply) CommandExecutionReply {
	var value *string
	if p.Value != nil {
		value = p.Value
	}

	var leaderID *int64
	if p.LeaderId != nil {
		leaderID = p.LeaderId
	}

	return CommandExecutionReply{
		Success:  p.GetSuccess(),
		Value:    value,
		LeaderID: leaderID,
	}
}
