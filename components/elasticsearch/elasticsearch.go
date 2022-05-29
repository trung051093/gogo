package elasticsearch

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v7"
)

type key string

type ElasticSearchSevice struct {
	client *elasticsearch.Client
}

var ElasticSearchServiceKey key = "ElasticSearchService"
var once sync.Once
var instance *ElasticSearchSevice

// singleton
func NewEsService(config elasticsearch.Config) (*ElasticSearchSevice, error) {
	client, err := elasticsearch.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return &ElasticSearchSevice{client: client}, nil
}

func GetIntance(config elasticsearch.Config) *ElasticSearchSevice {
	once.Do(func() {
		service, err := NewEsService(config)
		if err != nil {
			return
		}
		instance = service
	})
	return instance
}

func WithContext(ctx context.Context, es *ElasticSearchSevice) context.Context {
	return context.WithValue(ctx, ElasticSearchServiceKey, es)
}

func FromContext(ctx context.Context) (*ElasticSearchSevice, bool) {
	esService := ctx.Value(ElasticSearchServiceKey)
	if es, ok := esService.(*ElasticSearchSevice); ok {
		return es, true
	}
	return nil, false
}

func (es *ElasticSearchSevice) LogInfo(ctx context.Context) {
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
