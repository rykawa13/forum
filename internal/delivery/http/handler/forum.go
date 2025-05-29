package handler

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ForumHandler struct {
	useCase ForumUseCase
	logger  *zap.Logger
}

func NewForumHandler(useCase ForumUseCase, logger *zap.Logger) *ForumHandler {
	return &ForumHandler{
		useCase: useCase,
		logger:  logger,
	}
}

func (h *ForumHandler) RegisterRoutes(router *mux.Router) {
	// Routes will be implemented here
}
