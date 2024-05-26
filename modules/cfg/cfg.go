package cfg

import (
	"os"
	"fmt"
	"io/ioutil"
	"github.com/goccy/go-json"
	"github.com/akshaybharambe14/go-jsonc"
	// "github.com/tidwall/sjson"
)

var err error

type Config struct {
	CurrentTheme 	string
	IconSize   		int
	Layer     		string
	Position		string
	Blur			string
	Spacing			int
	Margin			int
}

func GetDefaultConfig() Config {
	config := Config{}

	config.CurrentTheme = "lotos"
	config.IconSize = 21
	config.Layer = "auto"
	config.Position = "bottom"
	config.Blur = "on"
	config.Spacing = 8
	config.Margin = 8

	return config
}

func ConnectConfig(jsoncFile string, isTheme bool) Config {
	// Read jsonc
	file, err := os.Open(jsoncFile)
	if err != nil {
		fmt.Println("Config file not found!\n", err, "\nLoad default config")
		return GetDefaultConfig()
	}
	defer file.Close()

	decoder := jsonc.NewDecoder(file)
	res, err := ioutil.ReadAll(decoder)
	fmt.Println(json.Valid(res), jsoncFile)

	config := Config{}
	if err = json.Unmarshal(res, &config); err != nil {
		fmt.Println("Config is incorrect!\n", err, "\nLoad default config")
		return GetDefaultConfig()
	}


	// Set default values ​​if not specified
	if config.CurrentTheme == "" && !isTheme{
		config.CurrentTheme = GetDefaultConfig().CurrentTheme
		fmt.Println("The theme is not set, the default theme is currently used - \"lotos\"")
	}

	if config.Layer == "" {
		config.Layer = GetDefaultConfig().CurrentTheme
	}

	isCorrect := config.Position == "left" || config.Position == "right"
	isCorrect2 := config.Position == "top" || config.Position == "bottom"
	isCorrect3 := isCorrect || isCorrect2

	if config.Position == "" || !isCorrect3 {
		if !isTheme {
			fmt.Println("Position incorrect or empty\nDefault position set")
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

	return config
}


type ItemList struct {
	List			[]map[string]string
}

func ReadItemList(jsonFile string) ItemList {
	file, _ := os.Open(jsonFile)
	defer file.Close()

	decoder := json.NewDecoder(file)

	itemList := ItemList{}
	err := decoder.Decode(&itemList)
	if err != nil {
	  fmt.Println("error: ", err)
	}
	
	return itemList
}