package crypto

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"go-live-broadcast-downloader/plugins/log"
)

// Sha256 sha256
func Sha256(bs []byte) string {
	h := sha256.New()
	h.Write(bs)
	return hex.EncodeToString(h.Sum(nil))
}

// Sha512 sha512
func Sha512(bs []byte) string {
	h := sha512.New()
	h.Write(bs)
	return hex.EncodeToString(h.Sum(nil))
}

// Sha512Iter 自定义sha512 方法，迭代n 次
func Sha512Iter(input, salt []byte, iterations int) string {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Sha512Iter").Msgf("Sha512Iter panic recover,%v", r)
		}
	}()

	h := sha512.New()
	h.Write(salt)
	h.Write(input)
	result := h.Sum(nil)
	for i := 1; i < iterations; i++ {
		h.Reset()
		h.Write(result)
		result = h.Sum(nil)
	}
	return hex.EncodeToString(result)
}
