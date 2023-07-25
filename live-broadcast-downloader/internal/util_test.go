package internal

import (
	"github.com/stretchr/testify/assert"
	"go-live-broadcast-downloader/plugins/file"
	cjson "go-live-broadcast-downloader/plugins/json"
	"testing"
)

func TestParseM3U8(t *testing.T) {
	filename := "/Users/ayakurayuki/Desktop/index_6m.m3u8"
	res := ParseM3U8(filename)
	bs, _ := cjson.JSON.MarshalIndent(res, "", "    ")
	t.Log(string(bs))
}

func TestCreateFolder(t *testing.T) {
	path := "/Users/ayakurayuki/Desktop/content-114514"
	err := CreateFolder(path)
	if err != nil {
		t.Fatal(err)
	}
	exist, err := file.IsPathExist(path)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, exist, "create failed")
}
