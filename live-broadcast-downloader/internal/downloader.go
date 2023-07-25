package internal

import (
	"errors"
	"fmt"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/file"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/misc"
	nhttp "github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/net/http"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/part"
	"log"
	"os"
	"path"
	"strings"
)

func DownloadFile(downloadUrl, saveTo, filename string) error {
	bs, err := nhttp.GetRaw(nil, downloadUrl, nil, nil, 60*1000, 2)
	if err != nil {
		return err
	}
	fullPath := path.Join(saveTo, filename)
	return os.WriteFile(fullPath, bs, os.ModePerm)
}

func Process(task *Task) error {
	m3u8Path := path.Join(task.SaveTo, task.Spec.PlaylistFilename)
	tsLinks := ParseM3U8(m3u8Path)
	if len(tsLinks) == 0 {
		return errors.New("empty playlist")
	}

	funcs := make([]misc.WorkFunc, 0)
	for indexRange := range part.Partition(len(tsLinks), 500) {
		bulkLinks := tsLinks[indexRange.Low:indexRange.High]
		funcs = append(funcs, func() error {
			var err0 error
			for _, link := range bulkLinks {
				tsPath := path.Join(task.SaveTo, link.Filename)
				exist, _ := file.IsPathExist(tsPath)
				if exist {
					log.Printf("skipped exist file: %s", link.Filename)
					continue
				}
				tsUrl := fmt.Sprintf("%s/%s", strings.TrimRight(task.Prefix, "/"), link.Filename)
				err0 = DownloadFile(tsUrl, task.SaveTo, link.Filename)
			}
			return err0
		})
	}
	if err := misc.MultiRun(funcs...); err != nil {
		return err
	}
	return nil
}

// misc...

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

func Validate(distPath string) error {
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
