package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host              string            `yaml:"host"`
	Port              string            `yaml:"port"`
	ListenPort        string            `yaml:"listenport"`
	Debug             bool              `yaml:"debug"`
	NumAnalizeWorkers int               `yaml:"numanalizeworkers`
	Scriptfile        map[string]string `yaml:"scriptfile"`
}

var (
	config Config
)

func Get(configname *string) *Config {
	yamlFile, err := ioutil.ReadFile(*configname)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	fmt.Println("Use " + *configname)
	return &config
}
