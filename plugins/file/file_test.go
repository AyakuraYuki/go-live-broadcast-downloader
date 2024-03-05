package file

import (
	cjson "github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/json"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestListDir(t *testing.T) {
	userHome, err := os.UserHomeDir()
	assert.NoError(t, err)
	list, err := ListDir(userHome)
	assert.NoError(t, err)
	j, _ := cjson.JSON.MarshalIndent(list, "", "    ")
	t.Log(string(j))
}

func TestWalkDir(t *testing.T) {
	userHome, err := os.UserHomeDir()
	assert.NoError(t, err)
	list := make([]string, 0)
	WalkDir(filepath.Join(userHome, "Documents"), &list)
	j, _ := cjson.JSON.MarshalIndent(list, "", "    ")
	t.Log(string(j))
}
