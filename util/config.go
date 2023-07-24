package util

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Cached   bool   `json:"cached"`
}

func CacheLoginDetails(config Config) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return errors.New("failed to get user config directory")
	}

	err = os.MkdirAll(configDir+"/limestone", os.ModePerm)
	if err != nil {
		return err
	}

	filePath := configDir + "/limestone/config.json"
	var dest *os.File

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		dest, err = os.Create(filePath)
		if err != nil {
			return err
		}

		defer dest.Close()
	} else {
		return err
	}

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

func (config *Config) GetLoginDetails() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	path := dir + "/limestone/config.json"
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	if config.Email != "" && config.Password != "" {
		config.Cached = true
	}

	return nil
}
