# go-live-broadcast-downloader

This is a program that downloads live broadcast archives from most of m3u8-base stream archives.

It supports the following platform currently:

- [Asobistage](https://asobistage.asobistore.jp/)
- [Eplus](https://ib.eplus.jp/)
- [Zaiko](https://zaiko.io/)
- [StreamPass]

This program does not use webdriver or any other headless browser. I made it download archive by using coroutines to boost the download speed.

## How to use

First thing first, you should have installed FFMpeg, we will merge the archive by using FFMpeg as soon as the archive clips downloaded.

> PS: I added some JSON configuration examples of different platforms. When your task configuration doesn't pass validation, you can get those examples.

```text
Usage of ./live-broadcast-downloader:
  -c string
        An absolute path of your task declaration JSON file.
  -config string
        An absolute path of your task declaration JSON file.
  -p string
        The name of Live Broadcast Platform, available values are [asobistage, eplus, zaiko, streampass].
  -plat string
        The name of Live Broadcast Platform, available values are [asobistage, eplus, zaiko, streampass].
  -proxy_host string
        Proxy server host or IP address. (default "127.0.0.1")
  -proxy_port int
        Proxy server port. (default 7890)
  -proxy_type string
        Proxy type, available types are [http, https, socks5].
  -retry int
        max retry times for each task (default 10)
  -threads int
        Declare the number of download threads (default 10)
  -verbose
        print more information when running

The JSON configuration should be like the following text:
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
```
