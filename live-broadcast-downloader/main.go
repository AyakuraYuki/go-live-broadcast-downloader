package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Xuanwo/go-locale"
	"golang.org/x/text/language"

	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/env"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/handler"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/encoding/json"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/file"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/localization"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/tools"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/verbose"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/model"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/platform"
)

func init() {
	var err error

	env.LocaleTag, err = locale.Detect()
	if err != nil {
		env.LocaleTag = language.English
	}
	localization.RegisterTag(env.LocaleTag)
	l10nDictionary := localization.GetLocalizationDictionary()

	flag.StringVar(&env.Platform, "p", "", l10nDictionary[localization.KeyPlatform])
	flag.StringVar(&env.Platform, "plat", "", l10nDictionary[localization.KeyPlatform])

	flag.StringVar(&env.TaskDefinitionFile, "c", "", l10nDictionary[localization.KeyTaskDefinitionFile])
	flag.StringVar(&env.TaskDefinitionFile, "config", "", l10nDictionary[localization.KeyTaskDefinitionFile])

	flag.StringVar(&env.ProxyHost, "proxy_host", "127.0.0.1", l10nDictionary[localization.KeyProxyHost])
	flag.IntVar(&env.ProxyPort, "proxy_port", 7890, l10nDictionary[localization.KeyProxyPort])
	flag.StringVar(&env.ProxyType, "proxy_type", "", l10nDictionary[localization.KeyProxyType])

	flag.IntVar(&env.Coroutines, "threads", 10, l10nDictionary[localization.KeyCoroutines])
	flag.IntVar(&env.MaxRetry, "retry", 10, l10nDictionary[localization.KeyMaxRetry])
	flag.BoolVar(&env.Verbose, "verbose", false, l10nDictionary[localization.KeyVerbose])

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		_, _ = fmt.Fprintf(w, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(w, "\n")
		_, _ = fmt.Fprintf(w, l10nDictionary[localization.KeyUsage])
	}

	json.RegisterFuzzyDecoders()
}

func main() {
	flag.Parse()
	var err error
	tools.ValidateFlags()
	if env.Verbose {
		verbose.EnableVerbose()
	}

	// load platform handler
	platformHandler, ok := handler.PlatformHandler[platform.Platform(strings.ToLower(env.Platform))]
	if !ok {
		log.Fatalf("platform %s currently not supported\n", env.Platform)
	}
	// load task definition
	jsonConfigContent := file.ReadString(env.TaskDefinitionFile)
	if jsonConfigContent == "" {
		log.Fatal("[error] empty config content, exit")
	}
	// load tasks and validate each task
	tasks := make([]*model.Task, 0)
	if err = json.JSON.Unmarshal([]byte(jsonConfigContent), &tasks); err != nil {
		log.Fatal("[error] invalid config content, exit")
	}
	if taskValidator, ok := handler.TaskValidator[platform.Platform(strings.ToLower(env.Platform))]; ok {
		for _, task := range tasks {
			if err = taskValidator(task); err != nil {
				log.Fatalf("[error] invalid task definition: %v\n", err)
			}
		}
	}

	// message (and environment params if enabled verbose)
	log.Println("This is a program that downloads live broadcast archives from asobistage, eplus, zaiko and other m3u8-base stream archives.")
	verbose.Log("Platform: %s", env.Platform)
	verbose.Log("Tasks:")
	for _, task := range tasks {
		verbose.Log("    - save to: %s", task.SaveTo)
		verbose.Log("    - page url: %s", task.PageUrl)
		verbose.Log("    - m3u8: %s", task.M3U8Url())
		if task.Spec.KeyName != "" {
			verbose.Log("    - key file: %s", task.KeyUrl())
		}
	}

	st := time.Now()
	for _, task := range tasks {
		// create dist dir
		if err = tools.CreateFolder(task.SaveTo); err != nil {
			log.Fatalf("[error] create folder: %v\n", err)
		}
		// download...
		retry, hitRetryLimit := 0, false
		for {
			if retry > env.MaxRetry {
				hitRetryLimit = true
				break
			}
			if err = platformHandler(task); err != nil {
				verbose.Log("Error: %v", err)
				log.Println("Error occurred, restarting...")
				time.Sleep(time.Second * 5)
				retry++
				continue
			} else {
				break
			}
		}

		// if hit retry limit, skip current task without validating resources
		if hitRetryLimit {
			log.Printf("Task hits retry limit, skipped. Task detail: %s\n", json.Stringify(task))
			continue
		}

		// check downloaded resources
		if err = tools.ValidateArchive(task.SaveTo); err != nil {
			log.Printf("Validate failed: %v\n", err)
		} else {
			log.Printf("Task done with: %v\n", task.Prefix)
		}
	}

	et := time.Now()
	log.Printf("[%s] Done!\n", et.Sub(st))
}
