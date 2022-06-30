package appctx

import (
	"fmt"
	"gogo/common"
	jaegerprovider "gogo/components/jaeger"
	graylog "gogo/components/log"
	"gogo/components/mailer"
	rabbitmqprovider "gogo/components/rabbitmq"
	storageprovider "gogo/components/storage"
	"os"
	"path"
	"runtime"

	es "github.com/elastic/go-elasticsearch/v8"
	redis "github.com/go-redis/redis/v8"
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
	Redis struct {
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		Port     int    `yaml:"port"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
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
	Minio struct {
		Host            string `yaml:"host"`
		Port            int    `yaml:"port"`
		AccessKeyID     string `yaml:"access_key_id"`
		SecretAccessKey string `yaml:"secret_access_key"`
		UseSSL          bool   `yaml:"useSSL"`
		PublicUrl       string `yaml:"public_url"`
	} `yaml:"minio"`
	JWT struct {
		Secret             string `yaml:"secret"`
		PasswordSaltLength int    `yaml:"pass_salt_length"`
		ExpireDays         int    `yaml:"expire_days"`
	} `yaml:"jwt"`
	GoogleOauth2 struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		RedirectUri  string `yaml:"redirect_uri"`
	} `yaml:"google_oauth2"`
	Jaeger struct {
		ServiceName       string `yaml:"service_name"`
		AgentEndpoint     string `yaml:"agent_endpoint"`
		CollectorEndpoint string `yaml:"collector_endpoint"`
	} `yaml:"jaeger"`
	Graylog struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		Key  string `yaml:"key"`
	} `yaml:"graylog"`
	Swagger struct {
		Title       string   `yaml:"title"`
		Description string   `yaml:"description"`
		Version     string   `yaml:"version"`
		Host        string   `yaml:"host"`
		Schemes     []string `yaml:"schemes"`
		BasePath    string   `yaml:"basepath"`
	} `yaml:"swagger"`
	Mail struct {
		Sender   string `yaml:"sender"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"mail"`
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

func (cfg *Config) GetRedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}
}

func (cfg *Config) GetStorageConfig() *storageprovider.StorageConfig {
	return &storageprovider.StorageConfig{
		Endpoint:        fmt.Sprintf("%s:%d", cfg.Minio.Host, cfg.Minio.Port),
		AccessKeyID:     cfg.Minio.AccessKeyID,
		SecretAccessKey: cfg.Minio.SecretAccessKey,
		UseSSL:          cfg.Minio.UseSSL,
		PublicUrl:       cfg.Minio.PublicUrl,
	}
}

func (cfg *Config) GetJaegerConfig() *jaegerprovider.JaegerConfig {
	return &jaegerprovider.JaegerConfig{
		ServiceName:       cfg.Jaeger.ServiceName,
		AgentEndpoint:     cfg.Jaeger.AgentEndpoint,
		CollectorEndpoint: cfg.Jaeger.CollectorEndpoint,
	}
}

func (cfg *Config) GetGraylogConfig() *graylog.GraylogConfig {
	return &graylog.GraylogConfig{
		Host: cfg.Graylog.Host,
		Port: cfg.Graylog.Port,
		Key:  cfg.Graylog.Key,
	}
}

func (cfg *Config) GetSwaggerConfig() *common.SwaggerInfo {
	return &common.SwaggerInfo{
		Host:        cfg.Swagger.Host,
		Title:       cfg.Swagger.Title,
		Description: cfg.Swagger.Description,
		Version:     cfg.Swagger.Version,
		Schemes:     cfg.Swagger.Schemes,
		BasePath:    cfg.Swagger.BasePath,
	}
}

func (cfg *Config) GetMailConfig() *mailer.MailConfig {
	return &mailer.MailConfig{
		Sender:   cfg.Mail.Sender,
		Host:     cfg.Mail.Host,
		Port:     cfg.Mail.Port,
		Username: cfg.Mail.Username,
		Password: cfg.Mail.Password,
	}
}
