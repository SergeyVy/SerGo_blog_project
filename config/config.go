package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Env         string `yaml:"env"`
	DatabaseURL string `yaml:"database_url"`
	JWTSecret   string `yaml:"jwt_secret"`
}

func MustLoad() *Config {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config/config.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("не удалось прочитать конфиг: %v ", err)
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("не удалось распарсить конфиг: %v", err)
	}

	// ⬇️ Переопределение из переменных окружения, если есть
	if envJWT := os.Getenv("JWT_SECRET"); envJWT != "" {
		cfg.JWTSecret = envJWT
	}

	return &cfg
}
