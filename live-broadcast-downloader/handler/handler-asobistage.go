package handler

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"

	"github.com/gorilla/websocket"
	"github.com/spf13/cast"

	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/downloader"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/encoding/json"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/file"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/verbose"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/model"
)

func asobistageTaskExample() string {
	task := &model.Task{
		Prefix: "https://xxx.cloudfront.net/path/to/stream/",
		SaveTo: "/path/to/save",
		Spec: &model.M3U8Spec{
			Filename: "index_6m.m3u8",
			KeyName:  "aes128.key",
		},
	}
	return fmt.Sprintf("\n( example: \n%s \n)", json.Prettify(task))
}

func asobistageTaskValidator(task *model.Task) error {
	if task.Prefix == "" {
		return errors.New("missing prefix" + asobistageTaskExample())
	}
	if task.SaveTo == "" {
		return errors.New("we don't know where you want to save the archive" + asobistageTaskExample())
	}
	if task.Spec == nil {
		return errors.New("missing spec" + asobistageTaskExample())
	}
	if task.Spec.KeyName == "" {
		return errors.New("asobistage requires a key file to handle crypto *.ts data" + asobistageTaskExample())
	}
	if task.Spec.Filename == "" {
		return errors.New("missing m3u8 filename" + asobistageTaskExample())
	}
	return nil
}

func asobistage(task *model.Task) error {
	if err := downloader.DownloadFile(task.KeyUrl(), task.SaveTo, task.Spec.KeyName); err != nil {
		return err
	} else {
		log.Printf("[asobistage] download key file successed, file: %s\n", path.Join(task.SaveTo, task.Spec.KeyName))
	}

	if err := downloader.DownloadFile(task.M3U8Url(), task.SaveTo, task.Spec.Filename); err != nil {
		return err
	} else {
		log.Printf("[asobistage] download m3u8 playlist successed, file: %s\n", path.Join(task.SaveTo, task.Spec.Filename))
	}

	if err := downloader.Process(task); err != nil {
		return err
	} else {
		log.Printf("[asobistage] successfully download all files in playlist\n")
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

	// download comments
	asobistageComments(task)

	return nil
}

func asobistageComments(task *model.Task) {
	if task.PageUrl == "" {
		return
	}
	re := regexp.MustCompile("https://asobistage.asobistore.jp/event/([^/]+)/archive/([^/]+)")
	groups := re.FindStringSubmatch(task.PageUrl)
	event := groups[1]
	day := groups[2]
	if event == "" || day == "" {
		return
	}

	saveTo := path.Join(task.SaveTo, fmt.Sprintf("%s_%s_comments.json", event, day))
	exist, _ := file.Exist(saveTo)
	if exist {
		_ = os.Remove(saveTo)
	}
	outputFile, err := os.Create(saveTo)
	if err != nil {
		log.Printf("[asobistage] [comments] error when creating comments file, breaked. err: %v\n", err)
		return
	}
	defer func(f *os.File) { _ = f.Close() }(outputFile)

	wssUrl := fmt.Sprintf("wss://replay.asobistore.jp/%s_%s_ch1/archive", event, day)
	_, _ = outputFile.WriteString("[")

	// fetch comments from websocket
	conn, _, err := websocket.DefaultDialer.Dial(wssUrl, nil)
	if err != nil {
		log.Printf("[asobistage] [comments] error when connecting to comment server, breaked. err: %v\n", err)
		return
	}
	defer func(conn *websocket.Conn) { _ = conn.Close() }(conn)

	_, _, _ = conn.ReadMessage()
	noneCount := 0
	for tick := 0; tick < MaxTick; tick++ {
		bs, _ := json.JSON.Marshal(map[string]string{"func": "archive-get", "time": cast.ToString(5 * tick)})
		if err0 := conn.WriteMessage(websocket.TextMessage, bs); err0 != nil {
			log.Printf("[asobistage] [comments] error when sending message to comment server, breaked. err: %v\n", err)
			return
		}
		_, p, _ := conn.ReadMessage()
		rsp := string(p)
		comments := rsp[12 : len(rsp)-2]
		if tick > 0 && len(comments) != 0 {
			_, _ = outputFile.WriteString(",")
		}
		_, _ = outputFile.WriteString(comments)

		if len(comments) == 0 {
			noneCount += 1
		} else {
			noneCount = 0
		}

		if noneCount > 19 {
			break
		}
		fmt.Printf("Downloading comments... Tick: %v, Sending: %s, Empty: %v\r", tick, string(bs), noneCount)
	}
	fmt.Println()

	_, _ = outputFile.WriteString("]")
	return
}
