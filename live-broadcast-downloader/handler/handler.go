package handler

import (
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/model"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/platform"
)

const (
	MaxTick = 20000
)

type platformHandlerFunc func(*model.Task) error
type taskValidatorFunc func(*model.Task) error

var PlatformHandler = map[platform.Platform]platformHandlerFunc{
	platform.Asobistage: asobistage,
	platform.Eplus:      eplus,
	platform.Zaiko:      zaiko,
	platform.StreamPass: streampass,
}

var TaskValidator = map[platform.Platform]taskValidatorFunc{
	platform.Asobistage: asobistageTaskValidator,
	platform.Eplus:      eplusTaskValidator,
	platform.Zaiko:      zaikoTaskValidator,
	platform.StreamPass: streampassTaskValidator,
}
