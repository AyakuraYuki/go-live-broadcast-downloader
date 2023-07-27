package main

import (
	"flag"
	"fmt"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/env"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/handler"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/model"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/tools"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/file"
	cjson "github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/json"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/localization"
	nhttp "github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/net/http"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/verbose"
	"github.com/Xuanwo/go-locale"
	"golang.org/x/text/language"
	"log"
	"os"
	"strings"
	"time"
)

func init() {
	var err error

	env.LocaleTag, err = locale.Detect()
	if err != nil {
		env.LocaleTag = language.English
	}
	l10nDictionary := localization.GetLocalizationDictionary(env.LocaleTag)

	flag.StringVar(&env.Platform, "p", "", l10nDictionary[localization.KeyPlatform])
	flag.StringVar(&env.Platform, "plat", "", l10nDictionary[localization.KeyPlatform])

	flag.StringVar(&env.TaskDefinitionFile, "c", "", l10nDictionary[localization.KeyTaskDefinitionFile])
	flag.StringVar(&env.TaskDefinitionFile, "config", "", l10nDictionary[localization.KeyTaskDefinitionFile])

	flag.StringVar(&env.ProxyHost, "proxy_host", "127.0.0.1", l10nDictionary[localization.KeyProxyHost])
	flag.IntVar(&env.ProxyPort, "proxy_port", 7890, l10nDictionary[localization.KeyProxyPort])
	flag.StringVar(&env.ProxyType, "proxy_type", "", l10nDictionary[localization.KeyProxyType])

	flag.IntVar(&env.Coroutines, "threads", 10, l10nDictionary[localization.KeyCoroutines])
	flag.IntVar(&env.MaxRetry, "retry", 10, l10nDictionary[localization.KeyMaxRetry])
	flag.BoolVar(&verbose.Verbose, "verbose", false, l10nDictionary[localization.KeyVerbose])

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		_, _ = fmt.Fprintf(w, "Usage of %s: %s -p <asobistage|eplus|zaiko> -c </path/to/config.json>\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(w, "\n")
		_, _ = fmt.Fprintf(w, l10nDictionary[localization.KeyUsage])
	}
}

func main() {
	flag.Parse()
	tools.ValidateFlags()

	var err error

	var proxyOption *nhttp.ProxyOption
	if env.ProxyType != "" {
		proxyOption = &nhttp.ProxyOption{
			Host:      env.ProxyHost,
			Port:      env.ProxyPort,
			ProxyType: env.ProxyType,
		}
	}

	platformHandler := handler.PlatformHandler[strings.ToLower(env.Platform)]
	if platformHandler == nil {
		log.Fatalf("platform %s currently not supported\n", env.Platform)
	}

	tasks := make([]*model.Task, 0)

	jsonConfigContent := file.ReadFile(env.TaskDefinitionFile)
	if jsonConfigContent == "" {
		log.Fatal("[error] empty config content, exit")
	}
	err = cjson.JSON.Unmarshal([]byte(jsonConfigContent), &tasks)
	if err != nil {
		panic(err)
	}
	// validate
	taskValidator := handler.TaskValidator[strings.ToLower(env.Platform)]
	for _, task := range tasks {
		if taskValidator != nil {
			err = taskValidator(task)
			if err != nil {
				panic(err)
			}
		}
	}

	log.Println("This is a program that downloads live broadcast archives from asobistage, eplus, zaiko and other m3u8-base stream archives.")
	verbose.Printf("Platform: %s", env.Platform)
	verbose.Println("Tasks:")
	for _, task := range tasks {
		verbose.Printf("    - save to: %s", task.SaveTo)
		verbose.Printf("    - page url: %s", task.PageUrl)
		verbose.Printf("    - m3u8: %s", task.M3U8Url())
		if task.Spec.KeyName != "" {
			verbose.Printf("    - key file: %s", task.KeyUrl())
		}
	}

	st := time.Now()
	for _, task := range tasks {
		// create dist dir
		err = tools.CreateFolder(task.SaveTo)
		if err != nil {
			panic(err)
		}
		// download...
		retry := 0
		hitRetryLimit := false
		for {
			if retry > env.MaxRetry {
				hitRetryLimit = true
				break
			}
			err = platformHandler(task, proxyOption)
			if err != nil {
				verbose.Printf("Error: %v", err)
				log.Println("Error occurred, restarting...")
				time.Sleep(time.Second * 5)
				retry++
				continue
			}
			break
		}

		// if hit retry limit, skip current task without validating resources
		if hitRetryLimit {
			bs, _ := cjson.JSON.Marshal(task)
			log.Printf("Task hits retry limit, skipped. Task detail: %s\n", string(bs))
			continue
		}

		// check downloaded resources
		err = tools.ValidateArchive(task.SaveTo)
		if err != nil {
			log.Printf("Validate failed: %v\n", err)
		} else {
			log.Printf("Task done with: %v\n", task.Prefix)
		}
	}

	et := time.Now()
	log.Printf("[%s] Done!\n", et.Sub(st))
}
