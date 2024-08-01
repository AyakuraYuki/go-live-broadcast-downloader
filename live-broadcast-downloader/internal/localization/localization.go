package localization

import (
	"golang.org/x/text/language"
)

var usageL10N map[language.Tag]map[string]string

var (
	supportedTags = []language.Tag{
		language.English,
		language.Chinese,
		language.SimplifiedChinese,
	}
	matcher = language.NewMatcher(supportedTags)
)

const (
	KeyUsage              = "Usage"
	KeyPlatform           = "Platform"
	KeyTaskDefinitionFile = "TaskDefinitionFile"
	KeyProxyHost          = "ProxyHost"
	KeyProxyPort          = "ProxyPort"
	KeyProxyType          = "ProxyType"
	KeyCoroutines         = "Coroutines"
	KeyMaxRetry           = "MaxRetry"
	KeyVerbose            = "Verbose"
)

func init() {
	usageL10N = make(map[language.Tag]map[string]string)
	usageL10N[language.Chinese] = chinese()
	usageL10N[language.English] = english()
}

var registeredTag = language.English

func RegisterTag(tag language.Tag) {
	registeredTag = tag
}

func GetLocalizationDictionary() map[string]string {
	_, index, c := matcher.Match(registeredTag)
	if index > len(supportedTags) || c < language.Low {
		return usageL10N[language.English]
	}

	matchedTag := supportedTags[index]
	switch matchedTag {
	case language.English:
		return usageL10N[language.English]
	case language.Chinese, language.SimplifiedChinese:
		return usageL10N[language.Chinese]
	}

	return usageL10N[language.English]
}

func GetOrDefault(key, defaultString string) string {
	dict := GetLocalizationDictionary()
	text, ok := dict[key]
	if !ok {
		return defaultString
	}
	return text
}

func Get(key string) string {
	return GetOrDefault(key, "")
}
