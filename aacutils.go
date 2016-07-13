package gaad

import "math"

func minInt(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

// Returns the maximum of both ints
func maxInt(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func aacRound(x float64) int {
	return int(math.Floor(x + 0.5))
}
