package crypto

import (
	"testing"
)

func TestMd5Str(t *testing.T) {
	if Md5Str("Shanghai") != "5466ee572bcbc75830d044e66ab429bc" {
		t.Fatal("should equal")
	}
}

func TestCrc32(t *testing.T) {
	if Crc32("Shanghai") != 1271261733 {
		t.Fatal("should equal")
	}
}

func TestRandString(t *testing.T) {
	r := RandString(12, LowerUpperDigit)

	if len(r) != 12 {
		t.Fatal("rand str should 12 len")
	}

	t.Log(r)
}

func TestRandString2(t *testing.T) {
	r := RandString(16, LowerUpperDigit)
	if len(r) != 16 {
		t.Fatal("rand str should 16 len")
	}
	t.Log(r)
}

func BenchmarkRandString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandString(12, LowerUpperDigit)
	}
}
