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

		ginCtx.JSON(http.StatusOK, common.SimpleSuccessResponse(user))
	}
}
```

2. Download Swag for Go by using:
```bash
// install swag
go install github.com/swaggo/swag/cmd/swag

// config GOPATH (macos)
export PATH=$(go env GOPATH)/bin:$PATH
```

3. Run the Swag at your Go api (for instance ~/modules/user/user_api.go), Swag will parse comments and generate required files(docs folder and docs/doc.go) at ~/docs.

```bash
// for only one api
swag init --parseDependency --parseInternal -d packages/rest_api
```

4. Download gin-swagger by using:

Import:

```bash
import "github.com/swaggo/gin-swagger" // gin-swagger middleware
import "github.com/swaggo/files" // swagger embed files
```

## Demo:
 [http://localhost:8080/swagger/index.html#/](http://localhost:8080/swagger/index.html#/)