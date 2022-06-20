# Swagger in go:
I believe that Swagger will make smooth communication for Back-end and Front-end. 

# Usage
## Start using it
1. Add comments to your API source code, [See Declarative Comments Format.](https://github.com/swaggo/swag/blob/master/README.md#declarative-comments-format)

```bash
// GetUser godoc
// @Summary      Get an user
// @Description  get user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  common.Response{data=usermodel.User} "desc"
// @Failure      400  {object}  common.AppError
// @Router       /api/v1/user/{id} [get]
func GetUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		id, err := strconv.Atoi(ginCtx.Param("id"))

		if err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		userRepo := NewUserRepository(appCtx.GetMainDBConnection())
		esService := appCtx.GetESService()
		userService := NewUserService(userRepo, esService)

		user, err := userService.GetUserTrace(ginCtx.Request.Context(), id)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(user))
	}
}
```


2. Download Swag for Go by using:
```bash
// install swag
go install github.com/swaggo/swag/cmd/swag

// config GOPATH (macos)
export PATH=$(go env GOPATH)/bin:$PATH

// test cmd
swag init -h
```


3. Download gin-swagger by using:

Import:

```bash
import "github.com/swaggo/gin-swagger" // gin-swagger middleware
import "github.com/swaggo/files" // swagger embed files
```

Create swagger route:

```bash
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

```


4. Run the Swag at your api (for instance ~/packages/rest_api), Swag will parse comments and generate required files(docs folder and docs/doc.go) at ~/docs.

```bash
swag init --parseDependency --parseInternal -d packages/rest_api
```
## Demo:
 [https://api.tdo.works/swagger/index.html#/](https://api.tdo.works/swagger/index.html#/)