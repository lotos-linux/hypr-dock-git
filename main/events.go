package main

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/glib"

	"hypr-dock/enternal/pkg/ipc"
)

var special = false

func initHandlers() {
	ipc.NewEventHandler("configreloaded", func(event string) {
		go ipc.AddLayerRule()
	})

	ipc.NewEventHandler("windowtitlev2", windowTitleHandler)

	ipc.NewEventHandler("openwindow", openwindowHandler)

	ipc.NewEventHandler("closewindow", closewindowHandler)

	ipc.NewEventHandler("activespecial", activatespecialHandler)
}

func windowTitleHandler(event string) {
	data := eventHandler(event, 2)
	address := "0x" + strings.TrimSpace(data[0])
	go changeWindowTitle(address, data[1])
}

func activatespecialHandler(event string) {
	data := eventHandler(event, 2)

	if data[0] == "special:special" {
		special = true
	}
	if data[0] != "special:special" {
		special = false
	}
}

func openwindowHandler(event string) {
	data := eventHandler(event, 4)
	address := "0x" + strings.TrimSpace(data[0])
	windowClient, err := ipc.SearchClientByAddress(address)
	if err != nil {
		fmt.Println(err)
	} else {
		glib.IdleAdd(func() {
			initNewItemInIPC(windowClient)
		})
	}
}

func closewindowHandler(event string) {
	data := eventHandler(event, 1)
	address := "0x" + strings.TrimSpace(data[0])
	glib.IdleAdd(func() {
		removeApp(address)
	})
}

func eventHandler(event string, n int) []string {
	parts := strings.SplitN(event, ">>", 2)
	dataClast := strings.TrimSpace(parts[1])
	dataParts := strings.SplitN(dataClast, ",", n)

	for i := range dataParts {
		dataParts[i] = strings.TrimSpace(dataParts[i])
	}

	return dataParts
}
