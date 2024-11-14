package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"log/slog"
	"os"
)

var (
	config Config
)

type Config struct {
	ApiUrl   string `yaml:"apiUrl"`
	ApiToken string `yaml:"apiToken"`
	Address  string `yaml:"httpAddress"`
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
	logger := SetupLogger(config.Env)
	logger.Debug("Configuration file parsed")
	return config, err
}

func LoadConfig() (Config, error) {
	file := flag.String("c", "./config.yaml", "Config path")
	flag.Parse()
	Config, err := ParseConfigFile(*file)
	if err != nil {
		logger := SetupLogger(config.Env)
		logger.Error("Configuration file not parsed")
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
