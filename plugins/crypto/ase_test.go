package crypto

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAESDecrypt(t *testing.T) {
	data := `abcd`
	fmt.Println(data)
	fmt.Println("base64 :", base64.StdEncoding.EncodeToString([]byte(data)))
	sign2 := "3d2d8d696b21d7735027be9009a3d2ed"
	iv := "5f2d8cf01d4f0be2"
	res, _ := AESEncrypt([]byte(data), []byte(sign2), []byte(iv))
	fmt.Println("md5:", Md5Str(string(res)), len(res))
	fmt.Println("加密:", base64.StdEncoding.EncodeToString(res))

	ori, err := AESDecrypt(res, []byte(sign2), []byte(iv))
	fmt.Println("解密:", len(ori), string(ori))
	fmt.Println(err)
}

func TestAESDecrypt2(t *testing.T) {
	data := `{"type":2,"data":{"checkCode":"02a81d5ff0b54c5aa14fc50634ddae75"},"version":"v2"}`
	key := `03ebc3aed25e658b1f9eaf1de78ffdfd`
	hKey, _ := hex.DecodeString(key)

	enc, _ := AESEncrypt([]byte(data), hKey, make([]byte, 16))
	hEnc := hex.EncodeToString(enc) // dfdce0760de51c2460f72aa0273acae3e55d12c96ad2be7870b1bbc3ec372a5c5d4073b6b0d02fe0f4a206f0c471b8b549b9eec57fceb864df0586c2300d6c60a8dacffa1026fcc87221ccd228595e95ab5aa4c0e4fd3eed6aa80068c5bd169b
	t.Log(hEnc)

	hhEnc, _ := hex.DecodeString(hEnc)
	dec, _ := AESDecrypt(hhEnc, hKey, make([]byte, 16))
	t.Log(string(dec))
}
