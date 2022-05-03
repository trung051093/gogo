package appctx

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Name     string `yaml:"name"`
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
		Port     int    `yaml:"port"`
		TimeZone string `yaml:"time_zone"`
		SSLMode  string `yaml:"ssl_mode"`
	} `yaml:"database"`
	JWT struct {
		Secret             string `yaml:"secret"`
		PasswordSaltLength int    `yaml:"pass_salt_length"`
		ExpireDays         int    `yaml:"expire_days"`
	} `yaml:"jwt"`
}

func GetConfig(cfg *Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		fmt.Sprintln(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Sprintln(err)
	}
}
