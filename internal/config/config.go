package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Auth     Authentication `json:"user_data"`
}

type Authentication struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

func (config *Config) CacheLoginDetails() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	dest, err := os.Create(path + "config.json")
	if err != nil {
		return err
	}

	defer dest.Close()

	encoded, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = dest.Write(encoded)
	if err != nil {
		return err
	}

	err = dest.Sync()

	return err
}

func GetLoginDetails() (Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	content, err := os.ReadFile(path + "config.json")
	if err != nil {
		return Config{}, err
	}

	var config Config
	
	err = json.Unmarshal(content, &config)
	if err != nil {
		return Config{}, err
	}

	if config.Email == "" || config.Password == "" {
		return Config{}, errors.New("email or password is empty")
	}

	return config, nil
}

func RemoveConfigDetails() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	return os.Remove(path + "config.json")
}

func getConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	path := dir + "/limestone/"

	return path, err
}
