// +build go1.9

package snowflake

import "time"

func (s *Snowflake) ts() (t time.Time, timestamp int64) {
	t = time.Now()
	d := int64(t.Sub(s.lastTime) / 1000000) // in milliseconds
	timestamp = s.lastTimestamp + d
	return
}
