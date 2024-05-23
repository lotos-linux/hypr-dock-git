package main

import (
	"fmt"
	"hypr-dock/modules/cfg"
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const version = "0.0.2alpha"

const CONFIG_DIR = "./configs/"
const THEMES_DIR = CONFIG_DIR + "themes/"
const MAIN_CONFIG = CONFIG_DIR + "config.jsonc"
const ITEMS_CONFIG = CONFIG_DIR + "items.json"

var config = cfg.ConnectConfig(MAIN_CONFIG)
var itemList = cfg.ReadItemList(ITEMS_CONFIG)

var err error
var app *gtk.Box

func main() {
	gtk.Init(nil)

	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		fmt.Println("Unable to create window:", err)
	}

	window.SetTitle("hypr-dock")
	orientation := setWindowProperty(window)

	err = addCssProvider(THEMES_DIR + config.CurrentTheme + "/style.css")	 
	if err != nil {
		fmt.Println(
			"CSS file not found, the default GTK theme is running!\n", err)
		app, _ = gtk.BoxNew(orientation, 5)
	} else {
		app, _ = gtk.BoxNew(orientation, 0)
		app.SetName("app")
	}

	renderItems(app)

	window.Add(app)
	window.Connect("destroy", func() {gtk.MainQuit()})
	window.ShowAll()
	gtk.Main()
}

func renderItems(app *gtk.Box) {
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

		app.Add(btns[item])
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
	} else {

		return err
	}
}

func setWindowProperty(window *gtk.Window) gtk.Orientation {
	LAYER_SHELL_LAYER := layershell.LAYER_SHELL_LAYER_BOTTOM
	LAYER_SHELL_EDGE := layershell.LAYER_SHELL_EDGE_LEFT
	APP_ORIENTATION := gtk.ORIENTATION_VERTICAL


	switch config.Layer {
	case "background":
		LAYER_SHELL_LAYER = layershell.LAYER_SHELL_LAYER_BACKGROUND
	case "bottom":
		LAYER_SHELL_LAYER = layershell.LAYER_SHELL_LAYER_BOTTOM
	case "top":
		LAYER_SHELL_LAYER = layershell.LAYER_SHELL_LAYER_TOP
	case "overlay":
		LAYER_SHELL_LAYER = layershell.LAYER_SHELL_LAYER_OVERLAY
	}

	switch config.Position {
	case "left":
		LAYER_SHELL_EDGE = layershell.LAYER_SHELL_EDGE_LEFT
	case "bottom":
		LAYER_SHELL_EDGE = layershell.LAYER_SHELL_EDGE_BOTTOM
		APP_ORIENTATION = gtk.ORIENTATION_HORIZONTAL
	case "right":
		LAYER_SHELL_EDGE = layershell.LAYER_SHELL_EDGE_RIGHT
	case "top":
		LAYER_SHELL_EDGE = layershell.LAYER_SHELL_EDGE_TOP
		APP_ORIENTATION = gtk.ORIENTATION_HORIZONTAL
	}

	layershell.InitForWindow(window)
	layershell.SetNamespace(window, "hypr-dock")
	layershell.SetLayer(window, LAYER_SHELL_LAYER)
	layershell.SetAnchor(window, LAYER_SHELL_EDGE, true)
	layershell.SetMargin(window, LAYER_SHELL_EDGE, config.Margin)

	return APP_ORIENTATION
}