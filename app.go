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
var addedItems []string

type buttonList struct {
	IndicatorImage		*gtk.Image
	Button				*gtk.Button
	ButtonBox			*gtk.Box
	ClientData			clientData
}
var addedWidget = make(map[string]*buttonList)

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

	for item := range len(pinnedApps) {
		addItem(pinnedApps[item])
		addedItems = append(addedItems, pinnedApps[item])
	}
	
	for item := range len(clients) {
		className := clients[item].Class
		if !slices.Contains(addedItems, className) {
			addItem(className)
			addIndicator(className)
			addedItems = append(addedItems, className)
		} else {
			addIndicator(className)
		}
	}
}

func addIndicator(className string) {
	widget := addedWidget[className]
	mainBox := widget.ButtonBox

	itemProp := widget.ClientData

	widget.IndicatorImage.Destroy()
	widget.Button.Destroy()

	newButton, _ := gtk.ButtonNew()
	image := createImage(itemProp.Icon, config.IconSize)

	newButton.SetImage(image)
	newButton.SetName(className)
	newButton.SetTooltipText(itemProp.Name)

	var newImage *gtk.Image
	imageName, _ := widget.IndicatorImage.GetName()
	if imageName == "empty" {
		newImage = createImage(
			THEMES_DIR + config.CurrentTheme + "/single.svg", config.IconSize - 10)
		newImage.SetName("single")

		newButton.Connect("clicked", func() {
			fmt.Println(itemProp.Exec)
		})

	} else {
		newImage = createImage(
			THEMES_DIR + config.CurrentTheme + "/multiple.svg", config.IconSize - 10)
		newImage.SetName("multiple")

		newButton.Connect("clicked", func() {
			fmt.Println(itemProp.Exec)
		})
	}

	addedWidget[className].IndicatorImage = newImage
	addedWidget[className].Button = newButton

	cancelHide(newButton)

	mainBox.Add(newButton)
	mainBox.Add(newImage)
	window.ShowAll()
}


func addItem(className string) {
	itemProp, err := getClientData(className)
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

	addedWidget[className] = &buttonList{
		IndicatorImage: indicatorImage,
		Button: button,
		ButtonBox: item,
		ClientData: itemProp,
	}

	item.Add(button)
	item.Add(indicatorImage)

	itemsBox.Add(item)
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
				"steam", size, gtk.ICON_LOOKUP_FORCE_SIZE)
			
		}

		image, err := gtk.ImageNewFromPixbuf(pixbuf)
		if err != nil {
			fmt.Println(err)
		}
		return image
	}

	// Create image in icon name
	pixbuf, err := iconTheme.LoadIcon(
		source, config.IconSize, gtk.ICON_LOOKUP_FORCE_SIZE)
	if err != nil {
		fmt.Println(err)
	}

	image, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		fmt.Println(err)
	}

	return image
}