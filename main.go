package main

import (
	"fmt"
	"flag"
	"os"
	"syscall"
	"time"
	"hypr-dock/cfg"
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/allan-simon/go-singleinstance"
	"strconv"
)

const version = "0.0.5-0-alpha"

// Only during development
const CONFIG_DIR = "/home/lots/hypr-dock/configs"
// const CONFIG_DIR = "./configs"
const THEMES_DIR = CONFIG_DIR + "/themes/"
const MAIN_CONFIG = CONFIG_DIR + "/config.jsonc"
const ITEMS_CONFIG = CONFIG_DIR + "/pinned.json"

var err error
var config cfg.Config
var window *gtk.Window
var detectArea *gtk.Window

var orientation gtk.Orientation

func initSettings() {
	configFile := flag.String("config", MAIN_CONFIG, "config file")

	config = cfg.ConnectConfig(*configFile, false)
	pinnedApps = cfg.ReadItemList(ITEMS_CONFIG)

	currentTheme := flag.String("theme", config.CurrentTheme, "theme")
	config.CurrentTheme = *currentTheme

	themeConfig := cfg.ConnectConfig(
		THEMES_DIR + config.CurrentTheme + "/" + config.CurrentTheme + ".jsonc", true)

	config.Blur = themeConfig.Blur
	config.Spacing = themeConfig.Spacing

	flag.Parse()
}

func main() {
	signalHandler()

	lockFilePath := fmt.Sprintf("%s/hypr-dock-%s.lock", tempDir(), os.Getenv("USER"))
	lockFile, err := singleinstance.CreateLockFile(lockFilePath)
	if err != nil {
		file, err := loadTextFile(lockFilePath)
		if err == nil {
			pidStr := file[0]
			pidInt, _ := strconv.Atoi(pidStr)
			syscall.Kill(pidInt, syscall.SIGUSR1)
		}
		os.Exit(0)
	}
	defer lockFile.Close()

	// Window build
	initSettings()
	gtk.Init(nil)

	window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		fmt.Println("Unable to create window:", err)
	}

	window.SetTitle("hypr-dock")
	orientation := setWindowProperty(window)

	err = addCssProvider(THEMES_DIR + config.CurrentTheme + "/style.css")	 
	if err != nil {
		fmt.Println(
			"CSS file not found, the default GTK theme is running!\n", err)
	}

	buildApp(orientation)

	window.Add(app)
	window.Connect("destroy", func() {gtk.MainQuit()})
	window.ShowAll()

	// Build detect area
	if config.Layer == "auto" {initDetectArea()}

	// Hyprland socket connect
	go initHyprEvents()


	gtk.Main()
}

func initDetectArea() {
	detectArea, _ = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	detectArea.SetName("detect")

	layershell.InitForWindow(detectArea)
	layershell.SetNamespace(detectArea, "dock-detect")
	layershell.SetAnchor(detectArea, Edge, true)
	layershell.SetMargin(detectArea, Edge, 0)
	layershell.SetLayer(detectArea, layershell.LAYER_SHELL_LAYER_TOP)

	switch orientation {
	case gtk.ORIENTATION_HORIZONTAL:
		detectArea.SetSizeRequest(config.IconSize * len(addedApps) * 2 - 20, 1)
	case gtk.ORIENTATION_VERTICAL:
		detectArea.SetSizeRequest(1, config.IconSize * len(addedApps) * 2 - 20)
	}

	detectArea.Connect("enter-notify-event", func(window *gtk.Window, e *gdk.Event) {
		event := gdk.EventCrossingNewFromEvent(e)
		isInWindow := event.Detail() == 3 || event.Detail() == 4 || true

		isCancelHide = 1
		if isInWindow {
			go func() {
				setLayer("top")
			}()
		}
	})

	detectArea.ShowAll()
}

func addCssProvider(cssFile string) error {
	cssProvider, _ := gtk.CssProviderNew()
	err := cssProvider.LoadFromPath(cssFile)

	if err == nil {
		screen, _ := gdk.ScreenGetDefault()

		gtk.AddProviderForScreen(
			screen, cssProvider, 
			gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

		return nil
	}

	return err
}

var Edge = layershell.LAYER_SHELL_EDGE_BOTTOM

func setWindowProperty(window *gtk.Window) gtk.Orientation {
	AppOreintation := gtk.ORIENTATION_HORIZONTAL
	Layer := layershell.LAYER_SHELL_LAYER_BOTTOM
	Edge = layershell.LAYER_SHELL_EDGE_BOTTOM

	switch config.Layer {
	case "background":
		Layer = layershell.LAYER_SHELL_LAYER_BACKGROUND
	case "bottom":
		Layer = layershell.LAYER_SHELL_LAYER_BOTTOM
	case "top":
		Layer = layershell.LAYER_SHELL_LAYER_TOP
	case "overlay":
		Layer = layershell.LAYER_SHELL_LAYER_OVERLAY
	}

	switch config.Position {
	case "left":
		Edge = layershell.LAYER_SHELL_EDGE_LEFT
		AppOreintation = gtk.ORIENTATION_VERTICAL
	case "bottom":
		Edge = layershell.LAYER_SHELL_EDGE_BOTTOM
	case "right":
		Edge = layershell.LAYER_SHELL_EDGE_RIGHT
		AppOreintation = gtk.ORIENTATION_VERTICAL
	case "top":
		Edge = layershell.LAYER_SHELL_EDGE_TOP
	}

	layershell.InitForWindow(window)
	layershell.SetNamespace(window, "hypr-dock")
	layershell.SetAnchor(window, Edge, true)
	layershell.SetMargin(window, Edge, 0)

	addLayerRule()

	if config.Layer == "auto" {
		layershell.SetLayer(window, Layer)
		autoLayer()
		return AppOreintation
	}

	layershell.SetLayer(window, Layer)
	return AppOreintation
}

func autoLayer() {
	window.Connect("enter-notify-event", func(window *gtk.Window, e *gdk.Event) {
		event := gdk.EventCrossingNewFromEvent(e)
		isInWindow := event.Detail() == 3 || event.Detail() == 4 || true

		isCancelHide = 1
		if isInWindow && !special {
			go func() {
				setLayer("top")
			}()
		}
	})

	window.Connect("leave-notify-event", func(window *gtk.Window, e *gdk.Event) {
		event := gdk.EventCrossingNewFromEvent(e)
		isInWindow := event.Detail() == 3 || event.Detail() == 4
		isCancelHide = 0

		if isInWindow {
			go func() {
				time.Sleep(time.Second / 3) 
				setLayer("bottom")
			}()
		}
	})
}

func setLayer(layer string) {
	switch layer {
	case "top":
		layershell.SetLayer(window, layershell.LAYER_SHELL_LAYER_TOP)
	case "bottom":
		if isCancelHide == 0 {
			layershell.SetLayer(window, layershell.LAYER_SHELL_LAYER_BOTTOM)
		}
		isCancelHide = 0
	}
}
