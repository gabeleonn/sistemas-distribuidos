package messages

import (
	"testing"
)

func TestRequestVoteArgumentsToProtoAndFromProto(t *testing.T) {
	args := RequestVoteArguments{
		Term:         1,
		CandidateID:  2,
		LastLogIndex: 3,
		LastLogTerm:  4,
	}

	protoArgs := args.ToProto()
	if protoArgs.GetTerm() != args.Term {
		t.Errorf("term mismatch: got %d, want %d", protoArgs.GetTerm(), args.Term)
	}
	if protoArgs.GetCandidateId() != args.CandidateID {
		t.Errorf("candidate ID mismatch: got %d, want %d", protoArgs.GetCandidateId(), args.CandidateID)
	}
	if protoArgs.GetLastLogIndex() != args.LastLogIndex {
		t.Errorf("last log index mismatch: got %d, want %d", protoArgs.GetLastLogIndex(), args.LastLogIndex)
	}
	if protoArgs.GetLastLogTerm() != args.LastLogTerm {
		t.Errorf("last log term mismatch: got %d, want %d", protoArgs.GetLastLogTerm(), args.LastLogTerm)
	}

	modelArgs := RequestVoteArgumentsFromProto(protoArgs)
	if modelArgs != args {
		t.Errorf("model arguments mismatch: got %+v, want %+v", modelArgs, args)
	}
}

func TestRequestVoteReplyToProtoAndFromProto(t *testing.T) {
	reply := RequestVoteReply{
		Term:        1,
		VoteGranted: true,
	}

	protoReply := reply.ToProto()
	if protoReply.GetTerm() != reply.Term {
		t.Errorf("term mismatch: got %d, want %d", protoReply.GetTerm(), reply.Term)
	}
	if protoReply.GetVoteGranted() != reply.VoteGranted {
		t.Errorf("vote granted mismatch: got %v, want %v", protoReply.GetVoteGranted(), reply.VoteGranted)
	}

	modelReply := RequestVoteReplyFromProto(protoReply)
	if modelReply != reply {
		t.Errorf("model reply mismatch: got %+v, want %+v", modelReply, reply)
	}
}
