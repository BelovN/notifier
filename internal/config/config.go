package config

import "os"

type Config struct {
	TelegramApiToken string
	DbPath           string
	HfaceApiToken    string
}

func LoadConfigFromEnv() *Config {
	return &Config{
		TelegramApiToken: os.Getenv("TELEGRAM_API_TOKEN"),
		DbPath:           os.Getenv("DATABASE_PATH"),
		HfaceApiToken:    os.Getenv("HFACE_API_TOKEN"),
	}
}
