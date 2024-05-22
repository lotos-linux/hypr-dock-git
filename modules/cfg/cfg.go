package cfg

import (
	"os"
	"fmt"
	"io/ioutil"
	"github.com/goccy/go-json"
	"github.com/akshaybharambe14/go-jsonc"
	// "github.com/tidwall/sjson"
)

type Config struct {
	CurrentTheme 	string
	IconSize   		int
	Layer     		string
	Position		string
	Margin			int
	Pinned			[]map[string]string
}

func ConnectConfig(jsoncFile string) Config {
	file, _ := os.Open(jsoncFile)
	defer file.Close()

	decoder := jsonc.NewDecoder(file)
	res, err := ioutil.ReadAll(decoder)

	config := Config{}
	if err = json.Unmarshal(res, &config); err != nil {
		fmt.Println("error while json Unmarshal: ", err)
	}

	return config
}


type ItemList struct {
	ItemList		[]map[string]string
}

func ReadItemList(jsonFile string) ItemList {
	file, _ := os.Open(jsonFile)
	defer file.Close()

	decoder := json.NewDecoder(file)

	itemList := ItemList{}
	err := decoder.Decode(&itemList)
	if err != nil {
	  fmt.Println("error:", err)
	}
	return itemList
}