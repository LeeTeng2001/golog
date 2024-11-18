package common

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type Float64Unmarshaler interface {
	UnmarshalFloat64(input float64) error
}

func FloatUnmarshallerHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Float64 {
			return data, nil
		}
		result := reflect.New(t).Interface()
		unmarshaller, ok := result.(Float64Unmarshaler)
		if !ok {
			return data, nil
		}
		if err := unmarshaller.UnmarshalFloat64(data.(float64)); err != nil {
			return nil, err
		}
		return result, nil
	}
}
