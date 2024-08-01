package downloader

import (
	"errors"
	"fmt"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/env"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/model"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/tools"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/file"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/misc"
	nhttp "github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/net/http"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/verbose"
	"github.com/samber/lo"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

func DownloadFile(downloadUrl, saveTo, filename string, proxy *nhttp.ProxyOption) error {
	var client *http.Client
	if proxy != nil {
		p, err := url.Parse(proxy.ProxyRawUrl())
		if err != nil {
			return err
		}
		client = &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 2048,
				IdleConnTimeout:     time.Minute * 5,
				Proxy:               http.ProxyURL(p),
			},
		}
	}
	bs, _, _, err := nhttp.GetRaw(client, downloadUrl, nil, nil, 60*1000, 2)
	if err != nil {
		return err
	}
	fullPath := path.Join(saveTo, filename)
	return os.WriteFile(fullPath, bs, os.ModePerm)
}

func Process(task *model.Task, proxy *nhttp.ProxyOption) error {
	m3u8Path := path.Join(task.SaveTo, task.Spec.Filename)
	tsLinks := tools.ParseM3U8(m3u8Path)
	if len(tsLinks) == 0 {
		return errors.New("empty playlist")
	}

	taskAmount := len(tsLinks)
	partitionSize := taskAmount / env.Coroutines
	counterChan := make(chan bool)
	defer close(counterChan)
	counter := uint64(0)
	go func() {
		for range counterChan {
			counter += 1
			if verbose.Verbose {
				fmt.Printf("downloading [%v / %v] \r", counter, taskAmount)
			}
		}
	}()

	funcs := make([]misc.WorkFunc, 0)
	var partitions [][]*model.TSLink
	if partitionSize == 0 {
		partitions = append(partitions, tsLinks)
	} else {
		partitions = lo.Chunk(tsLinks, partitionSize)
	}
	for _, partition := range partitions {
		partition := partition
		funcs = append(funcs, func() error {
			var err0 error
			for _, link := range partition {
				link := link
				tsPath := path.Join(task.SaveTo, link.Filename)
				exist, _ := file.IsPathExist(tsPath)
				if exist {
					verbose.Printf("skipped exist file: %s", link.Filename)
					counterChan <- true
					continue
				}
				tsUrl := fmt.Sprintf("%s/%s", strings.TrimRight(task.Prefix, "/"), link.Filename)
				if err0 = DownloadFile(tsUrl, task.SaveTo, link.Filename, proxy); err0 == nil {
					counterChan <- true
				}
			}
			return err0
		})
	}
	if err := misc.MultiRun(funcs...); err != nil {
		return err
	}
	return nil
}
