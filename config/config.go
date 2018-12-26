package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	BaseURL    string `json:"base_url"`
	SigningKey []byte `json:"signing_key"`
	Directory  string `json:"directory"`
	HMRC       struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		APIBase      string `json:"api_base"`
	} `json:"hmrc"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
