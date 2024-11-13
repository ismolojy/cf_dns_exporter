package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log/slog"
	"os"
)

var config Config

type Config struct {
	ApiUrl   string `yaml:"apiUrl"`
	ApiToken string `yaml:"apiToken"`
	Port     string `yaml:"httpPort"`
	LogPath  string `yaml:"logPath"`
	Env      string `yaml:"env"`
}

func ParseConfigFile(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	return config, err
}

func LoadConfig(configFile *string) (Config, error) {
	Config, err := ParseConfigFile(*configFile)
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

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
