package jsonreq

import (
	"github.com/gin-gonic/gin"
	cjson "go-live-broadcast-downloader/plugins/json"
	"go-live-broadcast-downloader/plugins/log"
)

// JSONDecode json 解析
func JSONDecode(c *gin.Context, val interface{}) (err error) {
	if val == nil {
		return nil
	}
	raw, err := c.GetRawData()
	if err != nil {
		log.Error("JSONDecode").Msg("%v", err)
		return
	}

	if err = cjson.JSON.Unmarshal(raw, val); err != nil {
		log.Error("JSONDecode").Msg("%v", err)
	}

	return err
}
