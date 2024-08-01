package verbose

import (
	"log"

	"github.com/AyakuraYuki/go-live-broadcast-downloader/live-broadcast-downloader/internal/colors"
)

var verbose = false

func EnableVerbose()  { verbose = true }
func DisableVerbose() { verbose = false }

func Log(format string, v ...any) {
	if verbose {
		log.Println(colors.White(format, v...))
	}
}
