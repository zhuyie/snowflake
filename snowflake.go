package snowflake

import (
	"log"
	"runtime"
	"sync"
	"time"
)

const (
	maxNode   = 1023
	maxSeq    = 4095
	timeShift = 22
	nodeShift = 12

	twitterEpoch int64 = 1288834974657
)

// Snowflake is a unique ID generation algorithm which opensourced by Twitter.
type Snowflake struct {
	locker        sync.Mutex
	epoch         int64 // the epoch we used
	lastTime      time.Time
	lastTimestamp int64  // timestamp in milliseconds, 41 bits
	node          uint16 // NodeID, 10 bits
	sequence      uint16 // sequence, 12 bits
}

// NewSnowflake creates and initializes a new Snowflake instance.
func NewSnowflake(node uint16) *Snowflake {
	return NewSnowflakeEpoch(node, twitterEpoch)
}

// NewSnowflakeEpoch creates and initializes a new Snowflake instance.
func NewSnowflakeEpoch(node uint16, epoch int64) *Snowflake {
	if node > maxNode {
		panic("node should between 0 - 1023")
	}
	return &Snowflake{epoch: epoch, node: node}
}

// Next generates a unique 64-bit integer.
func (s *Snowflake) Next() int64 {
	s.locker.Lock()
	defer s.locker.Unlock()

	var waitingForLifeToGetBackToNormal bool
	for {
		t, ts := s.ts()
		if ts < s.lastTimestamp {
			if !waitingForLifeToGetBackToNormal {
				waitingForLifeToGetBackToNormal = true
				log.Printf("snowflake: time moved backwards: %v ms", s.lastTimestamp-ts)
			}
			time.Sleep(time.Millisecond)
			continue
		}
		waitingForLifeToGetBackToNormal = false
		if ts > s.lastTimestamp {
			s.lastTime = t
			s.lastTimestamp = ts
			s.sequence = 0
		}
		if s.sequence < maxSeq {
			ID := s.lastTimestamp<<timeShift | int64(s.node)<<nodeShift | int64(s.sequence)
			s.sequence++
			return ID
		}
		// Retry
		runtime.Gosched()
	}

	// Never reached
}
