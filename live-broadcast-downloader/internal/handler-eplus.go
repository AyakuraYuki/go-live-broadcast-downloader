package internal

import (
	"errors"
	"log"
	"os/exec"
	"path"
)

func eplusTaskValidator(task *Task) error {
	if task.Prefix == "" {
		return errors.New("missing prefix")
	}
	if task.SaveTo == "" {
		return errors.New("we don't know where you want to save the archive")
	}
	if task.PageUrl == "" {
		return errors.New("eplus requires the page url")
	}
	if task.Spec == nil {
		return errors.New("missing spec")
	}
	if task.Spec.Filename == "" {
		return errors.New("missing m3u8 filename")
	}
	return nil
}

func eplus(task *Task) error {
	if err := DownloadFile(task.M3U8Url(), task.SaveTo, task.Spec.Filename); err != nil {
		return err
	} else {
		log.Printf("[eplus] download m3u8 playlist successed, file: %s\n", path.Join(task.SaveTo, task.Spec.Filename))
	}

	if err := Process(task); err != nil {
		return err
	} else {
		log.Printf("[eplus] successfully download all files in playlist\n")
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
