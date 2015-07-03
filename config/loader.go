package config

import (
	"encoding/json"
	"io/ioutil"
)

type PloymentConfig struct {
	Repositoryurl   string `json:"repositoryUrl"`
	Targetdirectory string `json:"targetDirectory"`
}

func FromFile(path string) (PloymentConfig, error) {

	var settings PloymentConfig
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return settings, err
	}

	if err = json.Unmarshal(configFile, &settings); err != nil {
		return settings, err
	}

	return settings, nil
}
