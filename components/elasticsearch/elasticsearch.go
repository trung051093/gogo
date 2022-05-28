package elasticsearch

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

type ElasticSearchService interface {
	LogInfo()
}

type elasticSearchSevice struct {
	client *elasticsearch.Client
}

func NewEsService(config elasticsearch.Config) *elasticSearchSevice {
	client, err := elasticsearch.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return &elasticSearchSevice{
		client: client,
	}
}

func (es *elasticSearchSevice) LogInfo() {
	var (
		r map[string]interface{}
	)
	// 1. Get cluster info
	res, err := es.client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))
}
