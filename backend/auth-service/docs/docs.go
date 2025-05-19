//go:generate swag init -g ../../cmd/main.go -o ./ --parseDependency --parseInternal

package docs

// @title Auth Service API
// @version 1.0
// @description API для микросервиса аутентификации

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

type swagger struct{}

var Swagger *swagger
