package verbose

import (
	"fmt"
	"log"
)

var (
	Verbose = false
)

func Printf(format string, v ...any) {
	if !Verbose {
		return
	}
	log.Printf(fmt.Sprintf("[verbose] %s", format), v...)
}

func Println(text string) {
	if !Verbose {
		return
	}
	log.Println(text)
}
