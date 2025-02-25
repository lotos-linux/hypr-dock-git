package main

import (
	"os"
	"fmt"
	"sync"
	"slices"
	"strings"
	"strconv"
	"errors"
	"hypr-dock/cfg"
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

func (this *appData) IsPinned() bool {
	for _, app := range pinnedApps {
		if app == this.ClassName {
			return true
		}
	}
	return false
}

func (this *appData) TogglePin() {
	changeType := ""

	if this.IsPinned() {
		changeType = "Remove"
		removeFromSliceByValue(&pinnedApps, this.ClassName)
		if this.Instances == 0 {
			removeItem(this.ClassName)
			delete(addedApps, this.ClassName)
		}
	} else {
		changeType = "Add"
		addToSlice(&pinnedApps, this.ClassName)
	}

	err := cfg.ChangeJsonPinnedApps(pinnedApps, ITEMS_CONFIG)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("File", ITEMS_CONFIG, "saved successfully! |",
					changeType+":", this.ClassName)
	}
}

var (addedApps = make(map[string]*appData)
	 mu        sync.Mutex)








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

	button.Connect("button-release-event", func(button *gtk.Button, e *gdk.Event) {
		event := gdk.EventButtonNewFromEvent(e)
		if event.Button() == 3 {
			menu := contextMenu(addedApps[className], true)
			menu.PopupAtWidget(button, gdk.GDK_GRAVITY_NORTH, gdk.GDK_GRAVITY_SOUTH, nil)
			return
		}
		if addedApps[className].Instances == 0 {
			launch(itemProp.Exec)
		}
		if addedApps[className].Instances == 1 {
			hyprctl("dispatch focuswindow address:" + addedApps[className].Windows[0]["Address"])
		}
		if addedApps[className].Instances > 1 {
			menu := contextMenu(addedApps[className], false)
			menu.PopupAtWidget(button, gdk.GDK_GRAVITY_NORTH, gdk.GDK_GRAVITY_SOUTH, nil)
		}
	})

	item.Add(button)
	item.Add(indicatorImage)

	itemsBox.Add(item)
	window.ShowAll()
}

func addIndicator(className string, ipcClient client) {
	thisApp := addedApps[className]

	appWindow := make(map[string]string)
	appWindow["Address"] = ipcClient.Address
	appWindow["Title"] = ipcClient.Title

	addedApps[className].Instances += 1 
	addedApps[className].Windows = append(
		addedApps[className].Windows, appWindow)

	thisApp.IndicatorImage.Destroy()

	indicatorPath := getIndicatorPath(thisApp.Instances)
	newImage := createImage(indicatorPath, config.IconSize - 10)

	addedApps[className].IndicatorImage = newImage
	thisApp.ButtonBox.Add(newImage)
	window.ShowAll()
}

func getIndicatorPath(instances int) string {
	var path string
	themeDir := THEMES_DIR + config.CurrentTheme + "/"

	switch {
	case instances == 0:
		path = themeDir + "empty.svg"
	case instances == 1:
		path = themeDir + "single.svg"
	case instances == 2:
		path = themeDir + "multiple.svg"
	case instances > 2:
		path = themeDir + "3.svg"
	}

	return path
}

func removeApp(address string) {
	thisApp, windowIndex, err := searhByAddress(address)
	if err != nil {
		fmt.Println(err)
		return
	}
	className := thisApp.ClassName

	listClients()

	if thisApp.Instances == 1 && !slices.Contains(pinnedApps, className) {
		removeItem(className)
		delete(addedApps, className)
		return
	}

	mainBox := thisApp.ButtonBox
	thisApp.IndicatorImage.Destroy()

	indicatorPath := getIndicatorPath(thisApp.Instances-1)
	newImage := createImage(indicatorPath, config.IconSize - 10)

	mainBox.Add(newImage)

	addedApps[className].Instances -= 1
	newWindows := removeFromSlice(addedApps[className].Windows, windowIndex)
	addedApps[className].Windows = newWindows
	addedApps[className].IndicatorImage = newImage
	
	window.ShowAll()
}

func removeItem(className string) {
	thisApp := addedApps[className]

	thisApp.ButtonBox.Destroy()

	window.ShowAll()
}

func contextMenu(app *appData, isContext bool) *gtk.Menu {
	var windowMenuItems []*gtk.MenuItem

	windows := app.Windows
	menu, err := gtk.MenuNew()
	if err != nil {
		fmt.Println(err)
	}

	for i, window := range windows {
		menuItem, err := buildContextItem(window["Title"], func() {
			go hyprctl("dispatch focuswindow address:" + window["Address"])
		})
		if err != nil {
			fmt.Println(err)

			if i == len(windows)-1 && isContext && i > 0 {
				windowMenuItems[i-1].SetMarginBottom(10)
			}

			continue
		}

		if i == len(windows)-1 && isContext{
			menuItem.SetMarginBottom(10)
		}

		windowMenuItems = append(windowMenuItems, menuItem)
		menu.Append(menuItem)
	}

	if isContext {
		labelText := "Pin"
		if app.IsPinned() {
			labelText = "Unpin"
		}

		menuItem, err := buildContextItem(labelText, func() {
			app.TogglePin()
		})
		if err != nil {
			fmt.Println(err)
		}

		menu.Append(menuItem)
	}

	menu.SetName("context-menu")
	menu.ShowAll()

	return menu
}

func buildContextItem(labelText string, connectFunc func()) (*gtk.MenuItem, error) {
	menuItem, err := gtk.MenuItemNew()
	if err != nil {
		return nil, err
	}

	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		return nil, err
	}

	label, err := gtk.LabelNew(labelText)
	if err != nil {
		return nil, err
	}

	if connectFunc != nil {
		menuItem.Connect("activate", connectFunc)
	}

	hbox.Add(label)
	menuItem.Add(hbox)

	return menuItem, nil
}

func cancelHide(button *gtk.Button) {
	button.Connect("enter-notify-event", func() {
		isCancelHide = 1
	})
}

func createImage(source string, size int) *gtk.Image {
	// Create image in file
	if strings.Contains(source, "/") {
		pixbuf, err := gdk.PixbufNewFromFileAtSize(source, size, size)
		if err != nil {
			fmt.Println(err)
			return returnMissingIcon(size)
		}

		return CreateImageFromPixbuf(pixbuf)
	}

	// Create image in icon name
	iconTheme, err := gtk.IconThemeGetDefault()
	if err != nil {
		fmt.Println("Unable to icon theme:", err)
		return returnMissingIcon(size)
	}

	pixbuf, err := iconTheme.LoadIcon(source, size, gtk.ICON_LOOKUP_FORCE_SIZE)
	if err != nil {
		fmt.Println(source, err)
		return returnMissingIcon(size)
	}

	return CreateImageFromPixbuf(pixbuf)
}

func CreateImageFromPixbuf(pixbuf *gdk.Pixbuf) *gtk.Image {
    image, err := gtk.ImageNewFromPixbuf(pixbuf)
    if err != nil {
        fmt.Println("Error creating image from pixbuf:", err)
        return nil
    }
    return image
}

func returnMissingIcon(size int) *gtk.Image {
	var iconPath string

	icon              := "icon-missing.svg"
	iconFromConfigDir := CONFIG_DIR + "/" + icon
	iconFromThemeDir  := CONFIG_DIR + THEMES_DIR + icon

	if fileExists(iconFromThemeDir) {
		iconPath = iconFromThemeDir
	} else if fileExists(iconFromConfigDir) {
		iconPath = iconFromConfigDir
	} 

	if iconPath == "" {
		fmt.Println("Unable to icon:", iconFromThemeDir)
		fmt.Println("Unable to icon:", iconFromConfigDir)
		return nil
	}

	pixbuf, err := gdk.PixbufNewFromFileAtSize(iconPath, size, size)	
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return CreateImageFromPixbuf(pixbuf)
}

func fileExists(path string) bool {
    _, err := os.Stat(path)
    return !errors.Is(err, os.ErrNotExist)
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

func removeFromSliceByValue(slice *[]string, value string) {
	index := -1
	for i, v := range *slice {
		if v == value {
			index = i
			break
		}
	}

	if index != -1 {
		*slice = append((*slice)[:index], (*slice)[index+1:]...)
	}
}

func addToSlice(slice *[]string, value string) {
	*slice = append(*slice, value)
}