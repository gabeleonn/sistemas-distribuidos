package raft

// LogEntry represents an entry in the Raft log, containing the index, term, and command to be executed on the state machine.
type LogEntry struct {
	Index   int64
	Term    int64
	Command Command
}
