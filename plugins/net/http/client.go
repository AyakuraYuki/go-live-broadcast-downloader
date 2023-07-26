package nhttp

import (
	"bytes"
	"context"
	"fmt"
	cjson "github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	// NoKeepAliveClient is http client without keep alive
	NoKeepAliveClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	defaultHTTPClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 2048,
			IdleConnTimeout:     time.Minute * 5,
		},
	}
	defaultTimeout    = 500
	defaultRetryCount = 2
)

// ========================================================================================================================
// Proxy

const (
	SOCKS5 = "socks5"
	HTTPS  = "https"
	HTTP   = "http"
)

type ProxyOption struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	ProxyType string `json:"proxyType"`
}

func (t *ProxyOption) ProxyRawUrl() string {
	u := url.URL{
		Scheme: t.ProxyType,
		Host:   fmt.Sprintf("%s:%d", t.Host, t.Port),
	}
	return u.String()
}

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

// ========================================================================================================================

func assembleRequestParams(params ...int) (timeout int, retry int) {
	timeout, retry = defaultTimeout, defaultRetryCount
	switch {
	case len(params) >= 2:
		timeout, retry = params[0], params[1]
	case len(params) >= 1:
		timeout = params[0]
	}
	return timeout, retry
}

// PostRaw do http post, returns bytes data, response headers, http code and function error
// params: []int{timeoutInMillis, retryTimes}
func PostRaw(client *http.Client, requestUrl string, header http.Header, requestBody interface{}, params ...int) (data []byte, rspHeader http.Header, httpCode int, err error) {
	timeout, retry := assembleRequestParams(params...)
	for i := 0; i < retry; i++ {
		data, rspHeader, httpCode, err = do(client, http.MethodPost, requestUrl, header, requestBody, timeout)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Printf("PostRaw err: %s\n", err)
	}
	return data, rspHeader, httpCode, err
}

// PostWithUnmarshal do http post with unmarshal
func PostWithUnmarshal(client *http.Client, requestUrl string, header http.Header, requestBody interface{}, resp interface{}, params ...int) error {
	data, _, _, err := PostRaw(client, requestUrl, header, requestBody, params...)
	if err != nil {
		return err
	}
	// for no resp needed request.
	if resp == nil {
		return nil
	}
	// for big int
	decoder := cjson.JSON.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()
	err = decoder.Decode(resp)
	if err != nil {
		log.Printf("PostWithUnmarshal err: %s, requestUrl: %s, respData: %s\n", err, requestUrl, string(data))
	}
	return err
}

// GetRaw get http raw
func GetRaw(client *http.Client, requestUrl string, header http.Header, requestBody interface{}, params ...int) (data []byte, rspHeader http.Header, httpCode int, err error) {
	timeout, retry := assembleRequestParams(params...)
	for i := 0; i < retry; i++ {
		data, rspHeader, httpCode, err = do(client, http.MethodGet, requestUrl, header, requestBody, timeout)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Printf("GetRaw err: %s\n", err)
	}
	return data, rspHeader, httpCode, err
}

// GetWithUnmarshal do http get with unmarshal
func GetWithUnmarshal(client *http.Client, requestUrl string, header http.Header, requestBody interface{}, resp interface{}, params ...int) error {
	data, _, _, err := GetRaw(client, requestUrl, header, requestBody, params...)
	if err != nil {
		return err
	}
	// for no resp needed request.
	if resp == nil {
		return nil
	}
	// for big int
	decoder := cjson.JSON.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()
	err = decoder.Decode(resp)
	if err != nil {
		log.Printf("GetWithUnmarshal err: %s, requestUrl: %s, respData: %s\n", err, requestUrl, string(data))
	}
	return err
}

// Head .
func Head(client *http.Client, requestUrl string, header http.Header, params ...int) (rspHeader http.Header, httpStatus int, err error) {
	timeout, retry := assembleRequestParams(params...)
	for i := 0; i < retry; i++ {
		_, rspHeader, httpStatus, err = do(client, http.MethodHead, requestUrl, header, nil, timeout)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Printf("Head err: %s\n", err)
	}
	return rspHeader, httpStatus, err
}

// do request and returns bytes data, response headers, http code and function error
func do(client *http.Client, method string, requestUrl string, header http.Header, requestBody interface{}, timeout int) ([]byte, http.Header, int, error) {
	if client == nil {
		client = defaultHTTPClient
	}

	var reader io.Reader
	switch v := requestBody.(type) {
	case url.Values:
		reader = strings.NewReader(v.Encode())
	case []byte:
		reader = bytes.NewBuffer(v)
	case string:
		reader = strings.NewReader(v)
	case io.Reader:
		reader = v
	default:
		buff := &bytes.Buffer{}
		err := cjson.JSON.NewEncoder(buff).Encode(v)
		if err != nil {
			return nil, nil, -1, err
		}
		reader = buff
	}

	if (method == http.MethodGet || method == http.MethodHead) && requestBody == nil {
		reader = nil
	}

	req, err := http.NewRequest(method, requestUrl, reader)
	if err != nil {
		return nil, nil, -1, err
	}

	if header != nil {
		req.Header = header
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeout))
	defer cancelFunc()
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, -1, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, -1, err
	}

	return data, resp.Header, resp.StatusCode, nil
}
