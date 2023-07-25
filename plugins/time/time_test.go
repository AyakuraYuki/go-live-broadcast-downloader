package time

import (
	cjson "go-live-broadcast-downloader/plugins/json"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	now := Time{time.Now()}
	j, err := cjson.JSON.Marshal(now)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	t.Log(string(j))

	unmarshal := Time{}
	err = cjson.JSON.UnmarshalFromString(string(j), &unmarshal)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	t.Logf("%+v\n", unmarshal)
}
