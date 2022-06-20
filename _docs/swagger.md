# Swagger in go:
I believe that Swagger will make smooth communication for Back-end and Front-end. 

# Demo:
 [https://api.tdo.works/swagger/index.html#/](https://api.tdo.works/swagger/index.html#/)

# Usage
## Start using it
1. Add comments to your API source code, [See Declarative Comments Format.](https://github.com/swaggo/swag/blob/master/README.md#declarative-comments-format)

[Code example](../modules/user/user_api.go#L15-L24) 

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

- [Import](../routes/swagger_routes.go#8-9) 

- [Create swagger route](../routes/swagger_routes.go#12) 

4. Run the Swag at your api (for instance ~/packages/rest_api), Swag will parse comments and generate required files(docs folder and docs/doc.go) at ~/docs.

```bash
swag init --parseDependency --parseInternal -d packages/rest_api
```