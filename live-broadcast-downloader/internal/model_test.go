package internal

import (
	"go-live-broadcast-downloader/plugins/consts"
	"testing"
)

var task = &Task{
	Prefix:  "https://example.com/m3u8/content-114514",
	SaveTo:  "/path/to/user-home",
	PageUrl: "https://example.com",
	Spec: &M3U8Spec{
		PlaylistFilename: "index_6m.m3u8",
		KeyName:          "aes128.key",
	},
}

func TestTask_KeyUrl(t *testing.T) {
	t.Log(task.KeyUrl())
}

func TestTask_M3U8Url(t *testing.T) {
	t.Log(task.M3U8Url())
}

func TestNewTSLink(t *testing.T) {
	links := []string{
		"https://www.example.com/video/playlist/content-114514/index_6m_00001.ts?time=1564987332",
		"https://www.example.com/video/playlist/content-114514/index_6m_00002.ts",
		"www.example.com/video/playlist/content-114514/index_6m_00003.ts",
		"index_6m_00004.ts",
	}
	for _, link := range links {
		tsLink := NewTSLink(link)
		t.Logf("%#v", tsLink)
	}
}

func TestProxyOption_GetProxyServer(t *testing.T) {
	p := &ProxyOption{
		Host:      "127.0.0.1",
		Port:      7890,
		ProxyType: consts.SOCKS5,
	}
	t.Log(p.GetProxyServer())
}
