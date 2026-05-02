package messages

import (
	pb "raft/autogen"
)

// RequestVoteArguments represents the arguments for a RequestVote RPC call.
type RequestVoteArguments struct {
	Term         int64
	CandidateID  int64
	LastLogIndex int64
	LastLogTerm  int64
}

// RequestVoteReply is the response to a RequestVote RPC call.
type RequestVoteReply struct {
	Term        int64
	VoteGranted bool
}

// ToProto converts RequestVoteArguments to its protobuf representation.
func (r RequestVoteArguments) ToProto() *pb.RequestVoteArguments {
	return &pb.RequestVoteArguments{
		Term:         r.Term,
		CandidateId:  r.CandidateID,
		LastLogIndex: r.LastLogIndex,
		LastLogTerm:  r.LastLogTerm,
	}
}

// RequestVoteArgumentsFromProto converts a protobuf RequestVoteArguments to its model representation.
func RequestVoteArgumentsFromProto(p *pb.RequestVoteArguments) RequestVoteArguments {
	return RequestVoteArguments{
		Term:         p.GetTerm(),
		CandidateID:  p.GetCandidateId(),
		LastLogIndex: p.GetLastLogIndex(),
		LastLogTerm:  p.GetLastLogTerm(),
	}
}

// ToProto converts RequestVoteReply to its protobuf representation.
func (r RequestVoteReply) ToProto() *pb.RequestVoteReply {
	return &pb.RequestVoteReply{
		Term:        r.Term,
		VoteGranted: r.VoteGranted,
	}
}

// RequestVoteReplyFromProto converts a protobuf RequestVoteReply to its model representation.
func RequestVoteReplyFromProto(p *pb.RequestVoteReply) RequestVoteReply {
	return RequestVoteReply{
		Term:        p.GetTerm(),
		VoteGranted: p.GetVoteGranted(),
	}
}
