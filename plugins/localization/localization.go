package localization

import (
	"golang.org/x/text/language"
)

var usageL10N map[language.Tag]map[string]string

var (
	languageSupport = []language.Tag{
		language.English,
		language.Chinese,
		language.SimplifiedChinese,
	}
	languageMatcher = language.NewMatcher(languageSupport)
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

func chinese() map[string]string {
	dict := make(map[string]string)
	dict[KeyUsage] = `任务定义文件是一个 JSON 配置文件，它的内容应该遵循如下的规则：
[
    {
        "prefix": "https://0000000000000.cloudfront.net/lalabit/mc/00000000-0000-0000-0000-000000000000/",
        "saveTo": "/home/username/archive/dist-path",
        "pageUrl": "https://live-broadcast-platform.host/link/to/archive/page-url",
        "spec": {
            "filename": "index.m3u8",
            "keyName": "aes128.key"
        }
    }
]

这是一个 JSON 数组内容，它定义了一组包含如下内容的任务：
1. m3u8 文件链接地址的前缀，它是一个去除 m3u8 文件名称后的链接地址
2. 一个你想要保存归档视频的文件夹完整路径
3. 视频播放页面的地址
4. 归档视频的解析度文件，一般不同的 m3u8 文件代表了不同的解析度，应该填写 m3u8 的文件名
5. 一个可选的加密解密文件的文件名，你能在 m3u8 文件内容中找到他

请遵循上面的声明，准备好你自己任务配置文件。
`
	dict[KeyPlatform] = "Live Broadcast 平台名称（asobistage, eplus, zaiko）"
	dict[KeyTaskDefinitionFile] = "任务定义文件的绝对路径"
	dict[KeyProxyHost] = "代理服务器主机地址/IP"
	dict[KeyProxyPort] = "代理服务器端口"
	dict[KeyProxyType] = "代理类型（http, https, socks5）"
	dict[KeyCoroutines] = "下载任务组的数量"
	dict[KeyMaxRetry] = "每个任务最大重试次数"
	dict[KeyVerbose] = "是否输出详细信息"
	return dict
}

func english() map[string]string {
	dict := make(map[string]string)
	dict[KeyUsage] = `The JSON configuration should be like the following text:
[
    {
        "prefix": "https://0000000000000.cloudfront.net/lalabit/mc/00000000-0000-0000-0000-000000000000/",
        "saveTo": "/home/username/archive/dist-path",
        "pageUrl": "https://live-broadcast-platform.host/link/to/archive/page-url",
        "spec": {
            "filename": "index.m3u8",
            "keyName": "aes128.key"
        }
    }
]

This is a JSON array that declares a bunch of tasks with:
1. m3u8 playlist link prefix(which means a link to m3u8 file but remove the m3u8 filename)
2. which local path you want to save the archive, and it should be an absolute path
3. the archive page URL link
4. which archive resolution spec that you want to download
5. an optional crypto key filename presents by the m3u8 playlist

Please prepare your own tasks config by using the format we declared.
`
	dict[KeyPlatform] = "The name of Live Broadcast Platform, available values are [asobistage, eplus, zaiko]."
	dict[KeyTaskDefinitionFile] = "An absolute path of your task declaration JSON file."
	dict[KeyProxyHost] = "Proxy server host or IP address."
	dict[KeyProxyPort] = "Proxy server port."
	dict[KeyProxyType] = "Proxy type, available types are [http, https, socks5]."
	dict[KeyCoroutines] = "Declare the number of download threads"
	dict[KeyMaxRetry] = "max retry times for each task"
	dict[KeyVerbose] = "print more information when running"
	return dict
}

func GetLocalizationDictionary(tag language.Tag) map[string]string {
	tagRegion, _ := tag.Region()
	matchedTag, _, _ := languageMatcher.Match(tag)
	matchedRegion, _ := matchedTag.Region()
	chineseRegion, _ := language.Chinese.Region()
	if tagRegion == matchedRegion && matchedRegion == chineseRegion {
		return usageL10N[language.Chinese]
	}
	return usageL10N[language.English]
}
