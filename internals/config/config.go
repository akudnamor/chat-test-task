package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	StoragePath string `yaml:"storage_path" env-required:"true"`
	Address     string `yaml:"address" env-default:"localhost:8000"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file is not exist")
	}

	var cfg Config

	// fill cfg struct
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("config file cant be read")
	}

	return &cfg
}
