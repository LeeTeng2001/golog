package common

import (
	"errors"
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
	// try to parse datetime format
	t, err := time.Parse(time.RFC3339, string(input))
	if err == nil {
		*n = TimeConverter(t)
		return nil
	}
	return errors.New("failed to parse time in string or float format")
}

func (n *TimeConverter) UnmarshalFloat64(input float64) error {
	*n = TimeConverter(time.UnixMicro(int64(input * 1000000)))
	return nil
}
