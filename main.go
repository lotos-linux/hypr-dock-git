package main

import (
	// "os/exec"
	"fmt"
	"flag"
	"time"
	"hypr-dock/modules/cfg"
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	// "github.com/gotk3/gotk3/glib"
)

const version = "0.0.2-1-dev"

// Only during development
const CONFIG_DIR = "./configs/"
const THEMES_DIR = CONFIG_DIR + "themes/"
const MAIN_CONFIG = CONFIG_DIR + "config.jsonc"
const ITEMS_CONFIG = CONFIG_DIR + "items.json"

var config cfg.Config
var itemList cfg.ItemList

var err error

var window *gtk.Window
var app *gtk.Box
var itemsBox *gtk.Box

func initSettings() {
	configFile := flag.String("config", MAIN_CONFIG, "config file")

	config = cfg.ConnectConfig(*configFile)
	itemList = cfg.ReadItemList(ITEMS_CONFIG)

	currentTheme := flag.String("theme", config.CurrentTheme, "theme")
	config.CurrentTheme = *currentTheme

	themeConfig := cfg.ConnectConfig(
		THEMES_DIR + config.CurrentTheme + "/" + config.CurrentTheme + ".jsonc")

	config.Blur = themeConfig.Blur
	config.Spacing = themeConfig.Spacing

	flag.Parse()
}

func main() {
	initSettings()
	go initHyprEvents()

	gtk.Init(nil)

	window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		fmt.Println("Unable to create window:", err)
	}

	window.SetTitle("hypr-dock")
	orientation := setWindowProperty(window)

	err = addCssProvider(THEMES_DIR + config.CurrentTheme + "/style.css")	 
	if err != nil {
		fmt.Println(
			"CSS file not found, the default GTK theme is running!\n", err)
	}

	buildApp(orientation)

	window.Add(app)
	window.Connect("destroy", func() {gtk.MainQuit()})
	window.ShowAll()
	gtk.Main()
}

func buildApp(orientation gtk.Orientation) {
	app, _ = gtk.BoxNew(orientation, 0)
	app.SetName("app")

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

	renderItems(itemsBox)
	app.Add(itemsBox)
}

func renderItems(itemsBox *gtk.Box) {
	iconTheme, err := gtk.IconThemeGetDefault()
	if err != nil {
		fmt.Println("Unable to icon theme:", err)
	}

	for item := range len(itemList.List) {
		pixbuf, err := iconTheme.LoadIcon(
			itemList.List[item]["icon"], config.IconSize, 
			gtk.ICON_LOOKUP_FORCE_SIZE)
		if err != nil {
			fmt.Println(err)
			return
		}

		btns := map[int]*gtk.Button{}
		btns[item], _ = gtk.ButtonNew()

		imgs := map[int]*gtk.Image{}
		imgs[item], _ = gtk.ImageNewFromPixbuf(pixbuf)
		btns[item].SetImage(imgs[item])

		btns[item].Connect("clicked", func() {
			go func() {

			}()
		})

		itemsBox.Add(btns[item])
	}
}

func addCssProvider(cssFile string) error {
	cssProvider, _ := gtk.CssProviderNew()
	err := cssProvider.LoadFromPath(cssFile)

	if err == nil {
		screen, _ := gdk.ScreenGetDefault()

		gtk.AddProviderForScreen(
			screen, cssProvider, 
			gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

		return nil
	}

	return err
}

func setWindowProperty(window *gtk.Window) gtk.Orientation {
	AppOreintation := gtk.ORIENTATION_VERTICAL
	Layer := layershell.LAYER_SHELL_LAYER_BOTTOM
	Edge := layershell.LAYER_SHELL_EDGE_LEFT

	switch config.Layer {
	case "background":
		Layer = layershell.LAYER_SHELL_LAYER_BACKGROUND
	case "bottom":
		Layer = layershell.LAYER_SHELL_LAYER_BOTTOM
	case "top":
		Layer = layershell.LAYER_SHELL_LAYER_TOP
	case "overlay":
		Layer = layershell.LAYER_SHELL_LAYER_OVERLAY
	}

	switch config.Position {
	case "left":
		Edge = layershell.LAYER_SHELL_EDGE_LEFT
	case "bottom":
		Edge = layershell.LAYER_SHELL_EDGE_BOTTOM
		AppOreintation = gtk.ORIENTATION_HORIZONTAL
	case "right":
		Edge = layershell.LAYER_SHELL_EDGE_RIGHT
	case "top":
		Edge = layershell.LAYER_SHELL_EDGE_TOP
		AppOreintation = gtk.ORIENTATION_HORIZONTAL
	}

	layershell.InitForWindow(window)
	layershell.SetNamespace(window, "hypr-dock")
	layershell.SetAnchor(window, Edge, true)
	layershell.SetMargin(window, Edge, 0)

	addLayerRule()

	if config.Layer == "auto" {
		layershell.SetLayer(window, Layer)
		autoLayer()
		return AppOreintation
	}

	layershell.SetLayer(window, Layer)
	return AppOreintation
}

func autoLayer() {
	window.Connect("enter-notify-event", func(window *gtk.Window, e *gdk.Event) {
		event := gdk.EventCrossingNewFromEvent(e)
		isInWindow := event.Detail() == 3 || event.Detail() == 4 || true

		if isInWindow {
			go func() {
				setLayer("top")
			}()
		}
	})

	window.Connect("leave-notify-event", func(window *gtk.Window, e *gdk.Event) {
		event := gdk.EventCrossingNewFromEvent(e)
		isInWindow := event.Detail() == 3 || event.Detail() == 4
		
		if isInWindow {
			go func() {
				time.Sleep(time.Second / 3) 
				setLayer("bottom")
			}()
		}
	})
}

func setLayer(layer string) {
	switch layer {
	case "top":
		layershell.SetLayer(window, layershell.LAYER_SHELL_LAYER_TOP)
	case "bottom":
		layershell.SetLayer(window, layershell.LAYER_SHELL_LAYER_BOTTOM)
	}
}