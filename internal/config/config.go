package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Meteo struct {
	BaseUrl   string        `yaml:"base_url"`
	Timeout   time.Duration `yaml:"timeout"`
	Latitude  string        `yaml:"latitude"`
	Longitude string        `yaml:"longitude"`
}

type HFace struct {
	Timeout time.Duration `yaml:"timeout"`
	BaseUrl string        `yaml:"base_url"`
}

type Config struct {
	TelegramApiToken string `env:"TELEGRAM_API_TOKEN"`
	HfaceApiToken    string `env:"HFACE_API_TOKEN"`

	DbPath   string `yaml:"database_path"`
	CronSync string `yaml:"cron"`

	HFace `yaml:"hface"`
	Meteo `yaml:"meteo"`
}

func LoadConfigFromEnv() *Config {
	configPath := "./config/default.yml"

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(fmt.Errorf("failed to load YAML config: %w", err))
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(fmt.Errorf("failed to load environment variables: %w", err))
	}

	return &cfg
}
