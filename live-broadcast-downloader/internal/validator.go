package internal

import (
	"errors"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/file"
	"log"
	"path"
	"strings"
)

func ValidateArchive(distPath string) error {
	exist, _ := file.IsPathExist(distPath)
	if !exist {
		return errors.New("please execute download script")
	}

	downloadedFiles := make([]string, 0)
	file.WalkDir(distPath, &downloadedFiles)
	tsLinks := make([]*TSLink, 0)
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
