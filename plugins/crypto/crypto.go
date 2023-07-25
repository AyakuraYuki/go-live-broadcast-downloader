package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"go-live-broadcast-downloader/plugins/hack"
	"hash/crc32"
	"io"
	"math/rand"
	"os"
	"time"
)

func Md5Str(str string) string {
	h := md5.New()
	h.Write(hack.Slice(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Md5Byte2Str(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func FileMd5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func Crc32(str string) uint32 {
	return crc32.ChecksumIEEE(hack.Slice(str))
}

const Lower = 1 << 0
const Upper = 1 << 1
const Digit = 1 << 2

const LowerUpper = Lower | Upper
const LowerDigit = Lower | Digit
const UpperDigit = Upper | Digit
const LowerUpperDigit = LowerUpper | Digit

const lower = "abcdefghijklmnopqrstuvwxyz"
const upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digit = "0123456789"
const symbol = "!#$%&()*+,-./:;<=>?@[]^_{|}~\\\"'"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RandString(size int, set int) string {
	charset := ""
	if set&Lower > 0 {
		charset += lower
	}
	if set&Upper > 0 {
		charset += upper
	}
	if set&Digit > 0 {
		charset += digit
	}

	lenAll := len(charset)

	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = charset[rand.Intn(lenAll)]
	}
	return string(buf)
}

func Sha1(data string) string {
	s := sha1.New()
	s.Write([]byte(data))
	return hex.EncodeToString(s.Sum([]byte("")))
}

func RandPassword(size int) string {
	charset := ""
	charset += lower
	charset += upper
	charset += digit
	charset += symbol
	lenAll := len(charset)
	if size < 8 {
		size = 8
	}
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		switch i {
		case 0:
			buf[i] = upper[rand.Intn(len(upper))]
			break
		case 1:
			buf[i] = lower[rand.Intn(len(lower))]
			break
		case 2:
			buf[i] = symbol[rand.Intn(len(symbol))]
			break
		case 3:
			buf[i] = digit[rand.Intn(len(digit))]
			break
		default:
			buf[i] = charset[rand.Intn(lenAll)]
		}
	}
	return string(buf)
}
