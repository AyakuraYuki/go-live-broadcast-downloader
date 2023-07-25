package crypto

import (
	"crypto/sha256"
	"testing"
)

func TestRC4(t *testing.T) {
	str := "中文"
	data, err := RC4([]byte(str), []byte("adsxMZWdWUyWjQyyvOWnOn7wGXSKUl6C"))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("rc4:", sha256.New().Sum(data))
	t.Log("rc4:", data)
}

func BenchmarkRC4(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RC4([]byte("12345678"), []byte("123abc"))
	}
}

func TestRc4EncodeBase64(t *testing.T) {

}

func TestRc4DecodeBase64(t *testing.T) {
	str := "qjqyYoEBAlQyvRFpBnGxhj0PpTpj6vJZ1gVTNcMNsHhn3pD8Do5Uj8GkgGF8t9nE2mYhGo15jirpVuLig/NNqJTvHezRFv3IPGMP1B5hocV9uG7FZkUJw1S1MSavZ5OsnVHAwvkOJVdDpmT76cFfJ1wQuM/SUCIRCs8oqvZgdPaN3I2qKk8EFFi+I71CQ0xxVdpbkWwRV+WPEYp9CzqyuVVBoirPqcSHBPdiEi1PlotqCqBfU+Wu09mbCL0cxUjLJBe0eFeTBXE7iFcG+VnauRmy5XViBwHhgqGH/ioiH3dcFjcW7mkgRF1QrfvHo6OLLHggItfz5Mc6PUqGU0QRezvf6EwCqMcuGd5pCzSyYs2WVVrCf333R7NNgdVUZwjhrLwIxNeyzCAul1wEJt886RlDuNlInxHMf4EL5GOzPloEiB02QC6nL26Av7syh2+kd0veZmpXJSCEZDeOkfHBPkpJI12cIuihd6d2vUaUejloUstu2DQTwaro072lcod5t5kFGbTMu7yUmob1fpoQwFIyNwA5DfHLv/Nr3hv40FdV3X7HaoRcnW+zPXxbRNw9BtDvjhbcQ7IkPhl0erMb8G9CKdo9UKGUGX4MqZ02FoK8+0wPcTtsOj2wdmbQhrhOoLF/JN4+38QKk39DtQ=="
	key := "dodo-admin"
	s, err := Rc4DecodeBase64(str, []byte(key))
	t.Log(string(s))
	t.Log(err)
}
