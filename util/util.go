package util

import "reflect"

func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}
	return reflect.ValueOf(v).IsNil()
}
