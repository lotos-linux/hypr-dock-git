package cfg

import (
	"hypr-dock/internal/pkg/utils"
	"io"
	"log"
	"os"
	"slices"

	"github.com/akshaybharambe14/go-jsonc"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	// "github.com/tidwall/sjson"
)

type Config struct {
	CurrentTheme  string
	IconSize      int
	Layer         string
	Position      string
	Blur          string
	Spacing       int
	SystemGapUsed string
	Margin        int
}

type ThemeConfig struct {
	Blur    string
	Spacing int
}

type ItemList struct {
	Pinned []string
}

func GetDefaultConfig() Config {
	return Config{
		CurrentTheme:  "lotos",
		IconSize:      21,
		Layer:         "auto",
		Position:      "bottom",
		Blur:          "true",
		Spacing:       8,
		SystemGapUsed: "true",
		Margin:        8,
	}
}

func ReadConfig(jsoncFile string) Config {
	// Read jsonc
	config := Config{}
	err := ReadJsonc(jsoncFile, &config)
	if err != nil {
		log.Println(err)
		log.Println("Load default config")
		return GetDefaultConfig()
	}

	// Set default values ​​if not specified
	if config.CurrentTheme == "" {
		config.CurrentTheme = GetDefaultConfig().CurrentTheme
		log.Println("The theme is not set, the default theme is currently used - \"lotos\"")
	}

	correctLayers := []string{"auto", "background", "bottom", "top", "overlay"}
	if !slices.Contains(correctLayers, config.Layer) {
		log.Println("Layer", "\""+config.Layer+"\"", "is incorrect or empty. Default layer set")
		config.Layer = GetDefaultConfig().CurrentTheme
	}

	correctPositions := []string{"left", "right", "top", "bottom"}
	if !slices.Contains(correctPositions, config.Position) {
		log.Println("Position", "\""+config.Layer+"\"", "is incorrect or empty. Default position set")
		config.Position = GetDefaultConfig().Position
	}

	correctBlurModes := []string{"true", "false"}
	if !slices.Contains(correctBlurModes, config.Blur) {
		config.Blur = GetDefaultConfig().Blur
	}

	correctSystemGapUsed := []string{"true", "false"}
	if !slices.Contains(correctSystemGapUsed, config.SystemGapUsed) {
		log.Println("SystemGapUsed", "\""+config.SystemGapUsed+"\"", "is incorrect or empty. Defailt value set")
		config.SystemGapUsed = GetDefaultConfig().SystemGapUsed
	}

	if config.Spacing < 1 {
		config.Spacing = GetDefaultConfig().Spacing
	}

	if config.IconSize < 1 {
		config.IconSize = GetDefaultConfig().IconSize
	}

	return config
}

func ReadTheme(jsoncFile string, config Config) *ThemeConfig {
	// Read jsonc
	themeConfig := ThemeConfig{}
	err := ReadJsonc(jsoncFile, &themeConfig)
	if err != nil {
		log.Println(err)
		log.Println("Load default config")
		return nil
	}

	// Set default values ​​if not specified
	correctBlurModes := []string{"true", "false"}
	if !slices.Contains(correctBlurModes, themeConfig.Blur) {
		log.Println("Blur", "\""+themeConfig.Blur+"\"", "is incorrect or empty. Default blur set")
		themeConfig.Blur = config.Blur
	}

	if themeConfig.Spacing < 0 {
		themeConfig.Spacing = config.Spacing
	}

	return &themeConfig
}

func ReadItemList(jsonFile string) []string {
	itemList := ItemList{}

	if !utils.FileExists(jsonFile) {
		itemList.Pinned = CreateEmptyPinnedFile(jsonFile)
		return itemList.Pinned
	}

	err := ReadJson(jsonFile, &itemList)
	if err != nil {
		log.Fatal(err)
	}

	return itemList.Pinned
}

func ReadJsonc(jsoncFile string, v interface{}) error {
	file, err := os.Open(jsoncFile)
	if err != nil {
		return errors.Wrapf(err, "file %q not found", jsoncFile)
	}
	defer file.Close()

	decoder := jsonc.NewDecoder(file)
	res, err := io.ReadAll(decoder)
	if err != nil {
		return errors.Wrapf(err, "file %q. io.ReadAll error", jsoncFile)
	}

	if err = json.Unmarshal(res, &v); err != nil {
		return errors.Wrapf(err, "file %q has a syntax error", jsoncFile)
	}

	return nil
}

func ChangeJsonPinnedApps(apps []string, jsonFile string) error {
	itemList := ItemList{
		Pinned: apps,
	}

	if err := WriteItemList(jsonFile, itemList); err != nil {
		log.Println("Error", jsonFile, "writing: ", err)
		return err
	}

	return nil
}

func ReadJson(jsonFile string, v interface{}) error {
	file, err := os.Open(jsonFile)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&v); err != nil {
		return err
	}

	return nil
}

func CreateEmptyPinnedFile(jsonFile string) []string {
	initialData := ItemList{
		Pinned: []string{},
	}

	if err := WriteItemList(jsonFile, initialData); err != nil {
		log.Fatalf("Failed to create file %q: %v", jsonFile, err)
		return nil
	}

	return initialData.Pinned
}

func WriteItemList(jsonFile string, data ItemList) error {
	file, err := os.Create(jsonFile)
	if err != nil {
		return errors.Wrapf(err, "failed to create file %q", jsonFile)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return errors.Wrapf(err, "failed to encode data to file %q", jsonFile)
	}

	return nil
}
