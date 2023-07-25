package file

import (
	cjson "go-live-broadcast-downloader/plugins/json"
	"testing"
)

func TestListDir(t *testing.T) {
	path := "/Users/ayakurayuki/Desktop"
	list, err := ListDir(path)
	if err != nil {
		t.Fatal(err)
	}
	j, _ := cjson.JSON.MarshalIndent(list, "", "    ")
	t.Log(string(j))
}

func TestWalkDir(t *testing.T) {
	path := "/Users/ayakurayuki/go/src/ay-go-scaffolding/plugins"
	list := make([]string, 0)
	WalkDir(path, &list)
	j, _ := cjson.JSON.MarshalIndent(list, "", "    ")
	t.Log(string(j))
}
