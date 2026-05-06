package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"raft/models/log"
	"sync"
)

// NodePersistentState represents the persistent state of a Raft node that needs to be saved to disk.
type NodePersistentState struct {
	CurrentTerm int64
	VotedFor    *int64
	Log         []log.LogEntry
}

var lock sync.RWMutex

// SavePersistentNodeState saves the given NodePersistentState to the specified path. This function is thread-safe.
func SavePersistentNodeState(filepath string, value NodePersistentState) error {
	lock.Lock()
	defer lock.Unlock()
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := marshal(value)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, &data)
	return err
}

// LoadPersistentNodeState loads the NodePersistentState from the specified path. This function is thread-safe.
func LoadPersistentNodeState(filepath string) (NodePersistentState, error) {
	lock.Lock()
	defer lock.Unlock()
	file, err := os.Open(filepath)
	if err != nil {
		return NodePersistentState{}, err
	}
	defer file.Close()
	data := bytes.Buffer{}
	_, err = io.Copy(&data, file)
	if err != nil {
		return NodePersistentState{}, err
	}
	return unmarshal(data)
}

func unmarshal(data bytes.Buffer) (NodePersistentState, error) {
	var value NodePersistentState
	buf := bytes.NewBuffer(data.Bytes())
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&value)
	if err != nil {
		return NodePersistentState{}, err
	}
	return value, nil
}

func marshal(value NodePersistentState) (bytes.Buffer, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(value)
	if err != nil {
		return bytes.Buffer{}, err
	}

	return buf, nil
}

// GetNodePath returns the file path for storing the persistent state of a node with the given ID.
func GetNodePath(nodeID int64) string {
	return fmt.Sprintf("%sgoraft_node_%d.gob", os.TempDir(), nodeID)
}
