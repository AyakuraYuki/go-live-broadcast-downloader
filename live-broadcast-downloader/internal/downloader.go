package internal

import (
	"errors"
	"fmt"
	"go-live-broadcast-downloader/plugins/misc"
	nhttp "go-live-broadcast-downloader/plugins/net/http"
	"go-live-broadcast-downloader/plugins/part"
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
				tsUrl := fmt.Sprintf("%s/%s", strings.TrimRight(task.Prefix, "/"), link.Filename)
				err0 = DownloadFile(tsUrl, task.SaveTo, link.Filename)
			}
			return err0
		})
	}
	err := misc.MultiRun(funcs...)
	return err
}
