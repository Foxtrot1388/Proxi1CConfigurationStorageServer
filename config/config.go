package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	ListenPort string `yaml:"listenport"`
	Debug      bool   `yaml:"debug"`
}

var (
	config Config
)

func Get(configname *string) *Config {
	yamlFile, err := ioutil.ReadFile("./" + *configname)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	return &config
}
