package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

// LoadConfig from file
func LoadConfig(configFile string) (cc CloudCredentials, err error) {
	if configFile == "" {
		return cc, errors.New("Must provide a config file")
	}

	file, err := os.Open(configFile)
	if err != nil {
		return cc, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return cc, err
	}

	err = json.Unmarshal(bytes, &cc)
	if err != nil {
		return cc, err
	}

	err = cc.Validate()
	if err != nil {
		return cc, err
	}

	return cc, nil
}
