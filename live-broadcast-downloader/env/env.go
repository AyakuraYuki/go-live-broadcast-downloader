package env

import "golang.org/x/text/language"

var (
	Platform           string
	TaskDefinitionFile string
	ProxyHost          string
	ProxyPort          int
	ProxyType          string
	Coroutines         int
	LocaleTag          language.Tag
	MaxRetry           int
	Verbose            bool
)
