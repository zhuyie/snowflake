package main

import(
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %v SnowflakeID\n", os.Args[0])
		os.Exit(2)
	}

	u, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		fmt.Printf("Invalid ID\n")
		os.Exit(1)
	}

	ts := u >> 22             // 41 bits
	node := (u >> 10) & 0x3ff // 10 bits
	sequence := u & 0xfff     // 12 bits

	unixTime := int64(ts + 1288834974657)  // twitterEpoch
	t := time.Unix(unixTime / 1000, (unixTime % 1000) * 1000 * 1000)

	fmt.Printf("ID       : %v\n", u)
	fmt.Printf("Time     : %v\n", t)
	fmt.Printf("Node     : %v\n", node)
	fmt.Printf("Sequence : %v\n", sequence)
}
