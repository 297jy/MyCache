package common

import "unsafe"

func Max(x int, y int) int {
	if x < y {
		return y
	} else {
		return x
	}
}

func ByteToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
