package router

import (
	"auth-service/internal/delivery/middleware"
	"auth-service/internal/delivery/user_handler"
	"auth-service/internal/usecase"
	"auth-service/pkg/jwt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	authUC       *usecase.AuthUseCase
	tokenManager jwt.TokenManager
}

func NewRouter(authUC *usecase.AuthUseCase, tokenManager jwt.TokenManager) *Router {
	return &Router{
		authUC:       authUC,
		tokenManager: tokenManager,
	}
}

func (r *Router) Init() *gin.Engine {
	router := gin.Default()

	// CORS настройка
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}

	router.Use(cors.New(config))

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(r.tokenManager, r.authUC)

	// Handlers
	userHandler := user_handler.NewUserHandler(r.authUC)

	// Public routes
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", userHandler.SignUp)
		auth.POST("/sign-in", userHandler.SignIn)
		auth.POST("/refresh", userHandler.RefreshToken)
		auth.POST("/logout", userHandler.Logout)
	}

	// Protected routes
	api := router.Group("/api", authMiddleware.Auth())
	{
		api.GET("/me", userHandler.GetMe)
		api.PUT("/me", userHandler.UpdateMe)

		// Admin endpoints moved directly under /api
		api.GET("/users", userHandler.GetUsers)
		api.PUT("/users/:id/role", userHandler.UpdateUserRole)
		api.GET("/stats", userHandler.GetStats)
	}

	return router
}

func (r *Router) Run(addr string) error {
	return r.Init().Run(addr)
}
