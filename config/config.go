package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)



type BotConfig struct {
	Token string  `envconfig:"BOT_TOKEN" required:"true"`
	ApiID int `envconfig:"API_ID" required:"true"`
	ApiHash string `envconfig:"API_HASH" required:"true"`
}

type DbConfig struct {
	DbUrl string `envconfig:"DATABASE_URL" required:"true"`
}

type LogConfig struct {
	Level string `envconfig:"LOG_LEVEL" required:"true" default:"info"`
	Filename string `envconfig:"LOG_FILENAME" required:"true"`
	MaxSize int `envconfig:"LOG_MAX_SIZE" required:"true"`
	MaxBackups int `envconfig:"LOG_MAX_BACKUPS" required:"true"`
	MaxAge int `envconfig:"LOG_MAX_AGE" required:"true"`
	Compress bool `envconfig:"LOG_COMPRESS" required:"true"`
}


type Config struct {
	BotConfig
	DbConfig
	LogConfig
}

func LoadConfig() (*Config, error) {
 _ = godotenv.Load()

 var cfg Config
 if err := envconfig.Process("", &cfg); err != nil {
	return nil, err
 }

 return &cfg, nil

}

