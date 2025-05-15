package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TimeAddMs      time.Duration `yaml:"TIME_ADDITION_MS" env-default:"100ms"`
	TimeSubMs      time.Duration `yaml:"TIME_SUBTRACTION_MS" env-default:"100ms"`
	TimeMulMs      time.Duration `yaml:"TIME_MULTIPLICATIONS_MS" env-default:"100ms"`
	TimeDivMs      time.Duration `yaml:"TIME_DIVISIONS_MS" env-default:"100ms"`
	Port           string        `yaml:"PORT" env-default:"8080"`
	ComputingPower int           `yaml:"COMPUTING_POWER" env-default:"2"`
	StoragePath    string        `yaml:"storage_path"`
	Secret         string        `yaml:"jwt_secret"`
}

func MustLoad() *Config {
	os.Setenv("CONFIG_PATH", "./config/config.yaml")
	config_path := os.Getenv("CONFIG_PATH")
	// fmt.Printf("PATH=%s", config_path)
	fmt.Println()
	if config_path == "" {
		config_path = "finalProject/config/config.yaml"
	}
	if _, err := os.Stat(config_path); os.IsNotExist(err) {
		// log.Fatalf("Config path is incorrect: %s", *config_path)
	}
	var cfg Config
	err := cleanenv.ReadConfig(config_path, &cfg)
	if err != nil {
		// log.Fatal("Cant read config file")
	}
	return &cfg
}
