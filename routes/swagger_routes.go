package routes

import (
	"user_management/components/appctx"
	"user_management/docs"

	"github.com/gin-gonic/gin" // swagger embed files
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func SwaggerRoutes(appCtx appctx.AppContext, router *gin.Engine) {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Example API"
	docs.SwaggerInfo.Description = "This is a sample server"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
