package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	Token                string   `yaml:"token"`
	ProxyAddr            string   `yaml:"proxy-addr"`
	ProxyUser            string   `yaml:"proxy-user"`
	ProxyPassword        string   `yaml:"proxy-password"`
	CloudJenkinsUrl      string   `yaml:"cloud-jenkins-url"`
	CloudJenkinsUser     string   `yaml:"cloud-jenkins-user"`
	CloudJenkinsPassword string   `yaml:"cloud-jenkins-password"`
	ChmodJenkinsUrl      string   `yaml:"chmod-jenkins-url"`
	ChmodJenkinsUser     string   `yaml:"chmod-jenkins-user"`
	ChmodJenkinsPassword string   `yaml:"chmod-jenkins-password"`
	TgUsers              []string `yaml:"tg-users"`
	GitRepoCloneDir      string   `yaml:"git_repo_clone_dir"`
	GitRepoUrl           string   `yaml:"git_repo_url"`
	GitUser              string   `yaml:"git_user"`
	GitEmail             string   `yaml:"git_email"`
	GitPassword          string   `yaml:"git_password"`
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
