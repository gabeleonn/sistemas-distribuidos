package raft

import "goraft/proto"

type AppendEntriesRequest struct {
	Term         int64
	LeaderID     int64
	PrevLogIndex int64
	PrevLogTerm  int64
	Entries      []LogEntry
	LeaderCommit int64
}

type AppendEntriesResponse struct {
	Term    int64
	Success bool
}

func (ae *AppendEntriesRequest) ToProto() *proto.AppendEntriesRequest {
	entriesProto := make([]*proto.LogEntry, len(ae.Entries))
	for i, entry := range ae.Entries {
		entriesProto[i] = entry.ToProto()
	}

	return &proto.AppendEntriesRequest{
		Term:         ae.Term,
		LeaderId:     ae.LeaderID,
		PrevLogIndex: ae.PrevLogIndex,
		PrevLogTerm:  ae.PrevLogTerm,
		Entries:      entriesProto,
		LeaderCommit: ae.LeaderCommit,
	}
}

func AppendEntriesRequestFromProto(p *proto.AppendEntriesRequest) *AppendEntriesRequest {
	entries := make([]LogEntry, len(p.Entries))
	for i, entryProto := range p.Entries {
		entry, err := LogEntryFromProto(entryProto)
		if err != nil {
			continue
		}
		entries[i] = *entry
	}

	return &AppendEntriesRequest{
		Term:         p.Term,
		LeaderID:     p.LeaderId,
		PrevLogIndex: p.PrevLogIndex,
		PrevLogTerm:  p.PrevLogTerm,
		Entries:      entries,
		LeaderCommit: p.LeaderCommit,
	}
}

func (ae *AppendEntriesResponse) ToProto() *proto.AppendEntriesResponse {
	return &proto.AppendEntriesResponse{
		Term:    ae.Term,
		Success: ae.Success,
	}
}

func AppendEntriesResponseFromProto(p *proto.AppendEntriesResponse) *AppendEntriesResponse {
	return &AppendEntriesResponse{
		Term:    p.Term,
		Success: p.Success,
	}
}

func NewAppendEntriesRequest(
	term int64,
	leaderID int64,
	prevLogIndex int64,
	prevLogTerm int64,
	entries []LogEntry,
	leaderCommit int64,
) *AppendEntriesRequest {
	return &AppendEntriesRequest{
		Term:         term,
		LeaderID:     leaderID,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      entries,
		LeaderCommit: leaderCommit,
	}
}
