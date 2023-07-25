package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/pkg/errors"
	"go-live-broadcast-downloader/plugins/log"
)

// GenerateRsaKey 生成RSA密钥对
func GenerateRsaKey(bits int) (privateKey, publicKey []byte, err error) {
	prvKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Error("GenerateRsaKey").Msgf("prv: %v", err)
		return
	}
	// 私钥部分
	derPKCS1 := x509.MarshalPKCS1PrivateKey(prvKey)
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: derPKCS1}
	privateKey = pem.EncodeToMemory(block)
	// 公钥部分
	pubKey := &prvKey.PublicKey
	derPKIX, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		log.Error("GenerateRsaKey").Msgf("pub: %v", err)
		return
	}
	block = &pem.Block{Type: "PUBLIC KEY", Bytes: derPKIX}
	publicKey = pem.EncodeToMemory(block)
	return
}

// RsaEncrypt 加密
func RsaEncrypt(data []byte, pubKey []byte) (string, error) {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return "", errors.New("key err! ")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Error("ParsePKIXPublicKey").Msgf("%v", err)
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)
	rbs, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		log.Error("EncryptPKCS1v15").Msgf("%v", err)
		return "", err
	}
	return string(rbs), nil
}

// RsaDecrypt 解密
func RsaDecrypt(decodeBytes []byte, privateKey []byte) (string, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("key err! ")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Error("ParsePKCS1PrivateKey").Msgf("%v", err)
		return "", err
	}
	rbs, err := rsa.DecryptPKCS1v15(rand.Reader, priv, decodeBytes)
	if err != nil {
		log.Error("ParsePKCS1PrivateKey").Msgf("%v", err)
		return "", err
	}
	return string(rbs), nil
}
