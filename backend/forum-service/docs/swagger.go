package docs

import "github.com/swaggo/swag"

// @title Forum Service API
// @version 1.0
// @description API сервиса форума

// @host localhost:8082
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Токен в формате: Bearer <token>

var SwaggerInfo = &swag.Spec{
	Version:     "1.0",
	Host:        "localhost:8082",
	BasePath:    "/",
	Title:       "Forum Service API",
	Description: "API сервиса форума",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
