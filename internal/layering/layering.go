package layering

import (
	"hypr-dock/internal/settings"
	"hypr-dock/internal/state"
	"log"
	"time"

	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func SetWindowProperty(window *gtk.Window, appState *state.State) (gtk.Orientation, layershell.LayerShellEdgeFlags) {
	oreintation := gtk.ORIENTATION_HORIZONTAL
	layer := layershell.LAYER_SHELL_LAYER_BOTTOM
	edge := layershell.LAYER_SHELL_EDGE_BOTTOM

	switch settings.Get().Position {
	case "left":
		edge = layershell.LAYER_SHELL_EDGE_LEFT
		oreintation = gtk.ORIENTATION_VERTICAL
	case "bottom":
		edge = layershell.LAYER_SHELL_EDGE_BOTTOM
	case "right":
		edge = layershell.LAYER_SHELL_EDGE_RIGHT
		oreintation = gtk.ORIENTATION_VERTICAL
	case "top":
		edge = layershell.LAYER_SHELL_EDGE_TOP
	}

	layershell.InitForWindow(window)
	layershell.SetNamespace(window, "hypr-dock")
	layershell.SetAnchor(window, edge, true)
	layershell.SetMargin(window, edge, 0)

	if settings.Get().Layer == "auto" {
		SetLayer("bottom", appState)
		AutoLayer(window, appState)
		return oreintation, edge
	}

	switch settings.Get().Layer {
	case "background":
		layer = layershell.LAYER_SHELL_LAYER_BACKGROUND
	case "bottom":
		layer = layershell.LAYER_SHELL_LAYER_BOTTOM
	case "top":
		layer = layershell.LAYER_SHELL_LAYER_TOP
	case "overlay":
		layer = layershell.LAYER_SHELL_LAYER_OVERLAY
	}

	layershell.SetLayer(window, layer)
	return oreintation, edge
}

func InitDetectArea(edge layershell.LayerShellEdgeFlags, appState *state.State) {
	detectArea, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("InitDetectArea(), gtk.WindowNew() | ", err)
	}
	detectArea.SetName("detect")

	layershell.InitForWindow(detectArea)
	layershell.SetNamespace(detectArea, "dock-detect")
	layershell.SetAnchor(detectArea, edge, true)
	layershell.SetMargin(detectArea, edge, 0)
	layershell.SetLayer(detectArea, layershell.LAYER_SHELL_LAYER_TOP)

	long := settings.Get().IconSize * len(appState.GetAddedApps().List) * 2

	switch appState.GetOrientation() {
	case gtk.ORIENTATION_HORIZONTAL:
		detectArea.SetSizeRequest(long, 1)
	case gtk.ORIENTATION_VERTICAL:
		detectArea.SetSizeRequest(1, long)
	}

	detectArea.Connect("enter-notify-event", func(window *gtk.Window, e *gdk.Event) {
		appState.SetPreventHide(false)

		go func() {
			SetLayer("top", appState)
		}()
	})

	detectArea.ShowAll()
	appState.SetDetectArea(detectArea)
}

func AutoLayer(window *gtk.Window, appState *state.State) {
	window.Connect("enter-notify-event", func(window *gtk.Window, e *gdk.Event) {
		event := gdk.EventCrossingNewFromEvent(e)
		isInWindow := event.Detail() == 3 || event.Detail() == 4

		if isInWindow {
			appState.SetPreventHide(true)
		}

		if isInWindow && !appState.GetSpecial() {
			go func() {
				SetLayer("top", appState)
				appState.SetPreventHide(true)
			}()
		}
	})

	window.Connect("leave-notify-event", func(window *gtk.Window, e *gdk.Event) {
		event := gdk.EventCrossingNewFromEvent(e)
		isInWindow := event.Detail() == 3 || event.Detail() == 4

		if isInWindow {
			appState.SetPreventHide(false)
		}

		if isInWindow && !appState.GetPreventHide() {
			go func() {
				time.Sleep(time.Second / 3)
				if !appState.GetPreventHide() {
					SetLayer("bottom", appState)
					appState.SetPreventHide(false)
				}
			}()
		}
	})
}

func SetLayer(layer string, appState *state.State) {
	window := appState.GetWindow()
	switch layer {
	case "top":
		layershell.SetLayer(window, layershell.LAYER_SHELL_LAYER_TOP)
	case "bottom":
		layershell.SetLayer(window, layershell.LAYER_SHELL_LAYER_BOTTOM)
	}
}
