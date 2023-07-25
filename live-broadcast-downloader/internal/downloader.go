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
	m3u8Path := path.Join(task.SaveTo, task.Spec.Filename)
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
