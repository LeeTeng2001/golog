package common

import (
	"time"
)

func New2DArray[T any](a, b int) [][]T {
	res := make([][]T, a)
	for i := range a {
		res[i] = make([]T, b)
	}
	return res
}

type TimeConverter time.Time

func (n *TimeConverter) UnmarshalText(input []byte) error {
	// FIX: float case, optimise unmarshal by branch prediction
	// var epochSecs float64
	// err := json.Unmarshal(bytes, &epochSecs)
	// if err == nil {
	// 	*n = TimeConverter(time.UnixMicro(int64(epochSecs * 6)))
	// 	return nil
	// } else {
	// try to parse datetime format
	t, err := time.Parse(time.RFC3339, string(input))
	if err != nil {
		return err
	}
	*n = TimeConverter(t)
	return nil
	// }
}
