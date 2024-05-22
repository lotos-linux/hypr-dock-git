package cfg

import (
	"os"
    "fmt"
    "github.com/goccy/go-json"
	// "github.com/dlasky/gotk3-layershell/layershell"
)

type Config struct {
	CurrentTheme  	string
    IconSize   		int
    Layer     		string
	Position		string
	Margin			int
}

func Connect(jsonFile string) Config {
	file, _ := os.Open(jsonFile)
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
	  fmt.Println("error:", err)
	}

	return config
}