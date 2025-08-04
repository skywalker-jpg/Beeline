package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Logger struct {
	Sink  string `yaml:"sink"`
	Level string `yaml:"level"`
}

type Server struct {
	URL       string `yaml:"url"`
	AuthToken string `yaml:"auth_token"`
	ServerURL string `yaml:"server_url"`
}

type AppConfig struct {
	Logger Logger `yaml:"logger"`
	Server Server `yaml:"server"`
}

func NewConfig(path string) (*AppConfig, error) {
	yamlConfig, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var appConfig AppConfig

	if err := yaml.Unmarshal(yamlConfig, &appConfig); err != nil {
		return nil, err
	}
	return &appConfig, nil
}
