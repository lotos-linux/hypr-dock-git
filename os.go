package main

import (
	"os"
	"fmt"
	"strings"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/gdk"
)

type clientData struct {
	Name 	string
	Icon	string
	Exec	string
}

const appDir = "/usr/share/applications/"

func getClientData(className string) (clientData, error) {
	allData, err := loadTextFile(searchDesktopFile(className))
	if err != nil {
		return clientData{}, err
	}

	data := clientData{
		Name: getClientOption(allData, "Name"),
		Icon: getClientOption(allData, "Icon"),
		Exec: getClientOption(allData, "Exec"),
	}

	return data, nil
}

func searchDesktopFile(className string) string{
	desktopFile := className + ".desktop"
	_, err := os.Stat(appDir + desktopFile)
	if err == nil {
		return appDir + desktopFile
	}

	files, _ := os.ReadDir(appDir)
	for _, file := range files {
		fileName := file.Name()
		
		// "krita" > "org.kde.krita.desktop" / "lutris" > "net.lutris.Lutris.desktop"
		if strings.Count(fileName, ".") > 1 && strings.Contains(fileName, className) {
			return appDir + fileName
		}
		// "VirtualBox Manager" > "virtualbox.desktop"
		if fileName == strings.Split(strings.ToLower(className), " ")[0] + ".desktop" {
			return appDir + fileName
		}
	}

	return ""
}

func getClientOption(allData []string, option string) string {
	for lineIndex := range len(allData) {
		line := allData[lineIndex]
		if strings.Contains(line, option + "=") {
			optionValue := strings.TrimPrefix(line, option + "=")
			return optionValue
		}
	}
	return ""
}

func loadTextFile(path string) ([]string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(bytes), "\n")
	var output []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			output = append(output, line)
		}

	}
	return output, nil
}

func createImage(source string) *gtk.Image {
	iconTheme, err := gtk.IconThemeGetDefault()
	if err != nil {
		fmt.Println("Unable to icon theme:", err)
	}

	// Create image in file
	if strings.Contains(source, "/") {
		pixbuf, err := gdk.PixbufNewFromFileAtSize(
			source, config.IconSize, config.IconSize)
		if err != nil {
			fmt.Println(err)
			pixbuf, _ = iconTheme.LoadIcon(
				"steam", config.IconSize, gtk.ICON_LOOKUP_FORCE_SIZE)
			
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