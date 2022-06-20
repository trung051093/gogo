package routes

import (
	"user_management/components/appctx"
	"user_management/docs"

	"github.com/gin-gonic/gin" // swagger embed files
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func SwaggerRoutes(appCtx appctx.AppContext, router *gin.Engine) {
	config := appCtx.GetConfig()
	swaggerInfo := config.GetSwaggerConfig()

	// programmatically set swagger info
	docs.SwaggerInfo.Title = swaggerInfo.Title
	docs.SwaggerInfo.Description = swaggerInfo.Description
	docs.SwaggerInfo.Version = swaggerInfo.Version
	docs.SwaggerInfo.Schemes = swaggerInfo.Schemes
	docs.SwaggerInfo.Host = swaggerInfo.Host
	docs.SwaggerInfo.BasePath = swaggerInfo.BasePath

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
