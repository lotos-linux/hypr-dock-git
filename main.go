package main

import (
	"strconv"
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const version = "0.0.1"

func main() {
	gtk.Init(nil)

	window, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("inclayer-go")

	addLayerShell(window)
	addCssProvider("/home/lotsmannc/repos/gotk/style.css")


	mainBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	mainBox.SetName("main-box")
	window.Add(mainBox)
	

	iconTheme, _ := gtk.IconThemeGetDefault()
	pixbuf, _ := iconTheme.LoadIcon(
		"system-file-manager", 22, gtk.ICON_LOOKUP_FORCE_SIZE)


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

func addLayerShell(window *gtk.Window) {
	layershell.InitForWindow(window)
	layershell.SetNamespace(window, "inclayer-go")
	layershell.SetLayer(window, layershell.LAYER_SHELL_LAYER_BOTTOM)
	layershell.SetAnchor(window, layershell.LAYER_SHELL_EDGE_LEFT, true)
	layershell.SetMargin(window, layershell.LAYER_SHELL_EDGE_LEFT, 10)
}

func increment(label *gtk.Label, inc int) {
	labelNum, _ := strconv.Atoi(label.GetLabel())
	label.SetLabel(strconv.Itoa(labelNum + inc))
}