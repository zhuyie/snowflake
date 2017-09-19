package snowflake

import (
	"fmt"
	"testing"
	"time"
)

func TestSnowflake(t *testing.T) {
	sf0 := NewSnowflake(0)
	sf1 := NewSnowflake(1)
	map0 := make(map[int64]bool)
	map1 := make(map[int64]bool)

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

func TestInvalidNode(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("TestInvalidNode panic: %v\n", err)
		}
	}()
	NewSnowflake(4096)
}

func TestEpoch(t *testing.T) {
	sf0 := NewSnowflakeEpoch(0, 100)
	v0 := sf0.Next()
	v1 := sf0.Next()
	if v0 == v1 {
		t.Fatalf("v0 equals v1, should different")
	}
}

func TestEpochInvalidNode(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("TestEpochInvalidNode panic: %v\n", err)
		}
	}()
	NewSnowflakeEpoch(4096, 100)
}

func TestVeryLargeEpoch(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("TestVeryLargeEpoch panic: %v\n", err)
		}
	}()
	sf0 := NewSnowflakeEpoch(0, 9223372036854770000)
	sf0.Next()
}

func TestTimeMovedBackwards(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("TestTimeMovedBackwards panic: %v\n", err)
		}
	}()
	sf0 := NewSnowflake(0)
	sf0.lastTime, sf0.lastTimestamp = sf0.ts()
	sf0.lastTime = sf0.lastTime.Add(time.Second)
	sf0.lastTimestamp += 1000
	sf0.Next()
}

func TestReachedMaxSequence(t *testing.T) {
	sf0 := NewSnowflake(0)
	for i := 0; i < 10000; i++ {
		v0 := sf0.Next()
		v1 := sf0.Next()
		if v0 == v1 {
			t.Fatalf("v0 equals v1, should different")
		}
	}
}

var result int64

func BenchmarkSnowflake(b *testing.B) {
	sf := NewSnowflake(99)
	var v int64
	for n := 0; n < b.N; n++ {
		// record the result prevent the compiler eliminating the function call.
		v = sf.Next()
	}
	// store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = v
}
