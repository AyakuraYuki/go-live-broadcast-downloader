package json

import (
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func RegisterFuzzyDecoders() {
	extra.RegisterFuzzyDecoders()
}

func Stringify(v any) string {
	raw, err := JSON.MarshalToString(v)
	if err != nil {
		return ""
	}
	return raw
}

func Prettify(v any) string {
	bs, err := JSON.MarshalIndent(v, "", "    ")
	if err != nil {
		return ""
	}
	return string(bs)
}
