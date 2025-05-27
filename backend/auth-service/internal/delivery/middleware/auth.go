package middleware

import (
	"auth-service/internal/usecase"
	"auth-service/pkg/jwt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenManager jwt.TokenManager
	authUC       *usecase.AuthUseCase
}

func NewAuthMiddleware(tokenManager jwt.TokenManager, authUC *usecase.AuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
		authUC:       authUC,
	}
}

func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			return
		}

		userID, err := m.tokenManager.Parse(headerParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Преобразуем ID в число
		id, err := strconv.Atoi(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			return
		}

		// Получаем пользователя из базы данных
		user, err := m.authUC.GetUserByID(c.Request.Context(), id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
