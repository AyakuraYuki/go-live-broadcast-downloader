package internal

import (
	"errors"
	"go-live-broadcast-downloader/plugins/file"
	"log"
	"os"
	builtinPath "path"
	"strings"
)

func ParseM3U8(filename string) []*TSLink {
	res := make([]*TSLink, 0)
	lines := file.ReadLines(filename)
	if len(lines) == 0 {
		return res
	}
	for _, line := range lines {
		if strings.Contains(line, ".ts") {
			res = append(res, NewTSLink(line))
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
	log.Printf("path created at %s\n", path)
	return nil
}

func Validate(path string) error {
	exist, _ := file.IsPathExist(path)
	if !exist {
		return errors.New("please execute download script")
	}

	downloadedFiles := make([]string, 0)
	file.WalkDir(path, &downloadedFiles)
	tsLinks := make([]*TSLink, 0)
	for _, downloadedFile := range downloadedFiles {
		if strings.HasSuffix(downloadedFile, ".m3u8") {
			tsLinks = ParseM3U8(builtinPath.Join(path, downloadedFile))
		}
	}

	missingAmount := 0
	for _, tsLink := range tsLinks {
		fullPath := builtinPath.Join(path, tsLink.Filename)
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
