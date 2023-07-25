package jsonrsp

import (
	"bytes"
	"compress/gzip"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	ay_const "go-live-broadcast-downloader/plugins/ay-const/trace"
	"go-live-broadcast-downloader/plugins/errcode"
	cjson "go-live-broadcast-downloader/plugins/json"
	"go-live-broadcast-downloader/plugins/log"
	"io"
	"net/http"
	"strings"
)

const (
	respCtxKey = "_respCtxKey"
)

type Response struct {
	Status  int64       `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"traceId"`
}

func ApiStatusCodeError(c *gin.Context, status errcode.ServerError) {
	ApiResponse(c, status, status.ErrMsg, nil)
}

func ApiError(c *gin.Context, status errcode.ServerError, message string) {
	ApiResponse(c, status, message, nil)
}

func ApiResponse(c *gin.Context, status errcode.ServerError, message string, data interface{}) {
	rsp := &Response{
		Status:  status.ErrCode,
		Message: message,
		Data:    data,
		TraceId: c.GetString(ay_const.TraceCtxKey),
	}
	if message == "" {
		rsp.Message = status.ErrMsg
	}
	if data == nil {
		rsp.Data = struct{}{}
	}

	bs, err := cjson.JSON.Marshal(rsp)
	if err != nil {
		log.Error("ApiResponse").Msgf("%v", err)
	}
	c.Set(respCtxKey, string(bs))

	// Compress (br -> gzip -> non-compress)
	if len(bs) > 1000 {
		var (
			encodings  = c.GetHeader("Accept-Encoding")
			encoding   string
			buf        = new(bytes.Buffer)
			compressor io.WriteCloser
		)
		if strings.Contains(encodings, "br") {
			encoding = "br"
			compressor = brotli.NewWriterLevel(buf, brotli.DefaultCompression)
		} else if strings.Contains(encodings, "gzip") {
			encoding = "gzip"
			compressor, err = gzip.NewWriterLevel(buf, gzip.DefaultCompression)
			if err != nil {
				log.Error("gzipNewWriterLevel").Msgf("%v", err)
				compressor = nil
			}
		}
		if compressor != nil {
			_, errW := compressor.Write(bs)
			if errW != nil {
				log.Error("compressorWrite").Msgf("%v", err)
			}
			errClose := compressor.Close()
			if errClose != nil {
				log.Error(" compressorClose").Msgf("%v", err)
			}

			if errW == nil && errClose == nil {
				c.Header("Content-Encoding", encoding)
				bs = buf.Bytes()
			}
		}
	}

	c.Render(http.StatusOK, render.Data{
		ContentType: "application/json; charset=utf-8",
		Data:        bs,
	})
	c.Abort()
}

func GetCtxStoreRspStr(c *gin.Context) string {
	return c.GetString(respCtxKey)
}
