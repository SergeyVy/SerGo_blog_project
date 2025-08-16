package config

import (
	"gopkg.in/yaml.v3"
	_ "gopkg.in/yaml.v3"
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
		log.Fatalf("не удалосьраспарсить конфиг: %v", err)
	}
	return &cfg
}
