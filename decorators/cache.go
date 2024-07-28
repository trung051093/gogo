package decorator

import (
	"bytes"
	"context"
	"gogo/common"
	"gogo/components/appctx"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v9"
	"github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

type dataCache struct {
	Code    int
	Headers http.Header
	Data    string
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type ApiHandler func(appctx.AppContext) func(*gin.Context)

func CacheRequest(
	appCtx appctx.AppContext,
	cacheName string,
	tls time.Duration,
) func(ApiHandler) func(*gin.Context) {
	cacheService := appCtx.GetCacheService()

	return func(apiHandler ApiHandler) func(*gin.Context) {
		return func(ginCtx *gin.Context) {
			ctx := context.TODO()
			url := ginCtx.Request.URL
			key := url.RequestURI()
			var data dataCache

			if err := cacheService.Get(ctx, key, &data); err != nil {
				blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ginCtx.Writer}
				ginCtx.Writer = blw
				apiHandler(appCtx)(ginCtx)
				go func(bodyWriter *bodyLogWriter) {
					if bodyWriter.Status() > 300 || bodyWriter.Status() < 200 {
						return
					}
					err := cacheService.Once(&cache.Item{
						Ctx:   ctx,
						Key:   key,
						Value: new(dataCache),
						TTL:   tls,
						Do: func(i *cache.Item) (interface{}, error) {
							return &dataCache{
								Code:    bodyWriter.Status(),
								Data:    bodyWriter.body.String(),
								Headers: bodyWriter.Header(),
							}, nil
						},
					})
					if err != nil {
						logrus.Errorln("CacheRequest error:", err)
					}
				}(blw)
			} else {
				var jsonData interface{}
				jsonErr := common.StringToJson(data.Data, &jsonData)
				if jsonErr != nil {
					apiHandler(appCtx)(ginCtx)
					return
				}
				for key, val := range data.Headers {
					ginCtx.Header(key, val[0])
				}
				ginCtx.JSON(data.Code, jsonData)
			}
		}
	}
}
