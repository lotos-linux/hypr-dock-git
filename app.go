package main

import (
	"fmt"
	"slices"
	"strings"
	"strconv"
	"errors"
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
	ClassName			string
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
	_, added := addedApps[className]
	if !slices.Contains(pinnedApps, className) && !added {
		addItem(className)
		addIndicator(className, ipcClient)
	} else {
		addIndicator(className, ipcClient)
	}
}

func removeItem(className string) {
	thisApp := addedApps[className]

	thisApp.ButtonBox.Destroy()

	window.ShowAll()
}

func removeApp(address string) {
	thisApp, windowIndex, err := searhByAddress(address)
	if err != nil {fmt.Println(err)}
	className := thisApp.ClassName
	fmt.Println(className)

	listClients()

	if thisApp.Instances == 1 && !slices.Contains(pinnedApps, className) {
		removeItem(className)
		delete(addedApps, className)
		return
	}

	mainBox := thisApp.ButtonBox
	thisApp.IndicatorImage.Destroy()

	var newImage *gtk.Image
	if thisApp.Instances == 1 && slices.Contains(pinnedApps, className) {
		newImage = createImage(
			THEMES_DIR + config.CurrentTheme + "/empty.svg", config.IconSize - 10)
	}
	
	if thisApp.Instances == 2 {
		newImage = createImage(
			THEMES_DIR + config.CurrentTheme + "/single.svg", config.IconSize - 10)
	}

	if thisApp.Instances == 3 {
		newImage = createImage(
			THEMES_DIR + config.CurrentTheme + "/multiple.svg", config.IconSize - 10)
	}

	if thisApp.Instances > 3 {
		newImage = createImage(
			THEMES_DIR + config.CurrentTheme + "/3.svg", config.IconSize - 10)
	}

	mainBox.Add(newImage)

	addedApps[className].Instances -= 1
	fmt.Println(addedApps[className].Instances)
	newWindows := removeFromSlice(addedApps[className].Windows, windowIndex)
	addedApps[className].Windows = newWindows
	addedApps[className].IndicatorImage = newImage
	
	window.ShowAll()
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

	addedApps[className] = &appData{
		IndicatorImage: indicatorImage,
		Button: button,
		ButtonBox: item,
		DesktopData: itemProp,
		Instances: 0,
		ClassName: className,
	}

	button.Connect("clicked", func() {
		if addedApps[className].Instances == 0 {
			launch(itemProp.Exec)
		}
		if addedApps[className].Instances == 1 {
			hyprctl("dispatch focuswindow address:" + addedApps[className].Windows[0]["Address"])
		}
		if addedApps[className].Instances > 1 {
			fmt.Println("more")
		}
	})

	item.Add(button)
	item.Add(indicatorImage)

	itemsBox.Add(item)
	window.ShowAll()
}

func addIndicator(className string, ipcClient client) {
	thisApp := addedApps[className]
	mainBox := thisApp.ButtonBox

	appWindow := make(map[string]string)
	appWindow["Address"] = ipcClient.Address
	appWindow["Title"] = ipcClient.Title

	addedApps[className].Instances += 1 
	addedApps[className].Windows = append(
		addedApps[className].Windows, appWindow)


	fmt.Println(addedApps[className].Instances)
	thisApp.IndicatorImage.Destroy()

	var newImage *gtk.Image
	if thisApp.Instances == 1 {
		newImage = createImage(
			THEMES_DIR + config.CurrentTheme + "/single.svg", config.IconSize - 10)

	} else if thisApp.Instances == 2 {
		newImage = createImage(
			THEMES_DIR + config.CurrentTheme + "/multiple.svg", config.IconSize - 10)
	} else {
		newImage = createImage(
			THEMES_DIR + config.CurrentTheme + "/3.svg", config.IconSize - 10)
	}

	addedApps[className].IndicatorImage = newImage
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
	fmt.Println(source)
	pixbuf, err := iconTheme.LoadIcon(
		source, size, gtk.ICON_LOOKUP_FORCE_SIZE)
	if err != nil {
		fmt.Println(source, err)
		pixbuf, _ = iconTheme.LoadIcon(
			"cancel", size, gtk.ICON_LOOKUP_FORCE_SIZE)
	}

	image, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		fmt.Println(err)
	}

	return image
}

func removeFromSlice(slice []map[string]string, s int) []map[string]string {
    return append(slice[:s], slice[s+1:]...)
}

func searhByAddress(address string) (*appData, int, error) {
	for _, data := range addedApps {
		for windowIndex, appWindow := range data.Windows {
			if appWindow["Address"] == address {
				return data, windowIndex, nil
			}
		}
	}

	err := errors.New("Window not found: " + address)

	return nil, 0, err
}