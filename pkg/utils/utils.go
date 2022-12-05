package utils

import (
	"bytes"
	"math/rand"
	"unsafe"
)

const abs = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandomString(length int) string {
	var buff bytes.Buffer
	for i := 0; i < length; i++ {
		buff.WriteByte(abs[rand.Intn(len(abs))])
	}

	return buff.String()
}

func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
