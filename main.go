package main

import (
	"fmt"
	// "strconv"
	"hypr-dock/modules/cfg"
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const version = "0.0.2alpha"

const CONFIG_DIR = "./configs/"
const THEMES_DIR = CONFIG_DIR + "themes/"
const MAIN_CONFIG = CONFIG_DIR + "main.jsonc"
const ITEMS_CONFIG = CONFIG_DIR + "items.json"

var config = cfg.ConnectConfig(MAIN_CONFIG)
var itemList = cfg.ReadItemList(ITEMS_CONFIG)

func main() {
	fmt.Println("Start")

	gtk.Init(nil)

	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		fmt.Println("Unable to create window:", err)
	}

	window.SetTitle("hypr-dock")
	addCssProvider(THEMES_DIR + config.CurrentTheme + "/style.css")


	orientation := setWindowProperty(window)			 


	mainBox, _ := gtk.BoxNew(orientation, 0)
	mainBox.SetName("main-box")
	window.Add(mainBox)

	renderItems(mainBox, window)


	window.Connect("destroy", func() {gtk.MainQuit()})
	window.ShowAll()
	gtk.Main()
}

func renderItems(mainBox *gtk.Box, window *gtk.Window) {
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

		fmt.Println(itemList.List[item]["title"])

		btns := map[int]*gtk.Button{}
		btns[item], _ = gtk.ButtonNew()

		imgs := map[int]*gtk.Image{}
		imgs[item], _ = gtk.ImageNewFromPixbuf(pixbuf)
		btns[item].SetImage(imgs[item])

		// btns[item].Connect("clicked", func() {
		// 	layershell.SetMargin(
		// 		window, layershell.LAYER_SHELL_EDGE_LEFT, config.IconSize * -2)
		// })

		mainBox.Add(btns[item])
	}
}

func addCssProvider(cssFile string) {
	cssProvider, _ := gtk.CssProviderNew()
	err := cssProvider.LoadFromPath(cssFile)
	if err == nil {
		screen, _ := gdk.ScreenGetDefault()
		gtk.AddProviderForScreen(
			screen, cssProvider,gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	}
}

func setWindowProperty(window *gtk.Window) gtk.Orientation {

	LAYER_SHELL_LAYER := layershell.LAYER_SHELL_LAYER_BOTTOM
	LAYER_SHELL_EDGE := layershell.LAYER_SHELL_EDGE_LEFT
	MAIN_BOX_ORIENTATION := gtk.ORIENTATION_VERTICAL


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
		MAIN_BOX_ORIENTATION = gtk.ORIENTATION_HORIZONTAL
	case "right":
		LAYER_SHELL_EDGE = layershell.LAYER_SHELL_EDGE_RIGHT
	case "top":
		LAYER_SHELL_EDGE = layershell.LAYER_SHELL_EDGE_TOP
		MAIN_BOX_ORIENTATION = gtk.ORIENTATION_HORIZONTAL
	}

	layershell.InitForWindow(window)
	layershell.SetNamespace(window, "hypr-dock")
	layershell.SetLayer(window, LAYER_SHELL_LAYER)
	layershell.SetAnchor(window, LAYER_SHELL_EDGE, true)
	layershell.SetMargin(window, LAYER_SHELL_EDGE, config.Margin)

	return MAIN_BOX_ORIENTATION
}