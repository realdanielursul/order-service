package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      `yaml:"app"`
		HTTP     `yaml:"http"`
		Postgres `yaml:"postgres"`
		Redis    `yaml:"redis"`
		Kafka    `yaml:"kafka"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}

	HTTP struct {
		Port string `yaml:"port"`
	}

	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		Database string `yaml:"database"`
		SSLMode  string `yaml:"ssl_mode"`
	}

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password" env:"REDIS_PASSWORD"`
		DB       int    `yaml:"db"`
	}

	Kafka struct {
		Host  string `yaml:"host"`
		Port  string `yaml:"port"`
		Topic string `yaml:"topic"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := cleanenv.UpdateEnv(cfg); err != nil {
		return nil, fmt.Errorf("error updating env file: %w", err)
	}

	return cfg, nil
}
