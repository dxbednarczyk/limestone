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

func CacheLoginDetails(config Config) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	err = os.MkdirAll(configDir+"/limestone", os.ModePerm)
	if err != nil {
		return err
	}

	filePath := configDir + "/limestone/config.json"
	var dest *os.File

	dest, err = os.Create(filePath)
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
	if err != nil {
		return err
	}

	return nil
}

func GetLoginDetails() (Config, error) {
	var config Config

	dir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, err
	}

	path := dir + "/limestone/config.json"
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

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
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	path := dir + "/limestone/config.json"

	return os.Remove(path)
}
