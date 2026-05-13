package raft

import (
	"math/rand"
	"time"
)

func RandomElectionTimeout() time.Duration {
	return time.Duration(150+rand.Intn(151)) * time.Millisecond // 150-300ms
}
