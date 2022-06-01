package appctx

import (
	"fmt"
	"os"
	"path"
	"runtime"
	rabbitmqprovider "user_management/components/rabbitmq"

	es "github.com/elastic/go-elasticsearch/v8"
	"gopkg.in/yaml.v3"
)

type Environment string

const (
	EnvLocal Environment = "local"
	EnvDev   Environment = "dev"
	EnvProd  Environment = "prod"
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
	RabbitMQ struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"rabbitmq"`
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

func GetFileConfig() string {
	environment := os.Getenv("env")
	rootDir := RootDir()
	if environment == "" {
		environment = string(EnvLocal)
	}
	file := fmt.Sprintf("%s/config_%s.yml", rootDir, environment)
	fmt.Println("File Config:", file)
	return file
}

func GetConfig() *Config {
	var cfg Config
	f, err := os.Open(GetFileConfig())
	if err != nil {
		fmt.Sprintln(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Sprintln(err)
	}
	return &cfg
}

func (cfg *Config) GetRabbitMQConfig() *rabbitmqprovider.RabbitmqConfig {
	return &rabbitmqprovider.RabbitmqConfig{
		Host: cfg.RabbitMQ.Host,
		Port: cfg.RabbitMQ.Port,
		User: cfg.RabbitMQ.Username,
		Pass: cfg.RabbitMQ.Password,
	}
}

func (cfg *Config) GetElasticSearchConfig() *es.Config {
	return &es.Config{
		Addresses: []string{cfg.ElasticSearch.Host},
		Username:  cfg.ElasticSearch.Username,
		Password:  cfg.ElasticSearch.Password,
	}
}
