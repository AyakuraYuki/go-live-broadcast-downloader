package tools

import (
	"errors"
	"flag"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/env"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/model"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/file"
	nhttp "github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/net/http"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/verbose"
	"log"
	"os"
	"path"
	"strings"
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
	exist, err := file.IsPathExist(path)
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
	verbose.Printf("path created at %s\n", path)
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
	if nhttp.MatchProxy(env.ProxyType) == "" {
		log.Fatal("[error] we are not support the proxy type that you presented, you can only use socks5, https or http proxy")
	}
}

func ValidateArchive(distPath string) error {
	exist, _ := file.IsPathExist(distPath)
	if !exist {
		return errors.New("please execute download script")
	}

	downloadedFiles := make([]string, 0)
	file.WalkDir(distPath, &downloadedFiles)
	tsLinks := make([]*model.TSLink, 0)
	for _, downloadedFile := range downloadedFiles {
		if strings.HasSuffix(downloadedFile, ".m3u8") {
			tsLinks = ParseM3U8(path.Join(distPath, downloadedFile))
		}
	}

	missingAmount := 0
	for _, tsLink := range tsLinks {
		fullPath := path.Join(distPath, tsLink.Filename)
		exist0, _ := file.IsPathExist(fullPath)
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
