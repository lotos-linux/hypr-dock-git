package h

import (
	"log"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func CreateImage(source string, size int) (*gtk.Image, error) {
	// Create image in file
	if strings.Contains(source, "/") {
		pixbuf, err := gdk.PixbufNewFromFileAtSize(source, size, size)
		if err != nil {
			log.Println(err)
			return ReturnMissingIcon(size), err
		}

		return CreateImageFromPixbuf(pixbuf), nil
	}

	// Create image in icon name
	iconTheme, err := gtk.IconThemeGetDefault()
	if err != nil {
		log.Println("Unable to icon theme:", err)
		return ReturnMissingIcon(size), err
	}

	pixbuf, err := iconTheme.LoadIcon(source, size, gtk.ICON_LOOKUP_FORCE_SIZE)
	if err != nil {
		log.Println(source, err)
		return ReturnMissingIcon(size), err
	}

	return CreateImageFromPixbuf(pixbuf), nil
}

func CreateImageFromPixbuf(pixbuf *gdk.Pixbuf) *gtk.Image {
	image, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		log.Println("Error creating image from pixbuf:", err)
		return nil
	}
	return image
}

func ReturnMissingIcon(size int) *gtk.Image {
	var iconPath string

	icon := "/icon-missing.svg"
	iconFromConfigDir := config.Consts["CONFIG_DIR"] + icon
	iconFromThemeDir := config.Consts["THEMES_DIR"] + icon

	if FileExists(iconFromThemeDir) {
		iconPath = iconFromThemeDir
	} else if FileExists(iconFromConfigDir) {
		iconPath = iconFromConfigDir
	}

	if iconPath == "" {
		log.Println("Unable to icon:", iconFromThemeDir)
		log.Println("Unable to icon:", iconFromConfigDir)
		return nil
	}

	pixbuf, err := gdk.PixbufNewFromFileAtSize(iconPath, size, size)
	if err != nil {
		log.Println(err)
		return nil
	}

	return CreateImageFromPixbuf(pixbuf)
}
