package cfg

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/akshaybharambe14/go-jsonc"
	"github.com/goccy/go-json"
	// "github.com/tidwall/sjson"
)

type Config struct {
	CurrentTheme string
	IconSize     int
	Layer        string
	Position     string
	Blur         string
	Spacing      int
	Margin       int
	Consts       map[string]string
}

func GetDefaultConfig() Config {
	return Config{
		CurrentTheme: "lotos",
		IconSize:     21,
		Layer:        "auto",
		Position:     "bottom",
		Blur:         "on",
		Spacing:      8,
		Margin:       8,
		Consts:       map[string]string{},
	}
}

func ConnectConfig(jsoncFile string, isTheme bool) Config {
	// Read jsonc
	file, err := os.Open(jsoncFile)
	if err != nil {
		log.Println("Config file not found!\n", err, "\nLoad default config")
		return GetDefaultConfig()
	}
	defer file.Close()

	decoder := jsonc.NewDecoder(file)
	res, _ := io.ReadAll(decoder)
	log.Println(json.Valid(res), jsoncFile)

	config := Config{}
	if err = json.Unmarshal(res, &config); err != nil {
		log.Println("Config is incorrect!\n", err, "\nLoad default config")
		return GetDefaultConfig()
	}

	// Set default values ​​if not specified
	if config.CurrentTheme == "" && !isTheme {
		config.CurrentTheme = GetDefaultConfig().CurrentTheme
		log.Println("The theme is not set, the default theme is currently used - \"lotos\"")
	}

	if config.Layer == "" {
		config.Layer = GetDefaultConfig().CurrentTheme
	}

	isCorrect := config.Position == "left" || config.Position == "right"
	isCorrect2 := config.Position == "top" || config.Position == "bottom"
	isCorrect3 := isCorrect || isCorrect2

	if config.Position == "" || !isCorrect3 {
		if !isTheme {
			log.Println("Position incorrect or empty\nDefault position set")
		}
		config.Position = GetDefaultConfig().Position
	}

	if config.Spacing == 0 {
		config.Spacing = GetDefaultConfig().Spacing
	}

	if config.IconSize == 0 {
		config.IconSize = GetDefaultConfig().IconSize
	}

	if config.Blur == "" {
		config.Blur = GetDefaultConfig().Blur
	}

	config.Consts = GetDefaultConfig().Consts

	return config
}

type ItemList struct {
	Pinned []string
}

func ReadItemList(jsonFile string) []string {
	file, _ := os.Open(jsonFile)
	defer file.Close()

	decoder := json.NewDecoder(file)

	itemList := ItemList{}
	err := decoder.Decode(&itemList)
	if err != nil {
		log.Println("error: ", err)
	}

	return itemList.Pinned
}

func ChangeJsonPinnedApps(apps []string, jsoncFile string) error {
	itemList := ItemList{
		Pinned: apps,
	}

	jsonData, err := json.MarshalIndent(itemList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	file, err := os.Create(jsoncFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
