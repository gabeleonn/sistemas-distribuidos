package raft

type CommandType string

const (
	SET CommandType = "SET"
	DEL CommandType = "DEL"
)

type Command struct {
	Type  CommandType
	Key   string
	Value *string // Value is optional, only needed for SET commands
}
