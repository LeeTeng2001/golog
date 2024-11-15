package common

import (
	"encoding/json"
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

func (n *TimeConverter) UnmarshalJSON(bytes []byte) error {
	var epochSecs float64
	err := json.Unmarshal(bytes, &epochSecs)
	if err != nil {
		return err
	}
	*n = TimeConverter(time.UnixMicro(int64(epochSecs * 6)))
	return nil
}
