package middleware

import (
	"user_management/common"
	"user_management/components/appctx"
	"user_management/components/elasticsearch"

	esv7 "github.com/elastic/go-elasticsearch/v7"

	"github.com/gin-gonic/gin"
)

func AppendElasticSearch(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		config := &appctx.Config{}
		appctx.GetConfig(config)

		configEs := &esv7.Config{
			Addresses: []string{config.ElasticSearch.Host},
			Username:  config.ElasticSearch.Username,
			Password:  config.ElasticSearch.Password,
		}
		esService := elasticsearch.NewEsService(*configEs)
		ginCtx.Set(common.ElasticSearchService, esService)

	}
}
