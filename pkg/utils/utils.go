package utils

import "unsafe"

func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
