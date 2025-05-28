package user_handler

import (
	"auth-service/internal/entity"
	"auth-service/internal/usecase"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	authUC *usecase.AuthUseCase
}

func NewUserHandler(authUC *usecase.AuthUseCase) *UserHandler {
	return &UserHandler{
		authUC: authUC,
	}
}

// @Summary Sign Up
// @Tags auth
// @Description Create new user account
// @Accept json
// @Produce json
// @Param input body entity.UserCreate true "user info"
// @Success 201 {object} map[string]interface{}
// @Failure 400,500 {object} map[string]interface{}
// @Router /auth/sign-up [post]
func (h *UserHandler) SignUp(c *gin.Context) {
	var input entity.UserCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authUC.SignUp(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary Sign In
// @Tags auth
// @Description Login with existing credentials
// @Accept json
// @Produce json
// @Param input body entity.SignInInput true "credentials"
// @Success 200 {object} entity.Tokens
// @Failure 400,401,500 {object} map[string]interface{}
// @Router /auth/sign-in [post]
func (h *UserHandler) SignIn(c *gin.Context) {
	var input entity.SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAgent := c.GetHeader("User-Agent")
	clientIP := c.ClientIP()

	tokens, err := h.authUC.SignIn(c.Request.Context(), &input, userAgent, clientIP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// @Summary Refresh Token
// @Tags auth
// @Description Refresh access token using refresh token
// @Accept json
// @Produce json
// @Param input body entity.RefreshInput true "refresh token"
// @Success 200 {object} entity.Tokens
// @Failure 400,401,500 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var input entity.RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAgent := c.GetHeader("User-Agent")
	clientIP := c.ClientIP()

	tokens, err := h.authUC.RefreshTokens(c.Request.Context(), input.RefreshToken, userAgent, clientIP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// @Summary Logout
// @Tags auth
// @Description Logout user
// @Accept json
// @Produce json
// @Param input body entity.RefreshInput true "refresh token"
// @Success 200 {object} map[string]interface{}
// @Failure 400,500 {object} map[string]interface{}
// @Router /auth/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	var input entity.RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authUC.Logout(c.Request.Context(), input.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Get Me
// @Tags profile
// @Description Get current user profile
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} entity.User
// @Failure 401,500 {object} map[string]interface{}
// @Router /api/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	user := c.MustGet("user").(*entity.User)
	c.JSON(http.StatusOK, user)
}

// @Summary Update Me
// @Tags profile
// @Description Update current user profile
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body entity.User true "user update info"
// @Success 200 {object} entity.User
// @Failure 400,401,500 {object} map[string]interface{}
// @Router /api/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	user := c.MustGet("user").(*entity.User)

	var updatedUser entity.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser.ID = user.ID
	if err := h.authUC.UpdateUser(c.Request.Context(), &updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// @Summary Update User Role
// @Tags admin
// @Description Update user's admin status
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param input body entity.User true "user info"
// @Success 200 {object} entity.User
// @Failure 400,401,403,500 {object} map[string]interface{}
// @Router /api/admin/users/{id}/role [put]
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var input entity.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authUC.UpdateUserRole(c.Request.Context(), userID, input.IsAdmin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Get Users
// @Tags admin
// @Description Get list of all users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} entity.User
// @Failure 401,403,500 {object} map[string]interface{}
// @Router /api/admin/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.authUC.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Get Stats
// @Tags admin
// @Description Get forum statistics
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} entity.Stats
// @Failure 401,403,500 {object} map[string]interface{}
// @Router /api/admin/stats [get]
func (h *UserHandler) GetStats(c *gin.Context) {
	stats := gin.H{
		"total_users":    0,
		"active_users":   0,
		"total_sessions": 0,
	}

	users, err := h.authUC.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	stats["total_users"] = len(users)
	activeUsers := 0
	for _, user := range users {
		if !user.IsBlocked {
			activeUsers++
		}
	}
	stats["active_users"] = activeUsers

	c.JSON(http.StatusOK, stats)
}

// @Summary Delete User
// @Tags admin
// @Description Delete user by ID
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400,401,403,500 {object} map[string]interface{}
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.authUC.DeleteUser(c.Request.Context(), userID); err != nil {
		if err.Error() == "cannot delete the last admin user" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get User by ID
// @Tags users
// @Description Get user information by ID
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entity.User
// @Failure 400,401,404,500 {object} map[string]interface{}
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.authUC.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ParseToken проверяет JWT токен и возвращает пользователя
func (h *UserHandler) ParseToken(ctx context.Context, token string) (*entity.User, error) {
	return h.authUC.ParseToken(ctx, token)
}
