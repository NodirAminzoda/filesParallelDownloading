package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  Server  `yaml:"server"`
	Adapter Adapter `yaml:"adapter"`
}

type Server struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	ReadTimeOut  time.Duration `yaml:"read_time_out"`
	WriteTimeOut time.Duration `yaml:"write_time_out"`
}

type Adapter struct {
	TimeOut time.Duration `yaml:"time_out"`
}

var path = "./internal/config/config.yml"

func MustReadConfig(log *zerolog.Logger) *Config {
	body, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config file")
	}

	cfg := &Config{}
	err = yaml.Unmarshal(body, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config file")
	}

	log.Info().Msg("Config file successfully loaded")

	return cfg
}
