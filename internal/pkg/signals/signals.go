package signals

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotk3/gotk3/gtk"
)

func Handler() {
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGTERM, syscall.SIGUSR1)

	go func() {
		for {
			signalU := <-signalChanel
			switch signalU {
			case syscall.SIGTERM:
				log.Println("Exit... (SIGTERM)")
				gtk.MainQuit()
			case syscall.SIGUSR1:
				log.Println("Exit... (SIGUSR1)")
				gtk.MainQuit()
			default:
				log.Println("Unknow signal")
			}
		}
	}()
}
