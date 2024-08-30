package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Maxmind struct {
		Key  string `json:"key"`
		ASN  string
		City string
	} `json:"maxmind"`
	Nameserver string `json:"nameserver"`
	Host       string `json:"host"`
	Port       string `json:"port"`
}

// Load configuration file
func LoadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		return config, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config, err
}
