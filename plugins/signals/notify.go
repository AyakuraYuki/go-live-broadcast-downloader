package signals

import (
	"go-live-broadcast-downloader/plugins/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GraceStop(callback func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	s := <-ch
	log.Info("GraceStop").Msgf("server shutdown at:%s,[%v]", time.Now().String(), s)

	callback()
	os.Exit(0)
}
