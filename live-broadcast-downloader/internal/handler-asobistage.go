package internal

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go-live-broadcast-downloader/plugins/file"
	cjson "go-live-broadcast-downloader/plugins/json"
	"go-live-broadcast-downloader/plugins/typeconvert"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
)

func asobistage(task *Task) error {
	if err := DownloadFile(task.KeyUrl(), task.SaveTo, task.Spec.KeyName); err != nil {
		return err
	} else {
		log.Printf("[asobistage] download key file successed, file: %s\n", path.Join(task.SaveTo, task.Spec.KeyName))
	}

	if err := DownloadFile(task.M3U8Url(), task.SaveTo, task.Spec.PlaylistFilename); err != nil {
		return err
	} else {
		log.Printf("[asobistage] download m3u8 playlist successed, file: %s\n", path.Join(task.SaveTo, task.Spec.PlaylistFilename))
	}

	if err := Process(task); err != nil {
		return err
	} else {
		log.Printf("[asobistage] successfully download all files in playlist\n")
	}

	// merge clips
	m3u8Path := path.Join(task.SaveTo, task.Spec.PlaylistFilename)
	videoPath := path.Join(task.SaveTo, "output.mp4")
	cmd := exec.Command("ffmpeg", "-allowed_extensions", "ALL", "-i", m3u8Path, "-c", "copy", videoPath)
	if err := cmd.Run(); err != nil {
		return err
	}

	// download comments
	asobistageComments(task)

	return nil
}

func asobistageComments(task *Task) {
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
	exist, _ := file.IsPathExist(saveTo)
	if exist {
		_ = os.Remove(saveTo)
	}
	outputFile, err := os.Create(saveTo)
	if err != nil {
		log.Printf("[asobistage] [comments] error when creating comments file, breaked. err: %v\n", err)
		return
	}
	defer outputFile.Close()

	wssUrl := fmt.Sprintf("wss://replay.asobistore.jp/%s_%s_ch1/archive", event, day)
	_, _ = outputFile.WriteString("[")

	// fetch comments from websocket
	conn, _, err := websocket.DefaultDialer.Dial(wssUrl, nil)
	if err != nil {
		log.Printf("[asobistage] [comments] error when connecting to comment server, breaked. err: %v\n", err)
		return
	}
	defer conn.Close()
	_, _, _ = conn.ReadMessage()
	noneCount := 0
	for tick := 0; tick < MaxTick; tick++ {
		bs, _ := cjson.JSON.Marshal(map[string]string{"func": "archive-get", "time": typeconvert.IntToString(5 * tick)})
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
