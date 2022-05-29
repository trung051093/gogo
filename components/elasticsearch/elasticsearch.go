package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type key string

type ElasticSearchSevice struct {
	client *elasticsearch.Client
}

var ElasticSearchServiceKey key = "ElasticSearchService"
var once sync.Once
var instance *ElasticSearchSevice
var instanceErr error

func NewEsService(config elasticsearch.Config) (*ElasticSearchSevice, error) {
	client, err := elasticsearch.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return &ElasticSearchSevice{client: client}, nil
}

// singleton
func GetIntance(config elasticsearch.Config) (*ElasticSearchSevice, error) {
	once.Do(func() {
		service, instanceErr := NewEsService(config)
		if instanceErr != nil {
			return
		}
		instanceErr = nil
		instance = service
	})
	return instance, instanceErr
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

func (es *ElasticSearchSevice) GetClient(ctx context.Context) *elasticsearch.Client {
	return es.client
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

func (es *ElasticSearchSevice) Index(ctx context.Context, index string, id string, data []byte) {
	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      index,
		Body:       bytes.NewReader(data),
		DocumentID: id,
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(ctx, es.client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%s", res.Status(), id)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
}

// delete document by ID
func (es *ElasticSearchSevice) Delete(ctx context.Context, index string, id string) {
	// Set up the request object.
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
		Refresh:    "true",
	}
	// Perform the request with the client.
	res, err := req.Do(ctx, es.client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Printf("[%s] Error delete document ID=%s", res.Status(), id)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
}
