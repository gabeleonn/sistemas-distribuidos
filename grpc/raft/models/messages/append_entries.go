package messages

import (
	pb "raft/autogen"
	"raft/models/node"
)

// AppendEntriesArguments represents the arguments for an AppendEntries RPC call.
type AppendEntriesArguments struct {
	Term         int64
	LeaderID     int64
	PrevLogIndex int64
	PrevLogTerm  int64
	Entries      []node.LogEntry
	LeaderCommit int64
}

// AppendEntriesReply is the response to an AppendEntries RPC call.
type AppendEntriesReply struct {
	Term    int64
	Success bool
}

// ToProto converts AppendEntriesArguments to its protobuf representation.
func (r AppendEntriesArguments) ToProto() *pb.AppendEntriesArguments {
	entries := make([]*pb.LogEntry, len(r.Entries))
	for i, entry := range r.Entries {
		entries[i] = entry.ToProto()
	}

	return &pb.AppendEntriesArguments{
		Term:         r.Term,
		LeaderId:     r.LeaderID,
		PrevLogIndex: r.PrevLogIndex,
		PrevLogTerm:  r.PrevLogTerm,
		Entries:      entries,
		LeaderCommit: r.LeaderCommit,
	}
}

// AppendEntriesArgumentsFromProto converts a protobuf AppendEntriesArguments to its model representation.
func AppendEntriesArgumentsFromProto(p *pb.AppendEntriesArguments) AppendEntriesArguments {
	entries := make([]node.LogEntry, len(p.GetEntries()))
	for i, entry := range p.GetEntries() {
		entries[i] = node.LogEntryFromProto(entry)
	}

	return AppendEntriesArguments{
		Term:         p.GetTerm(),
		LeaderID:     p.GetLeaderId(),
		PrevLogIndex: p.GetPrevLogIndex(),
		PrevLogTerm:  p.GetPrevLogTerm(),
		Entries:      entries,
		LeaderCommit: p.GetLeaderCommit(),
	}
}

// ToProto converts AppendEntriesReply to its protobuf representation.
func (r AppendEntriesReply) ToProto() *pb.AppendEntriesReply {
	return &pb.AppendEntriesReply{
		Term:    r.Term,
		Success: r.Success,
	}
}

// AppendEntriesReplyFromProto converts a protobuf AppendEntriesReply to its model representation.
func AppendEntriesReplyFromProto(p *pb.AppendEntriesReply) AppendEntriesReply {
	return AppendEntriesReply{
		Term:    p.GetTerm(),
		Success: p.GetSuccess(),
	}
}
