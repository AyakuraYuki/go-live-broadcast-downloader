package main

import (
	"flag"
	"fmt"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal"
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

var (
	platform           string
	taskDefinitionFile string
	proxyHost          string
	proxyPort          int
	proxyType          string

	localeTag language.Tag
	err       error
)

func init() {
	localeTag, err = locale.Detect()
	if err != nil {
		localeTag = language.English
	}
	l10nDictionary := localization.GetLocalizationDictionary(localeTag)

	flag.StringVar(&platform, "p", "", l10nDictionary[localization.KeyPlatform])
	flag.StringVar(&platform, "plat", "", l10nDictionary[localization.KeyPlatform])

	flag.StringVar(&taskDefinitionFile, "c", "", l10nDictionary[localization.KeyTaskDefinitionFile])
	flag.StringVar(&taskDefinitionFile, "config", "", l10nDictionary[localization.KeyTaskDefinitionFile])

	flag.StringVar(&proxyHost, "proxy_host", "127.0.0.1", l10nDictionary[localization.KeyProxyHost])
	flag.IntVar(&proxyPort, "proxy_port", 7890, l10nDictionary[localization.KeyProxyPort])
	flag.StringVar(&proxyType, "proxy_type", "", l10nDictionary[localization.KeyProxyType])

	flag.BoolVar(&verbose.Verbose, "verbose", false, l10nDictionary[localization.KeyVerbose])

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		_, _ = fmt.Fprintf(w, "Usage of %s: %s -p <asobistage|eplus|zaiko> -c </path/to/config.json>\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(w, "\n")
		_, _ = fmt.Fprintf(w, l10nDictionary[localization.KeyUsage])
	}
}

func validateFlags() {
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

	// validate proxy flags
	if proxyType == "" {
		return // without proxy
	}
	if nhttp.MatchProxy(proxyType) == "" {
		log.Fatal("[error] we are not support the proxy type that you presented, you can only use socks5, https or http proxy")
	}
}

func main() {
	flag.Parse()
	validateFlags()

	var proxyOption *nhttp.ProxyOption
	if proxyType != "" {
		proxyOption = &nhttp.ProxyOption{
			Host:      proxyHost,
			Port:      proxyPort,
			ProxyType: proxyType,
		}
	}

	platformHandler := internal.PlatformHandler[strings.ToLower(platform)]
	if platformHandler == nil {
		log.Fatalf("platform %s currently not supported\n", platform)
	}

	tasks := make([]*internal.Task, 0)

	jsonConfigContent := file.ReadFile(taskDefinitionFile)
	if jsonConfigContent == "" {
		log.Fatal("[error] empty config content, exit")
	}
	err = cjson.JSON.Unmarshal([]byte(jsonConfigContent), &tasks)
	if err != nil {
		panic(err)
	}
	// validate
	taskValidator := internal.TaskValidator[strings.ToLower(platform)]
	for _, task := range tasks {
		if taskValidator != nil {
			err = taskValidator(task)
			if err != nil {
				panic(err)
			}
		}
	}

	log.Println("This is a program that downloads live broadcast archives from asobistage, eplus, zaiko and other m3u8-base stream archives.")
	verbose.Printf("Platform: %s\n", platform)
	verbose.Println("Tasks:")
	for _, task := range tasks {
		verbose.Printf("    - save to: %s\n", task.SaveTo)
		verbose.Printf("    - page url: %s\n", task.PageUrl)
		verbose.Printf("    - m3u8: %s\n", task.M3U8Url())
		if task.Spec.KeyName != "" {
			verbose.Printf("    - key file: %s\n", task.KeyUrl())
		}
	}

	st := time.Now()
	for _, task := range tasks {
		// create dist dir
		err = internal.CreateFolder(task.SaveTo)
		if err != nil {
			panic(err)
		}
		// download...
		for {
			err = platformHandler(task, proxyOption)
			if err != nil {
				verbose.Printf("Error: %v\n", err)
				log.Println("Error occurred, restarting...")
				time.Sleep(time.Second * 5)
				continue
			}
			break
		}
		// check downloaded resources
		err = internal.ValidateArchive(task.SaveTo)
		if err != nil {
			log.Printf("Validate failed: %v\n", err)
		} else {
			log.Printf("Task done with: %v\n", task.Prefix)
		}
	}

	et := time.Now()
	log.Printf("[%s] Done!\n", et.Sub(st))
}
