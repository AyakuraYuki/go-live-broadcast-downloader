package crypto

import (
	"encoding/hex"
	"testing"
)

func TestSha512Iter(t *testing.T) {
	txt := "LiMin Dump"
	salt, _ := hex.DecodeString("75d67175f3c60a6f91dbe68cf655cfce")

	t.Log(txt)
	t.Log(salt)
	result := Sha512Iter([]byte(txt), salt, 1024)
	t.Log(result)
}
