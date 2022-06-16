package decorator

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"
	"user_management/common"
	"user_management/components/appctx"

	"github.com/gin-gonic/gin"
)

type dataReponse struct {
	Code interface{} `json:"code"`
	Data interface{} `json:"data"`
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func CacheRequest(
	appCtx appctx.AppContext,
	cacheName string,
	tls time.Duration,
	funcHandler func(appctx.AppContext) func(*gin.Context),
) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := context.Background()
		redisService := appCtx.GetRedisService()
		url := ginCtx.Request.URL
		key := fmt.Sprintf("%s:%s:%s", cacheName, url.Path, url.Query().Encode())

		if cacheString, err := redisService.GetStringValue(ctx, key); err != nil {
			blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ginCtx.Writer}
			ginCtx.Writer = blw
			funcHandler(appCtx)(ginCtx)
			go func(bodyWriter *bodyLogWriter) {
				json := &dataReponse{
					Code: bodyWriter.Status(),
					Data: bodyWriter.body.String(),
				}
				jsonStr, jsonErr := common.JsonToString(json)
				if jsonErr != nil {
					log.Println("Cache error:", jsonErr)
				}

				str, err := redisService.SetValue(
					ctx,
					key,
					jsonStr,
					tls,
				)
				if err != nil {
					log.Println("Cache error:", err)
				} else {
					log.Println("Cache SetValue:", str)
				}
			}(blw)
		} else {
			dataCache := &dataReponse{}
			err = common.StringToJson(cacheString, &dataCache)
			if err != nil {
				funcHandler(appCtx)(ginCtx)
				return
			}
			statusCode := int(dataCache.Code.(float64))
			data := map[string]interface{}{}
			jsonStr := dataCache.Data.(string)
			err = common.StringToJson(jsonStr, &data)
			if err != nil {
				funcHandler(appCtx)(ginCtx)
				return
			} else {
				ginCtx.JSON(statusCode, data)
			}
		}
	}
}
