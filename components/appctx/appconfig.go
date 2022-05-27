package appctx

import (
	"fmt"
	"os"
	"path"
	"runtime"

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
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Port     int    `yaml:"port"`
		TimeZone string `yaml:"time_zone"`
		SSLMode  string `yaml:"ssl_mode"`
	} `yaml:"database"`
	ElasticSearch struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"elasticsearch"`
	JWT struct {
		Secret             string `yaml:"secret"`
		PasswordSaltLength int    `yaml:"pass_salt_length"`
		ExpireDays         int    `yaml:"expire_days"`
	} `yaml:"jwt"`
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return d
}

func GetConfig(cfg *Config) {
	rootDir := RootDir()
	fmt.Sprintln("root dir:", rootDir)
	f, err := os.Open(rootDir + "/config.yml")
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
