package tools

import (
	"errors"
	"flag"
	"log"
	"os"
	"path"
	"strings"

	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/env"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/file"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/verbose"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/model"
)

const (
	SOCKS5 = "socks5"
	HTTPS  = "https"
	HTTP   = "http"
)

func ParseM3U8(filename string) []*model.TSLink {
	res := make([]*model.TSLink, 0)
	lines := file.ReadLines(filename)
	if len(lines) == 0 {
		return res
	}
	for _, line := range lines {
		if strings.Contains(line, ".ts") {
			res = append(res, model.NewTSLink(line))
		}
	}
	return res
}

func CreateFolder(path string) error {
	exist, err := file.Exist(path)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	verbose.Log("path created at %s", path)
	return nil
}

func ValidateFlags() {
	if env.Platform == "" && env.TaskDefinitionFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	if env.Platform == "" {
		log.Fatal("[error] please specific a platform (asobistage / eplus / zaiko)")
	}
	if env.TaskDefinitionFile == "" {
		log.Fatal("[error] please specific a path to task config json file")
	}

	// validate proxy flags
	if env.ProxyType == "" {
		return // without proxy
	}
	if MatchProxy(env.ProxyType) == "" {
		log.Fatal("[error] we are not support the proxy type that you presented, you can only use socks5, https or http proxy")
	}
}

func ValidateArchive(distPath string) error {
	exist, _ := file.Exist(distPath)
	if !exist {
		return errors.New("please execute download script")
	}

	downloadedFiles, _ := file.ListDir(distPath)
	tsLinks := make([]*model.TSLink, 0)
	for _, downloadedFile := range downloadedFiles {
		if strings.HasSuffix(downloadedFile, ".m3u8") {
			tsLinks = ParseM3U8(path.Join(distPath, downloadedFile))
		}
	}

	missingAmount := 0
	for _, tsLink := range tsLinks {
		fullPath := path.Join(distPath, tsLink.Filename)
		exist0, _ := file.Exist(fullPath)
		if !exist0 {
			log.Printf("missing file: %s\n", tsLink.Filename)
			missingAmount += 1
		}
	}
	if missingAmount == 0 {
		log.Println("download completed")
	} else {
		log.Printf("missing %v files\n", missingAmount)
	}
	return nil
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
