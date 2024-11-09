package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

var config Config

type Config struct {
	ApiUrl   string `yaml:"apiUrl"`
	ApiToken string `yaml:"apiToken"`
	Port     string `yaml:"httpPort"`
}

func ParseConfigFile(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	return config, err
}

func LoadConfig() (Config, error) {
	ConfigFile := flag.String("c", "./Config.yaml", "Config path")
	flag.Parse()
	Config, err := ParseConfigFile(*ConfigFile)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error loading Config: %v\n", err)
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error: %s", err)
			if err != nil {
				return Config, err
			}
		}
	}
	return Config, err
}
