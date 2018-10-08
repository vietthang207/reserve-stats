package common

import (
	"math/big"
	"time"
)

// GetTimepoint return current Unix Timestamp in millisecond with uint64 format
func GetTimepoint() int64 {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return int64(timestamp)
}

// BigToFloat converts a big int to float according to its number of decimal digits
// Example:
// - BigToFloat(1100, 3) = 1.1
// - BigToFloat(1100, 2) = 11
// - BigToFloat(1100, 5) = 0.11
func BigToFloat(b *big.Int, decimal int64) float64 {
	f := new(big.Float).SetInt(b)
	power := new(big.Float).SetInt(new(big.Int).Exp(
		big.NewInt(10), big.NewInt(decimal), nil,
	))
	res := new(big.Float).Quo(f, power)
	result, _ := res.Float64()
	return result
}

// TimepointMillisecToTime convert a timepoint in s to time.Time object
func TimepointMillisecToTime(t int64) time.Time {
	return time.Unix(0, int64(t)*int64(time.Millisecond))
}
