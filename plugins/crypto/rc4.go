package crypto

import (
	"crypto/rc4"
	"encoding/base64"
	"go-live-broadcast-downloader/plugins/log"
)

// RC4 rc4 xor
func RC4(src, key []byte) ([]byte, error) {
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		log.Error("RC4").Msgf("%v", err)
		return nil, err
	}

	dst := make([]byte, len(src))
	cipher.XORKeyStream(dst, src)
	return dst, nil
}

// Rc4EncodeBase64 rc4 加密后转base64
func Rc4EncodeBase64(src, key []byte) (string, error) {
	bs, err := RC4(src, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bs), nil
}

// Rc4DecodeBase64 base64 decode后解rc4
func Rc4DecodeBase64(base64Str string, key []byte) ([]byte, error) {
	bs, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		log.Error("Rc4DecodeBase64").Msgf("%v", err)
		return []byte{}, err
	}
	src, err := RC4(bs, key)
	if err != nil {
		return []byte{}, err
	}
	return src, nil
}
