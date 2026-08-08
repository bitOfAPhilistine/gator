package config

import (
	"os"
	"encoding/json"
)

configPath, _ := os.UserHomeDir() + "/.gatorconfig.json"

type Config struct {
	dbUrl string `json:"db_url"`
	currentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(name string) error {
	c.currentUserName = name
	
	return Write(c)
}

func Read() (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config
}

func Write(config *Config) error {
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	err = json.NewEncoder(file).Encode(config)
	if err != nil {
		return err
	}

	return nil
}