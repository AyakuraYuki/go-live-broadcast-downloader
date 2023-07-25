package osstool

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cavaliergopher/grab/v3"
	ay_const "go-live-broadcast-downloader/plugins/ay-const"
	ali_oss "go-live-broadcast-downloader/plugins/conf/ali-oss"
	"go-live-broadcast-downloader/plugins/crypto"
	"go-live-broadcast-downloader/plugins/log"
	"go-live-broadcast-downloader/plugins/misc"
	nhttp "go-live-broadcast-downloader/plugins/net/http"
	"go-live-broadcast-downloader/plugins/sequence"
	"go-live-broadcast-downloader/plugins/typeconvert"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func init() {
	grab.DefaultClient = &grab.Client{
		HTTPClient: &http.Client{
			Timeout: time.Second * 3,
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment,
				MaxIdleConnsPerHost: 16,
				IdleConnTimeout:     time.Minute,
			},
		},
	}
}

// DownloadImgs 下载图片到dodo 服务器[批量]
func DownloadImgs(imgUrls []string, ossConf ali_oss.OssConfig) (map[string]string, error) {
	var (
		mp     = make(map[string]string)
		locker sync.Mutex
	)

	funs := make([]misc.WorkFunc, 0)
	for _, imgUrl := range imgUrls {
		copyImgUrl := imgUrl // 不要删除，会引起问题
		funs = append(funs, func() error {
			newImgUrl, err := DownloadImg(copyImgUrl, ossConf, nil)
			locker.Lock()
			mp[copyImgUrl] = newImgUrl
			locker.Unlock()
			return err
		})
	}
	err := misc.MultiRun(funs...)
	return mp, err
}

// DownloadImgsWithFilterHost 下载图片到dodo 服务器[批量]，filterHost域名下的不下载
func DownloadImgsWithFilterHost(imgUrls map[string]struct{}, filterHosts map[string]struct{}, ossConf ali_oss.OssConfig) (map[string]string, error) {
	var (
		mp     = make(map[string]string)
		locker sync.Mutex
	)

	funs := make([]misc.WorkFunc, 0)
	for imgUrl := range imgUrls {
		copyImgUrl := imgUrl // 不要删除，会引起问题
		funs = append(funs, func() error {
			newImgUrl, err := DownloadImg(copyImgUrl, ossConf, filterHosts)
			locker.Lock()
			mp[copyImgUrl] = newImgUrl
			locker.Unlock()
			return err
		})
	}
	err := misc.MultiRun(funs...)
	return mp, err
}

// DownloadImg 下载图片到dodo 服务器[单个]
func DownloadImg(imgUrl string, ossConf ali_oss.OssConfig, filterHosts map[string]struct{}) (string, error) {
	u, err := url.Parse(imgUrl)
	if err != nil {
		log.Error("DownloadImg").Msgf("url:%s, err:%v", imgUrl, err)
		return "", errors.New(fmt.Sprintf("无法获取图片信息 [%s]", imgUrl))
	}
	// 获取文件大小，head 头
	header, httpStatus, err := nhttp.Head(nil, imgUrl, nil, 3000, 1)
	if err != nil {
		log.Error("DownloadImg").Msgf("url:%s, err:%v", imgUrl, err)
		return "", errors.New(fmt.Sprintf("无法获取图片信息 [%s]", imgUrl))
	}
	if httpStatus != http.StatusOK {
		log.Error("DownloadImg").Msgf("url:%s, httpStatus:%d", imgUrl, httpStatus)
		return "", errors.New(fmt.Sprintf("无法获取图片信息 [%s]", imgUrl))
	}

	// 判断是否需要过滤域名
	if len(filterHosts) > 0 {
		if _, hit := filterHosts[fmt.Sprintf("%s://%s", u.Scheme, u.Host)]; hit {
			return imgUrl, nil
		}
	}

	// 判断长度
	if header == nil || header.Get("Content-Length") == "" {
		log.Error("DownloadImg").Msgf("url:%s, no Content-Length", imgUrl)
		return "", errors.New(fmt.Sprintf("无法获取图片信息 [%s]", imgUrl))
	}
	// 20M 限制
	if size := typeconvert.StringToInt64(header.Get("Content-Length")); size > 1024*1024*20 || size <= 0 {
		log.Error("DownloadImg").Msgf("url:%s, Content-Length --> %d", imgUrl, size)
		return "", errors.New(fmt.Sprintf("图片超过20M [%s]", imgUrl))
	}

	// 下载文件
	localPath := fmt.Sprintf("/tmp/%d", sequence.ID())
	_, err = grab.Get(localPath, imgUrl)
	if err != nil {
		log.Error("DownloadImg").Msgf("url:%s, err:%v", imgUrl, err)
		return "", errors.New(fmt.Sprintf("图片下载失败 [%s]", imgUrl))
	}
	defer os.RemoveAll(localPath)
	tmpOssPath := fmt.Sprintf("tmp/download/out/%d", sequence.ID())
	_, err = OssUpload(localPath, tmpOssPath, ossConf)
	if err != nil {
		log.Error("DownloadImg").Msgf("url:%s, err:%v", imgUrl, err)
		return "", errors.New(fmt.Sprintf("图片下载失败 [%s]", imgUrl))
	}

	tmpOssFullPath := fmt.Sprintf("https://%s.%s/%s", ossConf.OssBucket, ossConf.OssEndpoint, tmpOssPath)
	// 获取原图的宽高
	imgInfo, err := GetImgMetaInfo(tmpOssFullPath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("无法获取图片信息 [%s]", imgUrl))
	}

	// 扩展名
	ext := ""
	format := strings.ToLower(imgInfo.Format.Value)
	if strings.Contains(format, "jpg") {
		ext = "jpg"
	} else if strings.Contains(format, "png") {
		ext = "png"
	} else if strings.Contains(format, "gif") {
		ext = "gif"
	} else if strings.Contains(format, "webp") {
		ext = "webp"
	} else if strings.Contains(format, "jpeg") {
		ext = "jpeg"
	} else if strings.Contains(format, "bmp") {
		ext = "bmp"
	} else {
		ext = format
	}

	if ext == "" {
		return "", errors.New(fmt.Sprintf("无法获取图片信息 [%s]", imgUrl))
	}

	// 移动零时文件目录
	newOssPath := fmt.Sprintf("dodo/%s.%s", crypto.Md5Str(tmpOssPath), ext)
	options := []oss.Option{
		oss.ContentType(GetImageContentType(ext)),
	}
	if err = OssMoveFile(tmpOssPath, newOssPath, ossConf, options...); err != nil {
		return "", errors.New(fmt.Sprintf("图片下载失败 [%s]", imgUrl))
	}
	return fmt.Sprintf("%s/%s", ay_const.CDNHost, newOssPath), nil
}

func GetImageContentType(format string) string {
	mp := map[string]string{
		"jpg":  "image/jpeg",
		"tiff": "image/tiff",
		"gif":  "image/gif",
		"jfif": "image/jpeg",
		"png":  "image/png",
		"tif":  "image/tiff",
		"ico":  "image/x-icon",
		"jpeg": "image/jpeg",
		"wbmp": "image/vnd.wap.wbmp",
		"webp": "image/webp",
		"fax":  "image/fax",
		"net":  "image/pnetvue",
		"jpe":  "image/jpeg",
		"rp":   "image/vnd.rn-realpix",
	}
	return mp[strings.ToLower(format)]

}
