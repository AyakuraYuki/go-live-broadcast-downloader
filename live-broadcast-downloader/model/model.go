package model

import (
	"log"
	"strings"

	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/urls"
)

// M3U8Spec declares the information of m3u8 playlist, presents m3u8 filename and crypto key filename
type M3U8Spec struct {
	Filename string `json:"filename"`
	KeyName  string `json:"keyName"`
	RawQuery string `json:"rawQuery"`
}

// Task declares the download task
type Task struct {
	Prefix  string    `json:"prefix"`
	SaveTo  string    `json:"saveTo"`
	PageUrl string    `json:"pageUrl"`
	Spec    *M3U8Spec `json:"spec"`
}

// KeyUrl assembles the download link of m3u8 crypto key, returns empty string if key is not defined
func (t *Task) KeyUrl() string {
	if t.Spec == nil || t.Spec.KeyName == "" {
		return ""
	}
	keyUrl := strings.Builder{}
	keyUrl.WriteString(t.Prefix)
	if !strings.HasSuffix(t.Prefix, "/") {
		keyUrl.WriteString("/")
	}
	keyUrl.WriteString(t.Spec.KeyName)
	return keyUrl.String()
}

// M3U8Url assembles the download link of m3u8 playlist file
func (t *Task) M3U8Url() string {
	m3u8Url := strings.Builder{}
	m3u8Url.WriteString(t.Prefix)
	if !strings.HasSuffix(t.Prefix, "/") {
		m3u8Url.WriteString("/")
	}
	m3u8Url.WriteString(t.Spec.Filename)
	if t.Spec.RawQuery != "" {
		m3u8Url.WriteString("?" + t.Spec.RawQuery)
	}
	return m3u8Url.String()
}

// TSLink declares the information of ts file in m3u8 playlist, presents link, base url, query string and the most important filename
type TSLink struct {
	Link        string `json:"link"`
	BaseUrl     string `json:"baseUrl"`
	QueryString string `json:"queryString"`
	Filename    string `json:"filename"`
}

// NewTSLink returns an object which contains the information of ts file
func NewTSLink(link string) *TSLink {
	t := &TSLink{
		Link:    link,
		BaseUrl: link,
	}
	t.parseLinkAttributes()
	t.parseFilename()
	return t
}

// parseLinkAttributes parses the link attributes when initializing
func (t *TSLink) parseLinkAttributes() {
	u, err := urls.ParseRawUrl(t.BaseUrl)
	if err != nil {
		log.Fatalf("[error] parse link attributes failed with resource: %#v", t)
		return
	}
	t.BaseUrl = strings.TrimLeft(strings.TrimRight(strings.ReplaceAll(u.String(), u.RawQuery, ""), "?"), "//")
	t.QueryString = u.RawQuery
}

// parseFilename parses the ts filename when initializing
func (t *TSLink) parseFilename() {
	if strings.Contains(t.BaseUrl, "/") {
		t.Filename = t.BaseUrl[strings.LastIndex(t.BaseUrl, "/")+1:]
	} else {
		t.Filename = t.BaseUrl
	}
}
