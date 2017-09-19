// +build !go1.9

package snowflake

import "time"

func (s *Snowflake) ts() (t time.Time, timestamp int64) {
	t = time.Now()
	tInMS := t.UnixNano() / 1000000 // in milliseconds
	if tInMS < s.epoch {
		panic("Time before Epoch")
	}
	timestamp = tInMS - s.epoch
	return
}
