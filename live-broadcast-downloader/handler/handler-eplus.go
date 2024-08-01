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

func eplusTaskExample() string {
	task := &model.Task{
		Prefix:  "https://eplus.stream.host/path/to/stream/",
		SaveTo:  "/path/to/save",
		PageUrl: "https://eplus.example.com/page/url/of/stream",
		Spec: &model.M3U8Spec{
			Filename: "index_6m.m3u8",
		},
	}
	return fmt.Sprintf("\n( example: \n%s \n)", json.Prettify(task))
}

func eplusTaskValidator(task *model.Task) error {
	if task.Prefix == "" {
		return errors.New("missing prefix" + eplusTaskExample())
	}
	if task.SaveTo == "" {
		return errors.New("we don't know where you want to save the archive" + eplusTaskExample())
	}
	if task.PageUrl == "" {
		return errors.New("eplus requires the page url" + eplusTaskExample())
	}
	if task.Spec == nil {
		return errors.New("missing spec" + eplusTaskExample())
	}
	if task.Spec.Filename == "" {
		return errors.New("missing m3u8 filename" + eplusTaskExample())
	}
	return nil
}

func eplus(task *model.Task) error {
	if err := downloader.DownloadFile(task.M3U8Url(), task.SaveTo, task.Spec.Filename); err != nil {
		return err
	} else {
		log.Printf("[eplus] download m3u8 playlist successed, file: %s\n", path.Join(task.SaveTo, task.Spec.Filename))
	}

	if err := downloader.Process(task); err != nil {
		return err
	} else {
		log.Printf("[eplus] successfully download all files in playlist\n")
	}

	// merge clips
	m3u8Path := path.Join(task.SaveTo, task.Spec.Filename)
	videoPath := path.Join(task.SaveTo, "output.mp4")
	cmd := exec.Command("ffmpeg", "-y", "-i", m3u8Path, "-c", "copy", videoPath)
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
