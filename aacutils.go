package aac

import "math"

func MinInt(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

// Returns the maximum of both ints
func MaxInt(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func AacRound(x float64) int {
	return int(math.Floor(x + 0.5))
}
