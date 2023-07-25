package consts

const (
	SOCKS5 = "socks5"
	HTTPS  = "https"
	HTTP   = "http"
)

func MatchProxy(proxyType string) string {
	if proxyType == "" {
		return ""
	}
	types := []string{SOCKS5, HTTPS, HTTP}
	for _, v := range types {
		if proxyType == v {
			return v
		}
	}
	return ""
}
