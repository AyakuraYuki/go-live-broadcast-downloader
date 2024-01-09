package handler

import (
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/model"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/consts"
	nhttp "github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/net/http"
)

const (
	MaxTick = 20000
)

type platformHandlerFunc func(*model.Task, *nhttp.ProxyOption) error
type taskValidatorFunc func(*model.Task) error

var PlatformHandler = map[string]platformHandlerFunc{
	consts.Asobistage: asobistage,
	consts.Eplus:      eplus,
	consts.Zaiko:      zaiko,
	consts.StreamPass: streampass,
}

var TaskValidator = map[string]taskValidatorFunc{
	consts.Asobistage: asobistageTaskValidator,
	consts.Eplus:      eplusTaskValidator,
	consts.Zaiko:      zaikoTaskValidator,
	consts.StreamPass: streampassTaskValidator,
}
