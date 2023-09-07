package utils

import (
	"fmt"
	"reflect"
)

// ConvertInterfaceToInt64 converts interface{} to int64.
//
//	@param value interface{}
//	@return int64
//	@return error
func ConvertInterfaceToInt64(value interface{}) (int64, error) {
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Float64:
		return int64(v.Float()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil
	default:
		return 0, fmt.Errorf("unhandled type: %s", v.Kind())
	}
}

// ConvertInterfaceToInt64Slice converts interface{} to []int64.
//
//	@param value interface{}
//	@return []int64
//	@return error
func ConvertInterfaceToInt64Slice(value interface{}) ([]int64, error) {
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Slice:
		slice := make([]int64, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			any := v.Index(i).Interface()
			if _, ok := any.(float64); ok {
				slice = append(slice, int64(any.(float64)))
			} else {
				slice = append(slice, any.(int64))
			}
		}

		return slice, nil
	default:
		return []int64{}, fmt.Errorf("unhandled type: %s", v.Kind())
	}
}
