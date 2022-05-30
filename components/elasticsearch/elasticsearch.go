package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	elasticsearchmodel "user_management/components/elasticsearch/model"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type key string

type ElasticSearchSevice interface {
	LogInfo(ctx context.Context)
	GetClient(ctx context.Context) *elasticsearch.Client
	Index(ctx context.Context, index string, id string, data []byte)
	Delete(ctx context.Context, index string, id string)
	Search(ctx context.Context, index string, query string) (*elasticsearchmodel.SearchResults, error)
}

type elasticSearchSevice struct {
	client *elasticsearch.Client
}

var ElasticSearchServiceKey key = "ElasticSearchService"
var once sync.Once
var instance *elasticSearchSevice
var instanceErr error

func NewEsService(config elasticsearch.Config) (*elasticSearchSevice, error) {
	client, err := elasticsearch.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return &elasticSearchSevice{client: client}, nil
}

// singleton
func GetIntance(config elasticsearch.Config) (*elasticSearchSevice, error) {
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

func WithContext(ctx context.Context, es ElasticSearchSevice) context.Context {
	return context.WithValue(ctx, ElasticSearchServiceKey, es)
}

func FromContext(ctx context.Context) (*elasticSearchSevice, bool) {
	esService := ctx.Value(ElasticSearchServiceKey)
	if es, ok := esService.(*elasticSearchSevice); ok {
		return es, true
	}
	return nil, false
}

func (es *elasticSearchSevice) GetClient(ctx context.Context) *elasticsearch.Client {
	return es.client
}

func (es *elasticSearchSevice) LogInfo(ctx context.Context) {
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

func (es *elasticSearchSevice) Index(ctx context.Context, index string, id string, data []byte) {
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
func (es *elasticSearchSevice) Delete(ctx context.Context, index string, id string) {
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

// Search returns results matching a query, paginated by after.
func (es *elasticSearchSevice) Search(ctx context.Context, index string, query string) (*elasticsearchmodel.SearchResults, error) {
	var results elasticsearchmodel.SearchResults

	res, err := es.client.Search(
		es.client.Search.WithIndex(index),
		es.client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return &results, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return &results, err
		}
		return &results, fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}

	type envelopeResponse struct {
		Took int
		Hits struct {
			Total struct {
				Value int
			}
			Hits []struct {
				ID         string          `json:"_id"`
				Source     json.RawMessage `json:"_source"`
				Highlights json.RawMessage `json:"highlight"`
				Sort       []interface{}   `json:"sort"`
			}
		}
	}

	var r envelopeResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return &results, err
	}

	results.Total = r.Hits.Total.Value

	if len(r.Hits.Hits) < 1 {
		results.Hits = []map[string]interface{}{}
		return &results, nil
	}

	for _, hit := range r.Hits.Hits {
		var data interface{}
		var highlights interface{}
		h := map[string]interface{}{
			"sort": hit.Sort,
		}

		if err := json.Unmarshal(hit.Source, &data); err != nil {
			return &results, err
		}

		if len(hit.Highlights) > 0 {
			if err := json.Unmarshal(hit.Highlights, &highlights); err != nil {
				return &results, err
			}
		}

		h["highlights"] = highlights
		h["data"] = data
		results.Hits = append(results.Hits, h)
	}

	return &results, nil
}

func BuildQuery(ctx context.Context, query string, after ...string) string {
	var b strings.Builder

	b.WriteString("{\n")
	b.WriteString(query)

	if len(after) > 0 && after[0] != "" && after[0] != "null" {
		b.WriteString(",\n")
		b.WriteString(fmt.Sprintf(`	"search_after": %s`, after))
	}

	b.WriteString("\n}")

	fmt.Printf("%s\n", b.String())
	return b.String()
}
