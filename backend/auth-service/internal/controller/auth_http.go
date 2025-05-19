package controller

import (
	"net/http"

	"github.com/forum-backend/auth-service/internal/entity"
	"github.com/forum-backend/auth-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AuthHTTPController struct {
	authUC usecase.AuthUseCase
}

func NewAuthHTTPController(authUC usecase.AuthUseCase) *AuthHTTPController {
	return &AuthHTTPController{authUC: authUC}
}

// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body entity.UserCreate true "User registration data"
// @Success 201 {object} entity.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/register [post]
func (c *AuthHTTPController) Register(ctx *gin.Context) {
	var input entity.UserCreate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authUC.Register(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// @Summary Login user
// @Description Authenticate user and get tokens
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body entity.UserLogin true "User credentials"
// @Success 200 {object} entity.Tokens
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (c *AuthHTTPController) Login(ctx *gin.Context) {
	var input entity.UserLogin
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := c.authUC.Login(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

// @Summary Refresh tokens
// @Description Get new access token using refresh token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param refresh_token body string true "Refresh token"
// @Success 200 {object} entity.Tokens
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/refresh [post]
func (c *AuthHTTPController) RefreshToken(ctx *gin.Context) {
	// Реализация обновления токенов
}
