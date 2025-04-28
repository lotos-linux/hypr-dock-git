package item

import (
	"log"

	"github.com/gotk3/gotk3/gtk"

	"hypr-dock/internal/pkg/desktop"
	"hypr-dock/internal/pkg/utils"
	"hypr-dock/internal/settings"
	"hypr-dock/pkg/ipc"
)

func (item *Item) WindowsMenu() (*gtk.Menu, error) {
	menu, err := gtk.MenuNew()
	if err != nil {
		log.Println(err)
	}

	desktopData := desktop.New(item.ClassName)

	AddWindowsItemToMenu(menu, item.Windows, desktopData)

	menu.SetName("windows-menu")
	menu.ShowAll()

	return menu, nil
}

func (item *Item) ContextMenu(settings settings.Settings) (*gtk.Menu, error) {
	menu, err := gtk.MenuNew()
	if err != nil {
		log.Println(err)
	}

	desktopData := desktop.New(item.ClassName)

	AddWindowsItemToMenu(menu, item.Windows, desktopData)

	if item.Instances != 0 {
		separator, err := gtk.SeparatorMenuItemNew()
		if err == nil {
			menu.Append(separator)
		} else {
			log.Println(err)
		}
	}

	launchMenuItem, err := BuildLaunchMenuItem(item, desktopData.Exec)
	if err == nil {
		menu.Append(launchMenuItem)
	} else {
		log.Println(err)
	}

	pinMenuItem, err := BuildPinMenuItem(item, settings)
	if err == nil {
		menu.Append(pinMenuItem)
	} else {
		log.Println(err)
	}

	menu.SetName("context-menu")
	menu.ShowAll()

	return menu, nil
}

func AddWindowsItemToMenu(menu *gtk.Menu, windows []map[string]string, desktopData *desktop.Desktop) {
	for _, window := range windows {
		menuItem, err := BuildContextItem(window["Title"], func() {
			go ipc.Hyprctl("dispatch focuswindow address:" + window["Address"])
		}, desktopData.Icon)

		if err != nil {
			log.Println(err)
			continue
		}

		menu.Append(menuItem)
	}
}

func BuildLaunchMenuItem(item *Item, exec string) (*gtk.MenuItem, error) {
	labelText := "New window"
	if item.Instances == 0 {
		labelText = "Open"
	}

	launchMenuItem, err := BuildContextItem(labelText, func() {
		utils.Launch(exec)
	})

	if err != nil {
		return nil, err
	}

	return launchMenuItem, nil
}

func BuildPinMenuItem(item *Item, settings settings.Settings) (*gtk.MenuItem, error) {
	labelText := "Pin"
	if item.IsPinned() {
		labelText = "Unpin"
	}

	menuItem, err := BuildContextItem(labelText, func() {
		item.TogglePin(settings)
	})

	if err != nil {
		return nil, err
	}

	return menuItem, nil
}

func BuildContextItem(labelText string, connectFunc func(), iconName ...string) (*gtk.MenuItem, error) {
	menuItem, err := gtk.MenuItemNew()
	if err != nil {
		return nil, err
	}

	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		return nil, err
	}

	if len(iconName) > 0 {
		icon, err := utils.CreateImage(iconName[0], 16)
		if err == nil {
			hbox.Add(icon)
		}
	}

	label, err := gtk.LabelNew(labelText)
	if err != nil {
		return nil, err
	}

	if connectFunc != nil {
		menuItem.Connect("activate", func() {
			// dispather()
			connectFunc()
		})
	}

	hbox.Add(label)
	menuItem.Add(hbox)

	return menuItem, nil
}
