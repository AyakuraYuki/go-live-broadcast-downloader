package cjson

import (
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func RegisterFuzzyDecoders() {
	extra.RegisterFuzzyDecoders()
}
