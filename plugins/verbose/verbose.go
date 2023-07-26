package verbose

import (
	"fmt"
	"log"
)

const (
	colorBlack   = uint8(iota + 30) // 黑
	colorRed                        // 红
	colorGreen                      // 绿
	colorYellow                     // 黄
	colorBlue                       // 蓝
	colorMagenta                    // 紫红
	colorCyan                       // 青蓝
	colorWhite                      // 白
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
	log.Println(fmt.Sprintf("[verbose] %s", text))
}
