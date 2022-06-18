package graylog

import (
	"fmt"
	"log"
	"os"

	grayloghook "github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/sirupsen/logrus"
)

type GraylogConfig struct {
	Host string
	Port int
	Key  string
}

func Integrate(config *GraylogConfig) {
	hostname, _ := os.Hostname()
	logger := logrus.New()

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	hook := grayloghook.NewGraylogHook(addr, map[string]interface{}{
		"_X-STREAM-ID": config.Key,
		"_host":        hostname,
	})

	logrus.AddHook(hook)
	logger.AddHook(hook)
	logger.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(logger.Writer())
}
