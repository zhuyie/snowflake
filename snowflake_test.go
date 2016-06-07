package snowflake

import (
	"testing"
	"time"
)

func TestSnowflake(t *testing.T) {
	sf0, _ := NewSnowflake(0)
	sf1, _ := NewSnowflake(1)
	map0 := make(map[uint64]bool)
	map1 := make(map[uint64]bool)

	for i := 0; i < 2000; i++ {
		v0 := sf0.Next()
		v1 := sf0.Next()
		v2 := sf1.Next()
		v3 := sf1.Next()

		if map0[v0] {
			t.Errorf("v0(%x) already exist in map0", v0)
		}
		if map0[v1] {
			t.Errorf("v1(%x) already exist in map0", v1)
		}
		map0[v0] = true
		map0[v1] = true

		if map1[v2] {
			t.Errorf("v2(%x) already exist in map1", v2)
		}
		if map1[v3] {
			t.Errorf("v3(%x) already exist in map1", v3)
		}
		map1[v2] = true
		map1[v3] = true

		if map0[v2] {
			t.Errorf("v2(%x) should not exist in map0", v2)
		}
		if map0[v3] {
			t.Errorf("v3(%x) should not exist in map0", v3)
		}

		time.Sleep(time.Millisecond)
	}
}

var result uint64

func BenchmarkSnowflake(b *testing.B) {
	sf, _ := NewSnowflake(99)
	var v uint64
	for n := 0; n < b.N; n++ {
		// record the result prevent the compiler eliminating the function call.
		v = sf.Next()
	}
	// store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = v
}
