package raft

// CommandType represents the type of a command in the log entry, which can be either SET or DELETE.
type CommandType string

// CommandType represents the type of a command in the log entry, which can be either SET or DELETE.
const (
	SET    CommandType = "SET"
	DELETE CommandType = "DELETE"
)

// Command represents a command to be executed on the state machine, including the type of command, the key it operates on, and an optional value for SET commands.
type Command struct {
	Type  CommandType
	Key   string
	Value *string // Value is optional, only needed for SET commands
}
