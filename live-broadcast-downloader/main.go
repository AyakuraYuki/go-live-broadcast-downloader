package main

import (
	"flag"
	"fmt"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/file"
	cjson "github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/json"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/localization"
	"github.com/Xuanwo/go-locale"
	"golang.org/x/text/language"
	"log"
	"os"
	"strings"
	"time"
)

var (
	platform           string
	taskDefinitionFile string
	localeTag          language.Tag
	err                error
)

func init() {
	localeTag, err = locale.Detect()
	if err != nil {
		localeTag = language.English
	}
	l10nDictionary := localization.GetLocalizationDictionary(localeTag)

	flag.StringVar(&platform, "plat", platform, l10nDictionary[localization.KeyPlatform])
	flag.StringVar(&taskDefinitionFile, "json", taskDefinitionFile, l10nDictionary[localization.KeyTaskDefinitionFile])

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage of %s: %s -plat <asobistage|eplus|zaiko> -json </path/to/config.json>\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, l10nDictionary[localization.KeyUsage])
	}
}

func main() {
	flag.Parse()
	if platform == "" && taskDefinitionFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	if platform == "" {
		log.Fatal("[error] please specific a platform (asobistage / eplus / zaiko)")
	}
	if taskDefinitionFile == "" {
		log.Fatal("[error] please specific a path to task config json file")
	}

	platformHandler := internal.PlatformHandler[strings.ToLower(platform)]
	if platformHandler == nil {
		log.Fatalf("platform %s currently not supported\n", platform)
	}
	taskValidator := internal.TaskValidator[strings.ToLower(platform)]

	tasks := make([]*internal.Task, 0)

	jsonConfigContent := file.ReadFile(taskDefinitionFile)
	if jsonConfigContent == "" {
		log.Fatal("[error] empty config content, exit")
	}
	err = cjson.JSON.Unmarshal([]byte(jsonConfigContent), &tasks)
	if err != nil {
		panic(err)
	}

	st := time.Now()
	log.Println("This is a program that downloads live broadcast archives from asobistage, eplus, zaiko and other m3u8-base stream archives.")
	for _, task := range tasks {
		// validate
		if taskValidator != nil {
			err = taskValidator(task)
			if err != nil {
				panic(err)
			}
		}
		// create dist dir
		err = internal.CreateFolder(task.SaveTo)
		if err != nil {
			panic(err)
		}
		// download...
		for {
			err = platformHandler(task)
			if err != nil {
				log.Printf("Error: %v\n", err)
				log.Println("Error occurred, restarting...")
				time.Sleep(time.Second * 5)
				continue
			}
			break
		}
		// check downloaded resources
		err = internal.Validate(task.SaveTo)
		if err != nil {
			log.Printf("Validate failed: %v\n", err)
		} else {
			log.Printf("Task done with: %v\n", task.Prefix)
		}
	}

	et := time.Now()
	log.Printf("[%s] Done!\n", et.Sub(st))
}
