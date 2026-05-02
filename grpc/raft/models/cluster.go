package models

type NodeStatus int

const (
	Stopped NodeStatus = iota
	Running
	Unreachable
)

type Cluster struct {
	Nodes map[int64]Metadata
}

type ClusterState struct {
	Nodes map[int64]NodeState
}

type NodeState struct {
	ID          int64
	Addr        string
	Role        NodeRole
	Term        int64
	LeaderID    *int64
	CommitIndex int64
	LastLogIndex int64
	LastApplied int64
	Status      NodeStatus
}
