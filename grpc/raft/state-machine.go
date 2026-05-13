package raft

type StateMachine interface {
	Apply(cmd Command) error
	Get(key string) (string, bool)
}
