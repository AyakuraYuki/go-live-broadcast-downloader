package internal

import (
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/consts"
	nhttp "github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/net/http"
)

const (
	MaxTick = 20000
)

type platformHandler func(*Task, *nhttp.ProxyOption) error
type taskValidator func(*Task) error

var PlatformHandler = map[string]platformHandler{
	consts.Asobistage: asobistage,
	consts.Eplus:      eplus,
	consts.Zaiko:      zaiko,
}

var TaskValidator = map[string]taskValidator{
	consts.Asobistage: asobistageTaskValidator,
	consts.Eplus:      eplusTaskValidator,
	consts.Zaiko:      zaikoTaskValidator,
}
