package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"go-live-broadcast-downloader/plugins/log"
)

// AESEncrypt 加密数据
func AESEncrypt(data, key, iv []byte) ([]byte, error) {
	if len(data) == 0 {
		log.Error("AESEncrypt").Msgf("data len=0")
		return nil, errors.New("data len=0")
	}
	defer func() {
		if r := recover(); r != nil {
			log.Error("AESEncrypt").Msgf("AESDecrypt panic recover")
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error("AESEncrypt").Msgf("%v", err)
		return nil, err
	}
	content := PKCS5Padding(data, block.BlockSize())
	encrypted := make([]byte, len(content))

	if len(iv) != block.BlockSize() {
		log.Error("AESEncrypt").Msg("cipher CBC: IV length must equal block size")
		return nil, errors.New("decrypt error -3")
	}

	aesDecrypt := cipher.NewCBCEncrypter(block, iv)
	aesDecrypt.CryptBlocks(encrypted, content)
	return encrypted, nil
}

// AESDecrypt 解密数据
func AESDecrypt(src, key, iv []byte) (data []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("AESEncrypt").Msgf("AESDecrypt panic recover")
		}
	}()
	if len(src) == 0 {
		log.Error("AESEncrypt").Msgf("data len=0")
		return nil, errors.New("data len=0")
	}
	decrypted := make([]byte, len(src))
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error("AESDecrypt").Msgf("%v", err)
		return nil, err
	}
	if len(iv) != block.BlockSize() {
		log.Error("AESDecrypt").Msg("cipher CBC: IV length must equal block size")
		return nil, errors.New("decrypt error -2")
	}
	aesDecrypt := cipher.NewCBCDecrypter(block, iv)
	aesDecrypt.CryptBlocks(decrypted, src)
	return PKCS5Trimming(decrypted), nil
}

// PKCS5Padding PKCS5包装
func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// PKCS5Trimming 解包装
func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
