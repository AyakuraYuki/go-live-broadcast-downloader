package nhttp

import (
	cjson "github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/json"
	"net/http"
	"testing"
)

func debugShowResponseResults(t *testing.T, data []byte, rspHeader http.Header, httpCode int) {
	t.Logf("http code: %v\n", httpCode)
	t.Logf("response data: %v\n", string(data))
	bs, _ := cjson.JSON.Marshal(rspHeader)
	t.Logf("headers: %v\n", string(bs))
}

func TestGetRaw(t *testing.T) {
	requestUrl := "https://v4.ipip.net/"
	data, header, code, err := GetRaw(nil, requestUrl, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	debugShowResponseResults(t, data, header, code)
}

func TestPostRaw(t *testing.T) {
	requestUrl := "https://apis.imdodo.com/island/v3/islands/10000"
	data, header, code, err := PostRaw(nil, requestUrl, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	debugShowResponseResults(t, data, header, code)
}

func TestHead(t *testing.T) {
	requestUrl := "https://www.bing.com/"
	header, code, err := Head(nil, requestUrl, nil, 10000, 3)
	if err != nil {
		t.Fatal(err)
	}
	debugShowResponseResults(t, nil, header, code)
}
