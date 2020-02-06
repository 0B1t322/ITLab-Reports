package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	DB *DBConfig	`json:"DbOptions"`
	Auth *AuthConfig	`json:"AuthOptions"`
	App *AppConfig		`json:"AppOptions"`
}

type DBConfig struct {
	Host 		string		`json:"host"`
	Port 		string		`json:"port"`
	DBName 		string		`json:"dbname"`
	CollectionName 	string 	`json:"collectionName"`
}
type AuthConfig struct {
	KeyURL		string		`json:"keyUrl"`
	Audience	string		`json:"audience"`
	TestKeyURL	string		`json:"testKeyUrl"`
}
type AppConfig struct {
	TestMode	bool	`json:"testMode"`
}

func GetConfig() *Config {
	var config Config
	data, err := ioutil.ReadFile("src/ITLabReports/api/config.json")
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Panic(err)
	}
	return &config
}
