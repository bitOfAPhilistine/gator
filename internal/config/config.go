package config

import (
	"os"
	"encoding/json"
)

var homeDir, _ = os.UserHomeDir()
var configPath = homeDir + "/.gatorconfig.json"

type Config struct {
	DbUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	
	return Write(c)
}

func Read() (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Write(config *Config) error {
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func Initialize(connectionURL string) error {
	initConfig := Config{
		DbUrl: connectionURL,
		CurrentUserName: "",
	}

	return Write(&initConfig)
}