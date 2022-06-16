package jaegerprovider

import (
	"log"

	"contrib.go.opencensus.io/exporter/jaeger"
)

type JaegerService interface {
	GetExporter() *jaeger.Exporter
}

type JaegerConfig struct {
	AgentEndpoint     string
	CollectorEndpoint string
	ServiceName       string
}

type jaegerService struct {
	exporter *jaeger.Exporter
}

func NewExporter(config *JaegerConfig) *jaegerService {
	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     config.AgentEndpoint,
		CollectorEndpoint: config.CollectorEndpoint,
		Process:           jaeger.Process{ServiceName: config.ServiceName},
	})
	if err != nil {
		log.Fatalf("Failed to create the Jaeger exporter: %v", err)
	}

	return &jaegerService{exporter: je}
}

func (j *jaegerService) GetExporter() *jaeger.Exporter {
	return j.exporter
}
