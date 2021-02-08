package common

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Test string `json:"test"`
	Jira struct {
		Host    	string `json:"host"`
		Login 		string `json:"login"`
		Token 		string `json:"token"`
	} `json:"jira"`
}

func LoadConfiguration(file string) Config {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	config := Config{}
	jsonParser.Decode(&config)

	return config
}