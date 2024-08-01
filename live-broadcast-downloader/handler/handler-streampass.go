package handler

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path"

	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/downloader"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/encoding/json"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/verbose"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/model"
)

func streampassTaskExample() string {
	task := &model.Task{
		Prefix: "https://xxx.cloudfront.net/path/to/stream/",
		SaveTo: "/path/to/save",
		Spec: &model.M3U8Spec{
			Filename: "index_6m.m3u8",
			KeyName:  "aes128.key",
			RawQuery: "__token=xxxxxx",
		},
	}
	return fmt.Sprintf("\n( example: \n%s \n)", json.Prettify(task))
}

func streampassTaskValidator(task *model.Task) error {
	if task.Prefix == "" {
		return errors.New("missing prefix" + streampassTaskExample())
	}
	if task.SaveTo == "" {
		return errors.New("we don't know where you want to save the archive" + streampassTaskExample())
	}
	if task.Spec == nil {
		return errors.New("missing spec" + streampassTaskExample())
	}
	if task.Spec.KeyName == "" {
		return errors.New("streampass requires a key file to handle crypto *.ts data" + streampassTaskExample())
	}
	if task.Spec.Filename == "" {
		return errors.New("missing m3u8 filename" + streampassTaskExample())
	}
	if task.Spec.RawQuery == "" {
		return errors.New("missing m3u8 query string, StreamPass required token to get *.m3u8 file" + streampassTaskExample())
	}
	return nil
}

func streampass(task *model.Task) error {
	if err := downloader.DownloadFile(task.KeyUrl(), task.SaveTo, task.Spec.KeyName); err != nil {
		return err
	} else {
		log.Printf("[streampass] download key file successed, file: %s\n", path.Join(task.SaveTo, task.Spec.KeyName))
	}

	if err := downloader.DownloadFile(task.M3U8Url(), task.SaveTo, task.Spec.Filename); err != nil {
		return err
	} else {
		log.Printf("[streampass] download m3u8 playlist successed, file: %s\n", path.Join(task.SaveTo, task.Spec.Filename))
	}

	if err := downloader.Process(task); err != nil {
		return err
	} else {
		log.Printf("[streampass] successfully download all files in playlist\n")
	}

	// merge clips
	m3u8Path := path.Join(task.SaveTo, task.Spec.Filename)
	videoPath := path.Join(task.SaveTo, "output.mp4")
	cmd := exec.Command("ffmpeg", "-allowed_extensions", "ALL", "-y", "-i", m3u8Path, "-c", "copy", videoPath)
	var cmdOut, cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	if err := cmd.Run(); err != nil {
		verbose.Log(cmdErr.String())
		return err
	} else {
		verbose.Log(cmdOut.String())
		verbose.Log(cmdErr.String())
	}

	return nil
}
