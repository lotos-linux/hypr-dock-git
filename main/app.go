package main

import (
	"errors"
	"log"
	"slices"
	"strconv"
	"sync"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"hypr-dock/enternal/appItem"
	"hypr-dock/enternal/pkg/h"
	"hypr-dock/enternal/pkg/ipc"
)

var itemsBox *gtk.Box
var addedApps = make(map[string]*appItem.Item)
var mu sync.Mutex

func buildApp(orientation gtk.Orientation) *gtk.Box {
	appItem.SetConfig(config)
	h.SetConfig(config)
	ipc.SetConfig(config)

	initHandlers()

	app, err := gtk.BoxNew(orientation, 0)
	if err != nil {
		log.Fatal(err)
	}

	app.SetName("app")

	strMargin := strconv.Itoa(config.Margin)
	css := "#app {margin-" + config.Position + ": " + strMargin + "px;}"

	marginProvider, _ := gtk.CssProviderNew()
	appStyleContext, _ := app.GetStyleContext()

	appStyleContext.AddProvider(
		marginProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	marginProvider.LoadFromData(css)

	itemsBox, _ = gtk.BoxNew(orientation, config.Spacing)
	itemsBox.SetName("items-box")

	switch orientation {
	case gtk.ORIENTATION_HORIZONTAL:
		itemsBox.SetMarginEnd(config.Spacing / 2)
		itemsBox.SetMarginStart(config.Spacing / 2)
	case gtk.ORIENTATION_VERTICAL:
		itemsBox.SetMarginBottom(config.Spacing / 2)
		itemsBox.SetMarginTop(config.Spacing / 2)
	}

	renderItems()
	app.Add(itemsBox)

	return app
}

func renderItems() {
	clients, _ := ipc.ListClients()

	// Render of pinned apps
	for _, className := range pinnedApps {
		initNewItemInClass(className)
	}

	// Render of running apps
	for _, ipcClient := range clients {
		initNewItemInIPC(ipcClient)
	}
}

func initNewItemInIPC(ipcClient ipc.Client) {
	className := ipcClient.Class
	if !slices.Contains(pinnedApps, className) && addedApps[className] == nil {
		initNewItemInClass(className)
	}

	addedApps[className].UpdateState(ipcClient)
	window.ShowAll()
}

func initNewItemInClass(className string) {
	item, err := appItem.New(className)
	if err != nil {
		log.Println(err)
		return
	}

	cancelHide(item.Button)
	appItemEventHandler(item)

	item.List = addedApps
	item.PinnedList = &pinnedApps
	addedApps[className] = item

	itemsBox.Add(item.ButtonBox)
	window.ShowAll()
}

func appItemEventHandler(item *appItem.Item) {
	item.Button.Connect("button-release-event", func(button *gtk.Button, e *gdk.Event) {
		event := gdk.EventButtonNewFromEvent(e)
		if event.Button() == 3 {
			menu, err := item.ContextMenu()
			if err != nil {
				log.Println(err)
				return
			}

			menu.PopupAtWidget(item.Button, gdk.GDK_GRAVITY_NORTH, gdk.GDK_GRAVITY_SOUTH, nil)

			return
		}

		if item.Instances == 0 {
			h.Launch(item.DesktopData.Exec)
		}
		if item.Instances == 1 {
			ipc.Hyprctl("dispatch focuswindow address:" + item.Windows[0]["Address"])
		}
		if item.Instances > 1 {
			menu, err := item.WindowsMenu()
			if err != nil {
				log.Println(err)
				return
			}

			menu.PopupAtWidget(item.Button, gdk.GDK_GRAVITY_NORTH, gdk.GDK_GRAVITY_SOUTH, nil)
		}
	})
}

func removeApp(address string) {
	item, windowIndex, err := searhByAddress(address)
	if err != nil {
		log.Println(err)
		return
	}

	ipc.ListClients()

	className := item.ClassName
	if item.Instances == 1 && !slices.Contains(pinnedApps, className) {
		item.ButtonBox.Destroy()
		delete(addedApps, className)
		return
	}

	item.RemoveLastInstance(windowIndex)

	window.ShowAll()
}

func searhByAddress(address string) (*appItem.Item, int, error) {
	for _, item := range addedApps {
		for windowIndex, window := range item.Windows {
			if window["Address"] == address {
				return item, windowIndex, nil
			}
		}
	}

	err := errors.New("Window not found: " + address)
	return nil, 0, err
}

func changeWindowTitle(address string, title string) {
	mu.Lock()
	defer mu.Unlock()

	for _, data := range addedApps {
		for _, appWindow := range data.Windows {
			if appWindow["Address"] == address {
				appWindow["Title"] = title
			}
		}
	}
}
