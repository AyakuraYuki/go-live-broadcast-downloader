package localization

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
            "keyName": "aes128.key",
            "rawQuery": "<raw_query_string_from_m3u8_link>"
        }
    }
]

This is a JSON array that declares a bunch of tasks with:
1. m3u8 playlist link prefix(which means a link to m3u8 file but remove the m3u8 filename)
2. which local path you want to save the archive, and it should be an absolute path
3. the archive page URL link
4. which archive resolution spec that you want to download
5. an optional crypto key filename presents by the m3u8 playlist
6. some m3u8 playlist link requires auth token, if that token passed by query string, you should copy the raw query string to "spec.rawQuery"

Please prepare your own tasks config by using the format we declared.
`
	dict[KeyPlatform] = "The name of Live Broadcast Platform, available values are [asobistage, eplus, zaiko, streampass]."
	dict[KeyTaskDefinitionFile] = "An absolute path of your task declaration JSON file."
	dict[KeyProxyHost] = "Proxy server host or IP address."
	dict[KeyProxyPort] = "Proxy server port."
	dict[KeyProxyType] = "Proxy type, available types are [http, https, socks5]."
	dict[KeyCoroutines] = "Declare the number of download threads"
	dict[KeyMaxRetry] = "max retry times for each task"
	dict[KeyVerbose] = "print more information when running"
	return dict
}
