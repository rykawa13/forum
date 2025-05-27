package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"auth-service/internal/delivery/user_handler"
	"auth-service/internal/entity"
)

type RouterHandler struct {
	userHandler *user_handler.UserHandler
}

func NewRouterHandler(userHandler *user_handler.UserHandler) *RouterHandler {
	return &RouterHandler{
		userHandler: userHandler,
	}
}

func (h *RouterHandler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Настраиваем CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.userHandler.SignUp)
		auth.POST("/sign-in", h.userHandler.SignIn)
		auth.POST("/refresh", h.userHandler.RefreshToken)
		auth.POST("/logout", h.authMiddleware(), h.userHandler.Logout)
	}

	api := router.Group("/api", h.authMiddleware())
	{
		api.GET("/me", h.userHandler.GetMe)
		api.PUT("/me", h.userHandler.UpdateMe)

		// Маршруты для управления пользователями
		users := api.Group("/users", h.adminMiddleware())
		{
			users.PUT("/:id/role", h.userHandler.UpdateUserRole)
		}
	}

	return router
}

// adminMiddleware проверяет, является ли пользователь администратором
func (h *RouterHandler) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*entity.User)
		if !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// authMiddleware проверяет JWT токен и устанавливает пользователя в контекст
func (h *RouterHandler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
			c.Abort()
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			c.Abort()
			return
		}

		if len(headerParts[1]) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
			c.Abort()
			return
		}

		user, err := h.userHandler.ParseToken(c.Request.Context(), headerParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
