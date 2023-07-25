package crypto

import (
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestGenerateRsaKey(t *testing.T) {
	privateKey, publicKey, err := GenerateRsaKey(1024)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(privateKey))
	t.Log(string(publicKey))
}

func TestRsaEncrypt(t *testing.T) {
	ctx := context.Background()
	privateKey, publicKey, _ := GenerateRsaKey(1024)
	cli := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6380",
	})

	cli.HMSet(ctx, "hs", map[string]string{
		"privateKey": string(privateKey),
		"publicKey":  string(publicKey),
	})

	publicKey2 := cli.HGet(ctx, "hs", "publicKey").Val()
	privateKey2 := cli.HGet(ctx, "hs", "privateKey").Val()

	t.Log("publicKey2", string(publicKey2) == string(publicKey))
	t.Log("privateKey2", string(privateKey2) == string(privateKey))

	bs, err := RsaEncrypt([]byte("aaa"), []byte(publicKey2))
	t.Log(err)
	t.Log(bs)
	t.Log("-----")
	bs2, err := RsaDecrypt([]byte(bs), []byte(privateKey2))
	t.Log(err)
	t.Log(bs2)
}

//func TestRsaEncrypt2(t *testing.T) {
//	ctx := context.Background()
//	d := dao.NewDao(config.GlobalConfig)
//
//	//// 生成RSA密钥对
//	//if privateKey, publicKey, err := GenerateRsaKey(2048); err == nil {
//	//	fields := map[string]interface{}{
//	//		dao.FieldPrivate: strings.TrimSpace(string(privateKey)),
//	//		dao.FieldPublic:  strings.TrimSpace(string(publicKey)),
//	//	}
//	//	d.SaveRsaKeyPair(ctx, fields, false)
//	//} else {
//	//	t.Fatal(err)
//	//	return
//	//}
//	//if privateKey, publicKey, err := GenerateRsaKey(2048); err == nil {
//	//	fields := map[string]interface{}{
//	//		dao.FieldPrivate: strings.TrimSpace(string(privateKey)),
//	//		dao.FieldPublic:  strings.TrimSpace(string(publicKey)),
//	//	}
//	//	d.SaveRsaKeyPair(ctx, fields, true)
//	//} else {
//	//	t.Fatal(err)
//	//	return
//	//}
//	//// 生成RSA密钥对结束
//
//	passwd := "dodoboys114514"
//	pub, _ := d.GetRsaPublicKey(ctx, false)
//	t.Logf("pub: %s\n", pub)
//	pri, _ := d.GetRsaPrivateKey(ctx)
//	t.Logf("pri: %s\n", pri)
//	enc, _ := RsaEncrypt([]byte(passwd), []byte(pub))
//	t.Logf("enc: %s\n", base64.StdEncoding.EncodeToString([]byte(enc)))
//	dec, _ := RsaDecrypt([]byte(enc), []byte(pri))
//	t.Logf("dec: %s\n", dec)
//
//	enc2 := "BypqiDLN2i06/+W1169qg04uU+Ls2fp77EUrl6LC/UGylFNnXKUsUs4rHmNupNR6ASXXFSVEI7sfDvKr842rIIkyYrPBKSVmCoPgURFNJSdN+3eVEPG2pm1YkxC7BIWjshp+UB0VylBFSPuW5u1VL7Q2VgbRf+suRRHOArUqIOY="
//	benc2, err := base64.StdEncoding.DecodeString(enc2)
//	if err != nil {
//		t.Fatal(err)
//	}
//	dec2, _ := RsaDecrypt(benc2, []byte(pri))
//	t.Log(dec2)
//}
