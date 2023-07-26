package handler

import (
	"errors"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/downloader"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/model"
	nhttp "github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/net/http"
	"log"
	"os/exec"
	"path"
)

func zaikoTaskValidator(task *model.Task) error {
	if task.Prefix == "" {
		return errors.New("missing prefix")
	}
	if task.SaveTo == "" {
		return errors.New("we don't know where you want to save the archive")
	}
	if task.Spec == nil {
		return errors.New("missing spec")
	}
	if task.Spec.Filename == "" {
		return errors.New("missing m3u8 filename")
	}
	return nil
}

func zaiko(task *model.Task, proxy *nhttp.ProxyOption) error {
	if err := downloader.DownloadFile(task.M3U8Url(), task.SaveTo, task.Spec.Filename, proxy); err != nil {
		return err
	} else {
		log.Printf("[zaiko] download m3u8 playlist successed, file: %s\n", path.Join(task.SaveTo, task.Spec.Filename))
	}

	if err := downloader.Process(task, proxy); err != nil {
		return err
	} else {
		log.Printf("[zaiko] successfully download all files in playlist\n")
	}

	// merge clips
	m3u8Path := path.Join(task.SaveTo, task.Spec.Filename)
	videoPath := path.Join(task.SaveTo, "output.mp4")
	cmd := exec.Command("ffmpeg", "-i", m3u8Path, "-c", "copy", videoPath)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
