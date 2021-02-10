package common

import (
	"encoding/json"
	"fmt"
	"os"
)

/**
	Configuration object
 */
type Config struct {
	Test string `json:"test"`
	Jira struct {
		Host    	string `json:"host"`
		Login 		string `json:"login"`
		Token 		string `json:"token"`
	} `json:"jira"`
}

type Names struct{
	NameList [] string
}

/**
	Load json configuration json file and return a Config structure
 */
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

/**
	Check if a name exists in the list. If not, append it.
 */
func (n * Names) AddName(name string) {
	for _, n := range n.NameList {
		if name == n {
			return
		}
	}
	n.NameList = append(n.NameList, name)
	return
}