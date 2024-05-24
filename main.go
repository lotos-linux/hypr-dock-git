package main

import (
	"fmt"
	"flag"
	"hypr-dock/modules/cfg"
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const version = "0.0.2alpha"

// Only during development
const CONFIG_DIR = "./configs/"
const THEMES_DIR = CONFIG_DIR + "themes/"
const MAIN_CONFIG = CONFIG_DIR + "config.jsonc"
const ITEMS_CONFIG = CONFIG_DIR + "items.json"

var config cfg.Config
var itemList cfg.ItemList

var err error
var app *gtk.Box

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
	}

	buildApp(orientation)

	window.Add(app)
	window.Connect("destroy", func() {gtk.MainQuit()})
	window.ShowAll()
	gtk.Main()
}

func buildApp(orientation gtk.Orientation) {
	app, _ = gtk.BoxNew(orientation, config.Spacing)
	app.SetName("app")

	switch orientation {
	case gtk.ORIENTATION_HORIZONTAL:
		app.SetMarginEnd(config.Spacing / 2)
		app.SetMarginStart(config.Spacing / 2)
	case gtk.ORIENTATION_VERTICAL:
		app.SetMarginBottom(config.Spacing / 2)
		app.SetMarginTop(config.Spacing / 2)
	}


	renderItems(app)
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
	layershell.SetLayer(window, Layer)
	layershell.SetAnchor(window, Edge, true)
	layershell.SetMargin(window, Edge, config.Margin)

	addLayerRule()

	return AppOreintation
}