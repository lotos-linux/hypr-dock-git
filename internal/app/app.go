package app

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"hypr-dock/internal/hypr/hyprOpt"
	"hypr-dock/internal/item"
	"hypr-dock/internal/pkg/utils"
	"hypr-dock/internal/settings"
	"hypr-dock/internal/state"
	"hypr-dock/pkg/ipc"
)

func BuildApp(orientation gtk.Orientation, appState *state.State) *gtk.Box {
	app, err := gtk.BoxNew(orientation, 0)
	if err != nil {
		log.Fatal(err)
	}

	addWindowMarginRule(app)
	app.SetName("app")

	itemsBox, _ := gtk.BoxNew(orientation, settings.Get().Spacing)
	itemsBox.SetName("items-box")

	switch orientation {
	case gtk.ORIENTATION_HORIZONTAL:
		itemsBox.SetMarginEnd(int(float64(settings.Get().Spacing) * 0.8))
		itemsBox.SetMarginStart(int(float64(settings.Get().Spacing) * 0.8))
	case gtk.ORIENTATION_VERTICAL:
		itemsBox.SetMarginBottom(int(float64(settings.Get().Spacing) * 0.8))
		itemsBox.SetMarginTop(int(float64(settings.Get().Spacing) * 0.8))
	}

	appState.SetItemsBox(itemsBox)
	renderItems(appState)
	app.Add(itemsBox)

	return app
}

func renderItems(appState *state.State) {
	clients, _ := ipc.GetClients()

	for _, className := range settings.PinnedApps {
		InitNewItemInClass(className, appState)
	}

	for _, ipcClient := range clients {
		InitNewItemInIPC(ipcClient, appState)
	}
}

func InitNewItemInIPC(ipcClient ipc.Client, appState *state.State) {
	className := ipcClient.Class
	if !slices.Contains(settings.PinnedApps, className) && appState.GetAddedApps().List[className] == nil {
		InitNewItemInClass(className, appState)
	}

	appState.GetAddedApps().List[className].UpdateState(ipcClient)
	appState.GetWindow().ShowAll()
}

func InitNewItemInClass(className string, appState *state.State) {
	item, err := item.New(className)
	if err != nil {
		log.Println(err)
		return
	}

	appItemEventHandler(item)

	item.List = appState.GetAddedApps().List
	item.PinnedList = &settings.PinnedApps
	appState.GetAddedApps().Add(className, item)

	appState.GetItemsBox().Add(item.ButtonBox)
	appState.GetWindow().ShowAll()
}

func appItemEventHandler(item *item.Item) {
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
			utils.Launch(item.DesktopData.Exec)
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

func RemoveApp(address string, appState *state.State) {
	item, windowIndex, err := searhByAddress(address, appState)
	if err != nil {
		log.Println(err)
		return
	}

	className := item.ClassName
	if item.Instances == 1 && !slices.Contains(settings.PinnedApps, className) {
		item.Remove()
		return
	}

	item.RemoveLastInstance(windowIndex)

	appState.GetWindow().ShowAll()
}

func searhByAddress(address string, appState *state.State) (*item.Item, int, error) {
	for _, item := range appState.GetAddedApps().List {
		for windowIndex, window := range item.Windows {
			if window["Address"] == address {
				return item, windowIndex, nil
			}
		}
	}

	err := errors.New("Window not found: " + address)
	return nil, 0, err
}

func ChangeWindowTitle(address string, title string, appState *state.State) {
	for _, data := range appState.GetAddedApps().List {
		for _, appWindow := range data.Windows {
			if appWindow["Address"] == address {
				appWindow["Title"] = title
			}
		}
	}
}

func addWindowMarginRule(app *gtk.Box) {
	position := settings.Get().Position
	var marginProvider *gtk.CssProvider

	switch settings.Get().SystemGapUsed {
	case "true":
		margin, err := hyprOpt.GetGap()
		if err != nil {
			log.Println(err, "\nSet margin in config")
			applyWindowMarginCSS(app, position, settings.Get().Margin)
		}

		marginProvider = applyWindowMarginCSS(app, position, margin[0])

		hyprOpt.GapChangeEvent(func(gap int) {
			utils.RemoveStyleProvider(app, marginProvider)
			marginProvider = applyWindowMarginCSS(app, position, gap)
			log.Println("Window margins updated successfully: ", gap)
		})
	case "false":
		applyWindowMarginCSS(app, position, settings.Get().Margin)
	}
}

func applyWindowMarginCSS(app *gtk.Box, position string, margin int) *gtk.CssProvider {
	css := fmt.Sprintf("#app {margin-%s: %dpx;}", position, margin)

	marginProvider, err := gtk.CssProviderNew()
	if err != nil {
		log.Printf("Failed to create CSS provider: %v", err)
		return nil
	}

	appStyleContext, err := app.GetStyleContext()
	if err != nil {
		log.Printf("Failed to get style context: %v", err)
		return nil
	}

	appStyleContext.AddProvider(marginProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	err = marginProvider.LoadFromData(css)
	if err != nil {
		log.Printf("Failed to load CSS data: %v", err)
		return nil
	}

	return marginProvider
}
