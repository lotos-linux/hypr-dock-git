package main

import (
	"fmt"
	"strconv"
	"hypr-dock/modules/cfg"
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const version = "0.0.2alpha"

const CONFIG_DIR = "./configs/"
const THEMES_DIR = CONFIG_DIR + "themes/"
const MAIN_CONFIG = CONFIG_DIR + "main.jsonc"

func main() {
	fmt.Println("Start")
	config := cfg.Connect(MAIN_CONFIG)

	gtk.Init(nil)

	window, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("hypr-dock")

	orientation := setWindowProperty(window, config.Layer, 
	                                 config.Position, config.Margin)
									 
	addCssProvider(THEMES_DIR + config.CurrentTheme + "/style.css")


	mainBox, _ := gtk.BoxNew(orientation, 0)
	mainBox.SetName("main-box")
	window.Add(mainBox)
	

	iconTheme, _ := gtk.IconThemeGetDefault()
	pixbuf, _ := iconTheme.LoadIcon(
		"system-file-manager", config.IconSize, gtk.ICON_LOOKUP_FORCE_SIZE)


	label, _ := gtk.LabelNew("0")
	label.SetName("number-box")

	
	for number := range 6 {
		btns := map[int]*gtk.Button{}
		btns[number], _ = gtk.ButtonNew()

		imgs := map[int]*gtk.Image{}
		imgs[number], _ = gtk.ImageNewFromPixbuf(pixbuf)
		btns[number].SetImage(imgs[number])

		btns[number].Connect("clicked", func() {
			increment(label, number + 1)
		})
		mainBox.Add(btns[number])
	}


	mainBox.Add(label)


	window.Connect("destroy", func() {gtk.MainQuit()})
	window.ShowAll()
	gtk.Main()
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

func setWindowProperty(window *gtk.Window, 
					   layer string, 
					   position string, 
					   margin int) gtk.Orientation {

	LAYER_SHELL_LAYER := layershell.LAYER_SHELL_LAYER_BOTTOM
	LAYER_SHELL_EDGE := layershell.LAYER_SHELL_EDGE_LEFT
	MAIN_BOX_ORIENTATION := gtk.ORIENTATION_VERTICAL


	switch layer {
	case "background":
		LAYER_SHELL_LAYER = layershell.LAYER_SHELL_LAYER_BACKGROUND
	case "bottom":
		LAYER_SHELL_LAYER = layershell.LAYER_SHELL_LAYER_BOTTOM
	case "top":
		LAYER_SHELL_LAYER = layershell.LAYER_SHELL_LAYER_TOP
	case "overlay":
		LAYER_SHELL_LAYER = layershell.LAYER_SHELL_LAYER_OVERLAY
	}

	switch position {
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
	layershell.SetMargin(window, LAYER_SHELL_EDGE, margin)

	return MAIN_BOX_ORIENTATION
}

func increment(label *gtk.Label, inc int) {
	labelNum, _ := strconv.Atoi(label.GetLabel())
	label.SetLabel(strconv.Itoa(labelNum + inc))
}