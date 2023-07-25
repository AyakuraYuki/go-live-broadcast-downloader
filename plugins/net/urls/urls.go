package urls

import (
	"net/url"
)

// TryParseRawUrl returns scheme, host, path, query string
func TryParseRawUrl(rawUrl string) (scheme string, host string, path string, query string, err error) {
	u, err := url.Parse(rawUrl)
	if err != nil || u.Host == "" {
		u0, err0 := url.Parse("https://" + rawUrl)
		if err0 != nil {
			return "", "", "", "", err0
		}
		return "", u0.Host, u0.Path, u0.RawQuery, nil
	}
	return u.Scheme, u.Host, u.Path, u.RawQuery, nil
}

// ParseRawUrl parses url with both valid url and invalid url
func ParseRawUrl(rawUrl string) (*url.URL, error) {
	u, err := url.Parse(rawUrl)
	if err != nil || u.Host == "" {
		u0, err0 := url.Parse("https://" + rawUrl)
		if err0 != nil {
			return nil, err0
		}
		u0.Scheme = ""
		return u0, nil
	}
	return u, nil
}
