/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
