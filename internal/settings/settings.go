package settings

import (
	"errors"
	"hypr-dock/internal/pkg/cfg"
	"hypr-dock/internal/pkg/flags"
	"hypr-dock/internal/pkg/utils"
	"hypr-dock/pkg/ipc"
	"log"
	"os"
	"path/filepath"
)

const RunMode = "normal"

// const RunMode = "dev"
const DefaultTheme = "lotos"

var ConfigDir string
var ConfigPath string
var PinnedPath string
var ThemesDir string
var CurrentThemeDir string
var CurrentThemeConfigPath string
var CurrentThemeStylePath string

var conf cfg.Config
var PinnedApps []string

func Get() cfg.Config {
	return conf
}

func setConfigDir(mode string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Home dir: " + err.Error())
	}

	runModes := map[string]func() string{}
	runModes["normal"] = func() string {
		return filepath.Join(homeDir, ".config/hypr-dock")
	}
	runModes["dev"] = func() string {
		return filepath.Join(homeDir, "repos/hypr-dock/configs")
	}

	return runModes[mode]()
}

func Init() error {
	flags := flags.Get(DefaultTheme)

	ConfigDir = setConfigDir(RunMode)

	PinnedPath = filepath.Join(ConfigDir, "pinned.json")
	PinnedApps = cfg.ReadItemList(PinnedPath)
	defaultConfigPath := filepath.Join(ConfigDir, "config.jsonc")

	if flags.Config == "~/.config/hypr-dock" {
		ConfigPath = defaultConfigPath
	} else {
		ConfigPath = flags.Config
	}

	conf = cfg.ReadConfig(ConfigPath)

	ThemesDir = filepath.Join(ConfigDir, "themes")
	CurrentThemeDir = filepath.Join(ThemesDir, conf.CurrentTheme)

	if !utils.FileExists(CurrentThemeDir) {
		log.Println("Current theme not found (", conf.CurrentTheme, "). Loading default theme")

		if conf.CurrentTheme == DefaultTheme {
			log.Println("Default theme not found")
			return errors.New("default theme not found")
		}

		conf.CurrentTheme = DefaultTheme
	}

	CurrentThemeStylePath = filepath.Join(CurrentThemeDir, "style.css")
	CurrentThemeConfigPath = filepath.Join(CurrentThemeDir, conf.CurrentTheme+".jsonc")

	themeConfig := cfg.ReadTheme(CurrentThemeConfigPath, conf)
	if themeConfig == nil {
		log.Println(CurrentThemeConfigPath, "not found. Load default values")
		return nil
	}

	conf.Blur = themeConfig.Blur
	conf.Spacing = themeConfig.Spacing

	if conf.Blur == "true" {
		enableBlur()

		ipc.AddEventListener("configreloaded", func(event string) {
			go enableBlur()
		}, true)
	}

	return nil
}

func enableBlur() {
	ipc.Hyprctl("keyword layerrule blur,hypr-dock")
	ipc.Hyprctl("keyword layerrule ignorealpha 0.4,hypr-dock")
}
