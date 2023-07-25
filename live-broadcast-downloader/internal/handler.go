package internal

import "go-live-broadcast-downloader/plugins/consts"

const (
	MaxTick = 20000
)

type platformHandler func(*Task) error
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
