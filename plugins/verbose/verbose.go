package verbose

import (
	"fmt"
	"log"
)

// font color
const (
	colorBlack = uint8(iota + 30)
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)

// font background color
const (
	bgColorBlack = uint8(iota + 40)
	bgColorRed
	bgColorGreen
	bgColorYellow
	bgColorBlue
	bgColorMagenta
	bgColorCyan
	bgColorWhite
)

// control flag
const (
	removeColorSuffix = "\x1b[0m"
)

var (
	Verbose = false
)

func Printf(format string, v ...any) {
	if !Verbose {
		return
	}
	log.Printf(colors(fmt.Sprintf("[verbose] %s", format), colorWhite), v...)
}

func Println(text string) {
	if !Verbose {
		return
	}
	log.Println(colors(fmt.Sprintf("[verbose] %s", text), colorWhite))
}

func colors(str string, colorCode uint8) string {
	return fmt.Sprintf("\x1b[%dm%s%s", colorCode, str, removeColorSuffix)
}
