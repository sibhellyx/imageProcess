package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port          string `mapstructure:"SERV_PORT"`
	InputDir      string `mapstructure:"STORAGE_INPUTDIR"`
	OutputDir     string `mapstructure:"STORAGE_OUTPUTDIR"`
	FileSize      int    `mapstructure:"FILE_SIZE"`
	QueueCapacity int    `mapstructure:"QUEUE_CAPACITY"`
	Workers       int    `mapstructure:"WORKERS"`
}

func LoadConfig() Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	v.SetDefault("SERV_PORT", "8080")
	v.SetDefault("STORAGE_INPUTDIR", "./storage/input")
	v.SetDefault("STORAGE_OUTPUTDIR", "./storage/output")
	v.SetDefault("FILE_SIZE", 5)
	v.SetDefault("QUEUE_CAPACITY", 10)
	v.SetDefault("WORKERS", 1)

	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults and environment variables")
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	} else {
		log.Printf("Using config file: %s", v.ConfigFileUsed())
	}

	var cfg Config

	err = v.Unmarshal(&cfg)
	if err != nil {
		log.Panic(err)
	}
	return cfg
}
