package main

import (
	"fmt"
	"slices"
	"strings"
	"strconv"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	// "github.com/gotk3/gotk3/glib"
)

var app *gtk.Box
var itemsBox *gtk.Box
var isCancelHide int
var pinnedApps []string

type appData struct {
	Instances			int
	Windows				[]map[string]string
	DesktopData			desktopData
	Button				*gtk.Button
	ButtonBox			*gtk.Box
	IndicatorImage		*gtk.Image
}
var addedApps = make(map[string]*appData)

func buildApp(orientation gtk.Orientation) {
	app, _ = gtk.BoxNew(orientation, 0)
	app.SetName("app")


	strMargin := strconv.Itoa(config.Margin)
	css := "#app {margin-"+config.Position+": "+strMargin+"px;}"

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

	renderItems(itemsBox)
	app.Add(itemsBox)
}

func renderItems(itemsBox *gtk.Box) {
	listClients()

	// Render of pinned apps
	for _, className := range pinnedApps {
		addItem(className)
	}
	
	// Render of running apps
	for _, ipcClient := range clients {
		addApp(ipcClient)
	}
}

func addApp(ipcClient client) {
	className := ipcClient.Class
	if !slices.Contains(pinnedApps, className) {
		addItem(className)
		addIndicator(className, ipcClient)
	} else {
		addIndicator(className, ipcClient)
	}
}

func addItem(className string) {
	itemProp, err := getDesktopData(className)
	if err != nil {
		fmt.Println(err)
	}

	item, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

	button, _ := gtk.ButtonNew()
	image := createImage(itemProp.Icon, config.IconSize)

	button.SetImage(image)
	button.SetName(className)
	button.SetTooltipText(itemProp.Name)

	cancelHide(button)

	indicatorImage := createImage(
		THEMES_DIR + config.CurrentTheme + "/empty.svg", config.IconSize - 10)
	indicatorImage.SetName("empty")

	button.Connect("clicked", func() {
		launch(itemProp.Exec)
	})

	addedApps[className] = &appData{
		IndicatorImage: indicatorImage,
		Button: button,
		ButtonBox: item,
		DesktopData: itemProp,
		Instances: 0,
	}

	item.Add(button)
	item.Add(indicatorImage)

	itemsBox.Add(item)
	window.ShowAll()
}

func addIndicator(className string, ipcClient client) {
	thisApp := addedApps[className]
	mainBox := thisApp.ButtonBox

	itemProp := thisApp.DesktopData

	appWindow := make(map[string]string)
	appWindow["Address"] = ipcClient.Address
	appWindow["Title"] = ipcClient.Title

	addedApps[className].Instances += 1 
	addedApps[className].Windows = append(
		addedApps[className].Windows, appWindow)


	thisApp.IndicatorImage.Destroy()
	thisApp.Button.Destroy()

	newButton, _ := gtk.ButtonNew()
	image := createImage(itemProp.Icon, config.IconSize)

	newButton.SetImage(image)
	newButton.SetName(className)
	newButton.SetTooltipText(itemProp.Name)
	// newButton.SetTooltipText(appWindow["Title"])

	var newImage *gtk.Image
	if thisApp.Instances == 1 {
		newImage = createImage(
			THEMES_DIR + config.CurrentTheme + "/single.svg", config.IconSize - 10)

		newButton.Connect("clicked", func() {
			fmt.Println(thisApp.Instances, thisApp.Windows)
			hyprctl("dispatch focuswindow address:" + ipcClient.Address)
		})

	} else {
		newImage = createImage(
			THEMES_DIR + config.CurrentTheme + "/multiple.svg", config.IconSize - 10)

		newButton.Connect("clicked", func() {
			fmt.Println(thisApp.Instances, thisApp.Windows)
			hyprctl(
				"dispatch focuswindow address:" + ipcClient.Address)
		})
	}

	addedApps[className].IndicatorImage = newImage
	addedApps[className].Button = newButton

	cancelHide(newButton)

	mainBox.Add(newButton)
	mainBox.Add(newImage)
	window.ShowAll()
}

func cancelHide(button *gtk.Button) {
	button.Connect("enter-notify-event", func() {
		isCancelHide = 1
	})
}

func createImage(source string, size int) *gtk.Image {
	iconTheme, err := gtk.IconThemeGetDefault()
	if err != nil {
		fmt.Println("Unable to icon theme:", err)
	}

	// Create image in file
	if strings.Contains(source, "/") {
		pixbuf, err := gdk.PixbufNewFromFileAtSize(
			source, size, size)
		if err != nil {
			fmt.Println(err)
			pixbuf, _ = iconTheme.LoadIcon(
				"cancel", size, gtk.ICON_LOOKUP_FORCE_SIZE)
			
		}

		image, err := gtk.ImageNewFromPixbuf(pixbuf)
		if err != nil {
			fmt.Println(err)
		}
		return image
	}

	// Create image in icon name
	pixbuf, err := iconTheme.LoadIcon(
		source, size, gtk.ICON_LOOKUP_FORCE_SIZE)
	if err != nil {
		fmt.Println(err)
	}

	image, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		fmt.Println(err)
	}

	return image
}