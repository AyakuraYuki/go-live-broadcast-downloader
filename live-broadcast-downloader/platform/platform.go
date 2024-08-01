package platform

import "strings"

type Platform string

func (platform Platform) String() string { return strings.ToLower(string(platform)) }

const (
	Asobistage Platform = "asobistage"
	Eplus      Platform = "eplus"
	Zaiko      Platform = "zaiko"
	StreamPass Platform = "streampass"

	Common Platform = "common" // common m3u8 sources
)
