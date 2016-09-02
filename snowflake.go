package snowflake

import (
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

	maxMoveBackwards int64 = 3000 // milliseconds
)

// Snowflake is a unique ID generation algorithm which opensourced by Twitter.
type Snowflake struct {
	locker   sync.Mutex
	epoch    int64  // the epoch we used
	lastTime int64  // timestamp in milliseconds, 41 bits
	node     uint16 // NodeID, 10 bits
	sequence uint16 // sequence, 12 bits
}

// NewSnowflake creates and initializes a new Snowflake instance.
func NewSnowflake(node uint16) *Snowflake {
	if node > maxNode {
		panic("node should between 0 - 1023")
	}
	return &Snowflake{epoch: twitterEpoch, node: node}
}

// NewSnowflakeEpoch creates and initializes a new Snowflake instance.
func NewSnowflakeEpoch(node uint16, epoch int64) *Snowflake {
	if node > maxNode {
		panic("node should between 0 - 1023")
	}
	return &Snowflake{epoch: epoch, node: node}
}

func (s *Snowflake) timeGen() int64 {
	now := time.Now().UnixNano() / 1000000 // in milliseconds
	if now < s.epoch {
		panic("Time before Epoch")
	}
	return now - s.epoch
}

// Next generates a unique 64-bit integer.
func (s *Snowflake) Next() int64 {
	s.locker.Lock()
	defer s.locker.Unlock()

	for {
		now := s.timeGen()
		if now < s.lastTime {
			if s.lastTime - now <= maxMoveBackwards {
				time.Sleep(time.Millisecond)
				continue
			}
			panic("Time moved backwards")
		}
		if now > s.lastTime {
			s.lastTime = now
			s.sequence = 0
		}
		if s.sequence < maxSeq {
			ID := s.lastTime<<timeShift | int64(s.node)<<nodeShift | int64(s.sequence)
			s.sequence++
			return ID
		}
		// Retry
		runtime.Gosched()
	}

	// Never reached
}
