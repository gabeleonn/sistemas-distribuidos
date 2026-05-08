package raft

// StateMachine defines the interface for a state machine that can apply commands from the Raft log.
// Implementations of this interface will execute the logic for handling SET and DELETE commands on the state machine's state.
type StateMachine interface {
	Apply(cmd Command) error
}
