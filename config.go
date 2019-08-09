package archive_builder_bot

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	Token string `json:"token"`
}

func GetConfig(cfgPath string) (*Config, error) {
	filename, err := filepath.Abs(cfgPath)
	if err != nil {
		return nil, err
	}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
