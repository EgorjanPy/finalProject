package config

import (
	"flag"
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
}

func MustLoad() *Config {
	config_path := flag.String("CONFIG_PATH", "./config/config.yaml", "")
	flag.Parse()
	if *config_path == "" {
		// log.Fatalf("Config path is empty")
	}
	if _, err := os.Stat(*config_path); os.IsNotExist(err) {
		// log.Fatalf("Config path is incorrect: %s", *config_path)
	}
	var cfg Config
	err := cleanenv.ReadConfig(*config_path, &cfg)
	if err != nil {
		// log.Fatal("Cant read config file")
	}
	return &cfg
}
