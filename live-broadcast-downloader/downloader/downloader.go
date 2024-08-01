package downloader

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"

	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/env"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/file"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/misc"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/tools"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/verbose"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/model"
)

var client *resty.Client

func init() {
	client = resty.New()
	client.SetRetryCount(3)
	if env.ProxyType != "" {
		proxy := url.URL{
			Scheme: env.ProxyType,
			Host:   fmt.Sprintf("%s:%d", env.ProxyHost, env.ProxyPort),
		}
		client.SetProxy(proxy.String())
	}
}

func DownloadFile(downloadUrl, saveTo, filename string) error {
	dst := filepath.Join(saveTo, filename)
	_, err := client.R().SetOutput(dst).Get(downloadUrl)
	return err
}

func Process(task *model.Task) error {
	m3u8Path := path.Join(task.SaveTo, task.Spec.Filename)
	tsLinks := tools.ParseM3U8(m3u8Path)
	if len(tsLinks) == 0 {
		return errors.New("empty playlist")
	}

	taskAmount := len(tsLinks)
	size := taskAmount / env.Coroutines
	var partitions [][]*model.TSLink
	if size > 0 {
		partitions = lo.Chunk(tsLinks, size)
	} else {
		partitions = append(partitions, tsLinks)
	}

	funcs := make([]misc.WorkFunc, 0)
	var counter atomic.Int32
	for _, partition := range partitions {
		partition := partition
		funcs = append(funcs, func() error {
			for _, link := range partition {
				link := link
				if exist, _ := file.Exist(path.Join(task.SaveTo, link.Filename)); exist {
					verbose.Log("skipped exist file: %s", link.Filename)
					continue
				}
				tsUrl := fmt.Sprintf("%s/%s", strings.TrimRight(task.Prefix, "/"), link.Filename)
				if err0 := DownloadFile(tsUrl, task.SaveTo, link.Filename); err0 == nil {
					fmt.Printf("downloading [%v / %v] \r", counter.Add(1), taskAmount)
				}
			}
			return nil
		})
	}
	if err := misc.MultiRun(funcs...); err != nil {
		return err
	}
	return nil
}
