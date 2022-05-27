package main

import (
	"user_management/components/appctx"
	"user_management/components/elasticsearch"

	es "github.com/elastic/go-elasticsearch/v7"
)

func main() {
	config := &appctx.Config{}
	appctx.GetConfig(config)

	configEs := &es.Config{
		Addresses: []string{config.ElasticSearch.Host},
		Username:  config.ElasticSearch.Username,
		Password:  config.ElasticSearch.Password,
	}
	esService := elasticsearch.NewEsService(*configEs)
	esService.LogInfo()
}
