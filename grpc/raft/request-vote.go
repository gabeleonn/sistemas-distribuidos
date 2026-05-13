package raft

import "goraft/proto"

type RequestVoteRequest struct {
	Term         int64
	CandidateID  int64
	LastLogIndex int64
	LastLogTerm  int64
}

type RequestVoteResponse struct {
	Term        int64
	VoteGranted bool
}

func (rv *RequestVoteRequest) ToProto() *proto.RequestVoteRequest {
	return &proto.RequestVoteRequest{
		Term:         rv.Term,
		CandidateId:  rv.CandidateID,
		LastLogIndex: rv.LastLogIndex,
		LastLogTerm:  rv.LastLogTerm,
	}
}

func RequestVoteRequestFromProto(p *proto.RequestVoteRequest) *RequestVoteRequest {
	return &RequestVoteRequest{
		Term:        p.Term,
		CandidateID: p.CandidateId,
	}
}

func (rv *RequestVoteResponse) ToProto() *proto.RequestVoteResponse {
	return &proto.RequestVoteResponse{
		Term:        rv.Term,
		VoteGranted: rv.VoteGranted,
	}
}

func RequestVoteResponseFromProto(p *proto.RequestVoteResponse) *RequestVoteResponse {
	return &RequestVoteResponse{
		Term:        p.Term,
		VoteGranted: p.VoteGranted,
	}
}
