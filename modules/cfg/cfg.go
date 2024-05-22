package cfg

import (
	"os"
	"fmt"
	"io/ioutil"
	"github.com/goccy/go-json"
	"github.com/akshaybharambe14/go-jsonc"
)

type Config struct {
	CurrentTheme	string
	IconSize   		int
	Layer     		string
	Position		string
	Margin			int
}

func Connect(jsonFile string) Config {
	file, _ := os.Open(jsonFile)
	defer file.Close()

	decoder := jsonc.NewDecoder(file)
	res, err := ioutil.ReadAll(decoder)

	config := Config{}
	if err = json.Unmarshal(res, &config); err != nil {
		fmt.Println("error while json Unmarshal: ", err)
	}

	return config
}