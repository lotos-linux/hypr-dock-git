package h

import (
	"os"
	"path/filepath"
)

const (
	ConfigDir   = ".config/hypr-dock"
	ThemesDir   = "themes"
	MainConfig  = "config.jsonc"
	ItemsConfig = "pinned.json"
)

func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Не удалось определить домашнюю директорию: " + err.Error())
	}
	return filepath.Join(homeDir, ConfigDir)
}

func GetThemesPath() string {
	return filepath.Join(GetConfigPath(), ThemesDir)
}

func GetMainConfigPath() string {
	return filepath.Join(GetConfigPath(), MainConfig)
}

func GetItemsConfigPath() string {
	return filepath.Join(GetConfigPath(), ItemsConfig)
}
