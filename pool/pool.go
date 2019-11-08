package pool

import (
	"sync"
)

var bytePool sync.Pool

func GetBytesLen(length uint32) []byte {
	if v := bytePool.Get(); v != nil {
		e := v.([]byte)
		e = e[:0]
		for i := uint32(0); i < length; i++ {
			e = append(e, 0)
		}
		return e
	}
	e := make([]byte, length)
	return e
}
func PutBytes(e []byte) {
	bytePool.Put(e)
}
