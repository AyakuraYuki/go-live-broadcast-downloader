package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

// HmacSHA1 HmacSHA1
func HmacSHA1(data, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	h := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(h)
}
