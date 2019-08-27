package rslice

import "github.com/wudiliujie/common/convert"

func Join(data []int32, sep string) string {
	ret := ""
	for i, v := range data {
		if i > 0 {
			ret += sep
		}
		ret += convert.ToString(v)
	}
	return ret
}
func Join64(data []int64, sep string) string {
	ret := ""
	for i, v := range data {
		if i > 0 {
			ret += sep
		}
		ret += convert.ToString(v)
	}
	return ret
}
