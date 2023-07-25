package main

import (
	"flag"
	"go-live-broadcast-downloader/live-broadcast-downloader/internal"
	"go-live-broadcast-downloader/plugins/file"
	cjson "go-live-broadcast-downloader/plugins/json"
	"log"
	"time"
)

var (
	platform           string
	taskDefinitionFile string
)

func asobistage(task *internal.Task) error {
	if err := internal.DownloadFile(task.KeyUrl(), task.SaveTo, task.Spec.KeyName); err != nil {
		return err
	}
	if err := internal.DownloadFile(task.M3U8Url(), task.SaveTo, task.Spec.PlaylistFilename); err != nil {
		return err
	}
	return nil
}

func main() {
	log.Println("This is a program that downloads live broadcast archives from asobistage, eplus, zaiko and other m3u8-base stream archives.")

	flag.StringVar(&platform, "plat", platform, "Live Broadcast 平台（asobistage, eplus, zaiko）")
	flag.StringVar(&taskDefinitionFile, "json", taskDefinitionFile, "任务定义文件")

	if platform == "" {
		log.Fatal("[error] please specific a platform (asobistage / eplus / zaiko)")
	}
	if taskDefinitionFile == "" {
		log.Fatal("[error] please specific a path to task config json file")
	}

	var err error
	tasks := make([]*internal.Task, 0)

	jsonConfigContent := file.ReadFile(taskDefinitionFile)
	if jsonConfigContent == "" {
		log.Fatal("[error] empty config content, exit")
	}
	err = cjson.JSON.Unmarshal([]byte(jsonConfigContent), &tasks)
	if err != nil {
		panic(err)
	}

	for _, task := range tasks {
		err = internal.CreateFolder(task.SaveTo)
		if err != nil {
			panic(err)
		}
		for {
			err = asobistage(task)
			if err != nil {
				log.Printf("Error: %v\n", err)
				log.Println("Error occurred, restarting...")
				time.Sleep(time.Second * 5)
				continue
			}
			err = internal.Process(task)
			if err != nil {
				log.Printf("Error: %v\n", err)
				log.Println("Error occurred, restarting...")
				time.Sleep(time.Second * 5)
				continue
			}
			break
		}
		err = internal.Validate(task.SaveTo)
		if err != nil {
			log.Printf("Validate failed: %v\n", err)
		} else {
			log.Printf("Task done with: %v\n", task.Prefix)
		}
	}
}
