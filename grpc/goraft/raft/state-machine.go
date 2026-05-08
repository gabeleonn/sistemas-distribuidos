package raft

type StateMachine interface {
	Apply(cmd Command) error
}
