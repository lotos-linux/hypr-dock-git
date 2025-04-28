package layering

import (
	"hypr-dock/internal/state"
	"log"
	"strings"
	"time"

	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func SetWindowProperty(appState *state.State) {
	window := appState.GetWindow()
	settings := appState.GetSettings()

	layershell.InitForWindow(window)
	layershell.SetNamespace(window, "hypr-dock")

	ChangePosition(settings.Position, appState)
	ChangeLayer(settings.Layer, appState)
}

func ChangeLayer(layer string, appState *state.State) {
	window := appState.GetWindow()
	if window == nil {
		return
	}

	layers := map[string]layershell.LayerShellLayerFlags{
		"background": layershell.LAYER_SHELL_LAYER_BACKGROUND,
		"bottom":     layershell.LAYER_SHELL_LAYER_BOTTOM,
		"top":        layershell.LAYER_SHELL_LAYER_TOP,
		"overlay":    layershell.LAYER_SHELL_LAYER_OVERLAY,
	}

	if layer == "auto" {
		layershell.SetLayer(window, layers["bottom"])
		layershell.SetExclusiveZone(window, 0)
		AutoLayer(appState)
		InitDetectArea(appState)
		return
	}

	DisableAutoLayer(appState)
	if strings.Contains(layer, "exclusive") {
		exLayer := strings.Split(layer, "-")[1]
		layershell.SetLayer(window, layers[exLayer])
		layershell.AutoExclusiveZoneEnable(window)
		return
	}

	layershell.SetLayer(window, layers[layer])
}

func ChangePosition(position string, appState *state.State) {
	window := appState.GetWindow()
	if window == nil {
		return
	}

	oreintations := map[string]gtk.Orientation{
		"bottom": gtk.ORIENTATION_HORIZONTAL,
		"top":    gtk.ORIENTATION_HORIZONTAL,
		"left":   gtk.ORIENTATION_VERTICAL,
		"right":  gtk.ORIENTATION_VERTICAL,
	}

	edges := map[string]layershell.LayerShellEdgeFlags{
		"bottom": layershell.LAYER_SHELL_EDGE_BOTTOM,
		"top":    layershell.LAYER_SHELL_EDGE_TOP,
		"left":   layershell.LAYER_SHELL_EDGE_LEFT,
		"right":  layershell.LAYER_SHELL_EDGE_RIGHT,
	}

	layershell.SetAnchor(window, edges[position], true)
	layershell.SetMargin(window, edges[position], 0)

	appState.SetOrientation(oreintations[position])
	appState.SetEdge(edges[position])
}

func InitDetectArea(appState *state.State) {
	edge := appState.GetEdge()
	window := appState.GetWindow()

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

	long := appState.GetSettings().IconSize * appState.GetList().Len() * 2

	switch appState.GetOrientation() {
	case gtk.ORIENTATION_HORIZONTAL:
		detectArea.SetSizeRequest(long, 1)
	case gtk.ORIENTATION_VERTICAL:
		detectArea.SetSizeRequest(1, long)
	}

	detectArea.Connect("enter-notify-event", func(detectWindow *gtk.Window, e *gdk.Event) {
		appState.SetPreventHide(false)
		go func() {
			layershell.SetLayer(window, layershell.LAYER_SHELL_LAYER_TOP)
		}()
	})

	detectArea.ShowAll()
	appState.SetDetectArea(detectArea)
}

func AutoLayer(appState *state.State) {
	DisableAutoLayer(appState)
	window := appState.GetWindow()

	enterSig := window.Connect("enter-notify-event", func(window *gtk.Window, e *gdk.Event) {
		event := gdk.EventCrossingNewFromEvent(e)
		isInWindow := event.Detail() == 3 || event.Detail() == 4

		if isInWindow {
			appState.SetPreventHide(true)
		}

		if isInWindow && !appState.GetSpecial() {
			go func() {
				layershell.SetLayer(window, layershell.LAYER_SHELL_LAYER_TOP)
				appState.SetPreventHide(true)
			}()
		}
	})
	appState.AddSignalHandler("enter", enterSig)

	leaveSig := window.Connect("leave-notify-event", func(window *gtk.Window, e *gdk.Event) {
		DispathLeaveEvent(window, e, appState)
	})
	appState.AddSignalHandler("leave", leaveSig)
}

func DispathLeaveEvent(window *gtk.Window, e *gdk.Event, appState *state.State) {
	var isInWindow bool
	if e != nil {
		event := gdk.EventCrossingNewFromEvent(e)
		isInWindow = event.Detail() == 3 || event.Detail() == 4
	} else {
		isInWindow = true
	}

	if isInWindow {
		appState.SetPreventHide(false)
	}

	if isInWindow && !appState.GetPreventHide() {
		go func() {
			time.Sleep(time.Second / 3)
			if !appState.GetPreventHide() {
				layershell.SetLayer(window, layershell.LAYER_SHELL_LAYER_BOTTOM)
				appState.SetPreventHide(false)
			}
		}()
	}
}

func DisableAutoLayer(appState *state.State) {
	detectArea := appState.GetDetectArea()
	if detectArea != nil {
		detectArea.Destroy()
		appState.SetDetectArea(nil)
	}

	window := appState.GetWindow()

	appState.RemoveSignalHandler("enter", window)
	appState.RemoveSignalHandler("leave", window)

	appState.SetPreventHide(false)
}
