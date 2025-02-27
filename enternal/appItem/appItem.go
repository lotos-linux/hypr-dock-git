package appItem

import (
	"log"
	"slices"

	"github.com/gotk3/gotk3/gtk"

	"hypr-dock/enternal/pkg/cfg"
	"hypr-dock/enternal/pkg/desktop"
	"hypr-dock/enternal/pkg/h"
	"hypr-dock/enternal/pkg/ipc"
)

var config cfg.Config

type Item struct {
	Instances      int
	Windows        []map[string]string
	DesktopData    *desktop.Desktop
	ClassName      string
	Button         *gtk.Button
	ButtonBox      *gtk.Box
	IndicatorImage *gtk.Image
	List           map[string]*Item
	PinnedList     *[]string
}

func New(className string) (*Item, error) {
	desktopData := desktop.New(className)

	item, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	button, err := gtk.ButtonNew()
	if err == nil {
		image, err := h.CreateImage(desktopData.Icon, config.IconSize)
		if err == nil {
			button.SetImage(image)
		} else {
			log.Println(err)
		}

		button.SetName(className)
		button.SetTooltipText(desktopData.Name)

		item.Add(button)
	} else {
		log.Println(err)
	}

	indicatorImage, err := GetIndicatorImage(0)
	if err == nil {
		item.Add(indicatorImage)
	} else {
		log.Println(err)
	}

	return &Item{
		IndicatorImage: indicatorImage,
		Button:         button,
		ButtonBox:      item,
		DesktopData:    desktopData,
		Instances:      0,
		ClassName:      className,
		List:           nil,
		PinnedList:     nil,
	}, nil
}

func (item *Item) RemoveLastInstance(windowIndex int) {
	item.IndicatorImage.Destroy()

	newImage, err := GetIndicatorImage(item.Instances - 1)
	if err == nil {
		item.ButtonBox.Add(newImage)
	}

	item.Instances -= 1
	item.Windows = h.RemoveFromSlice(item.Windows, windowIndex)
	item.IndicatorImage = newImage
}

func (item *Item) UpdateState(ipcClient ipc.Client) {
	appWindow := map[string]string{
		"Address": ipcClient.Address,
		"Title":   ipcClient.Title,
	}

	if item.IndicatorImage != nil {
		item.IndicatorImage.Destroy()
	}

	indicatorImage, err := GetIndicatorImage(item.Instances + 1)
	if err == nil {
		item.ButtonBox.Add(indicatorImage)
	}

	item.Windows = append(item.Windows, appWindow)
	item.IndicatorImage = indicatorImage
	item.Instances += 1
}

func (item *Item) IsPinned() bool {
	return slices.Contains(*item.PinnedList, item.ClassName)
}

func (item *Item) TogglePin() {
	if item.IsPinned() {
		h.RemoveFromSliceByValue(&*item.PinnedList, item.ClassName)
		if item.Instances == 0 {
			item.ButtonBox.Destroy()
			delete(item.List, item.ClassName)
		}
		log.Println("Remove:", item.ClassName)
	} else {
		h.AddToSlice(&*item.PinnedList, item.ClassName)
		log.Println("Add:", item.ClassName)
	}

	err := cfg.ChangeJsonPinnedApps(*item.PinnedList, config.Consts["ITEMS_CONFIG"])
	if err != nil {
		log.Println("Error: ", err)
	} else {
		log.Println("File", config.Consts["ITEMS_CONFIG"], "saved successfully!", item.ClassName)
	}
}

func GetIndicatorImage(instances int) (*gtk.Image, error) {
	var path string
	themeDir := config.Consts["THEMES_DIR"] + config.CurrentTheme + "/"

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

	return h.CreateImage(path, config.IconSize-10)
}

func SetConfig(inpConfig cfg.Config) {
	config = inpConfig
}
