package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Admin  []string `json:"admin"`
	Status string   `json:"status"`
}

func GetConfig() Config {

	// Open our jsonFile
	fileName := "config.json"
	file, err := os.Open(fileName) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully Opened ", fileName)

	bb, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	json.Unmarshal(bb, &config)

	// log.Println(config)

	return config
}
