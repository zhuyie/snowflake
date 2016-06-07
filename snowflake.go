package snowflake

import (
	"errors"
	"runtime"
	"sync"
	"time"
)

const (
	maxNode   = 1023
	maxSeq    = 4095
	timeShift = 22
	nodeShift = 12
)

// Snowflake is a unique ID generation algorithm which opensourced by Twitter.
type Snowflake struct {
	locker   sync.Mutex
	lastTime uint64 // timestamp in microseconds, 41 bits
	node     uint16 // NodeID, 10 bits
	sequence uint16 // sequence, 12 bits
}

// NewSnowflake creates and initializes a new Snowflake instance.
func NewSnowflake(node uint16) (*Snowflake, error) {
	if node > maxNode {
		return nil, errors.New("InvalidArgs: node should between 0 - 1023")
	}
	s := &Snowflake{node: node}
	return s, nil
}

// Next generates a unique 64-bit integer.
func (s *Snowflake) Next() uint64 {
	s.locker.Lock()
	defer s.locker.Unlock()

	for {
		now := uint64(time.Now().UnixNano() / 1000) // in microseconds
		if now < s.lastTime {
			panic("timestamp rollback")
		}
		if now > s.lastTime {
			s.lastTime = now
			s.sequence = 0
		}
		if s.sequence < maxSeq {
			s.sequence++
			return s.lastTime<<timeShift | uint64(s.node)<<nodeShift | uint64(s.sequence)
		}
		// Retry
		runtime.Gosched()
	}

	// Never reached
}
