package osstool

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	ali_oss "go-live-broadcast-downloader/plugins/conf/ali-oss"
	"go-live-broadcast-downloader/plugins/log"
	nhttp "go-live-broadcast-downloader/plugins/net/http"
	"strings"
)

// OssUpload oss 文件上传(参数 dst 不要以 / 开头)
func OssUpload(local, dst string, ossConf ali_oss.OssConfig, options ...oss.Option) (string, error) {
	cli, err := oss.New(ossConf.OssEndpoint, ossConf.OssAccessKeyID, ossConf.OssAccessKeySecret)
	if err != nil {
		log.Error("OssUpload").Msgf("%v", err)
		return "", err
	}
	bucket, err := cli.Bucket(ossConf.OssBucket)
	if err != nil {
		log.Error("OssUpload").Msgf("%v", err)
		return "", err
	}

	if strings.HasPrefix(dst, "/") {
		// 去掉前面的 /
		dst = strings.Replace(dst, "/", "", 1)
	}

	err = bucket.PutObjectFromFile(dst, local, options...)
	if err != nil {
		log.Error("OssUpload").Msgf("%v", err)
		return "", err
	}
	return dst, err
}

// GetImgMetaInfo 获取图片信息
func GetImgMetaInfo(ossUrlPath string) (*AliImageInfo, error) {
	infoUrl := fmt.Sprintf("%s?x-oss-process=image/info", ossUrlPath)
	imgInfo := &AliImageInfo{}
	err := nhttp.GetWithUnmarshal(nil, infoUrl, nil, nil, &imgInfo, 3000, 2)
	if err != nil {
		log.Error("GetImgMetaInfo").Msgf("%v", err)
		return imgInfo, err
	}
	return imgInfo, nil
}

// OssMoveFile oss 文件移动(同一个bucket)
func OssMoveFile(src, target string, ossConf ali_oss.OssConfig, options ...oss.Option) error {
	cli, err := oss.New(ossConf.OssEndpoint, ossConf.OssAccessKeyID, ossConf.OssAccessKeySecret)
	if err != nil {
		log.Error("OssMoveFile").Msgf("%v", err)
		return err
	}
	bucket, err := cli.Bucket(ossConf.OssBucket)
	if err != nil {
		log.Error("OssMoveFile").Msgf("%v", err)
		return err
	}
	// 去掉前面的 /
	if strings.HasPrefix(src, "/") {
		src = strings.Replace(src, "/", "", 1)
	}
	if strings.HasPrefix(target, "/") {
		target = strings.Replace(target, "/", "", 1)
	}
	if src == "" || target == "" {
		return nil
	}

	_, err = bucket.CopyObject(src, target)
	if err != nil {
		log.Error("OssMoveFile").Msgf("%v", err)
		return err
	}

	// 设置文件的meta
	_ = bucket.SetObjectMeta(target, options...)

	// 删除原文件
	err = bucket.DeleteObject(src)
	if err != nil {
		log.Error("OssMoveFile").Msgf("%v", err)
		return err
	}
	return nil
}
