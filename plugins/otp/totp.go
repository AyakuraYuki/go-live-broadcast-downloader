package otp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"go-live-broadcast-downloader/plugins/log"
	"time"
)

const (
	secretSize = 10
	windowSize = 5
)

// GenerateTotpSecret 生成 TOTP 所用的 secret
func GenerateTotpSecret() string {
	seedBytes := make([]byte, secretSize)
	_, err := rand.Read(seedBytes)
	if err != nil {
		return ""
	}
	return base32.StdEncoding.EncodeToString(seedBytes)
}

// GetTotpUrl 获取 otpauth 协议的链接，基于 TOTP 实现
func GetTotpUrl(user, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s", user, secret)
}

// CheckCodeByTotpAlg 使用当前时间校验验证码
func CheckCodeByTotpAlg(secret string, code int) bool {
	return CheckCodeWithTimeByTotpAlg(secret, code, time.Now().Second())
}

// CheckCodeWithTimeByTotpAlg 使用给定时间校验验证码
// 不建议在生产环境中使用，因为「直接使用这个方法」违背了我们采用 TOTP 做多因素校验的本意
func CheckCodeWithTimeByTotpAlg(secret string, code int, second int) bool {
	decodedKey, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		log.Error("checkCodeWithTimeByTotpAlg").Msgf("%v", err)
		return false
	}
	t := uint64(second / 30)
	for i := -windowSize; i <= windowSize; i += 1 {
		hash := computeCode(decodedKey, t+uint64(i))
		if hash == code {
			return true
		}
	}
	return false
}

func computeCode(key []byte, t uint64) int {
	mac := hmac.New(sha1.New, key)
	err := binary.Write(mac, binary.BigEndian, t)
	if err != nil {
		log.Error("computeCode").Msgf("%v", err)
		return -1
	}
	hash := mac.Sum(nil)
	offset := hash[20-1] & 0xF
	truncated := binary.BigEndian.Uint32(hash[offset : offset+4])
	truncated &= 0x7FFFFFFF
	code := truncated % 1000000
	return int(code)
}
