package http

import (
	"forum-service/internal/entity"
	"forum-service/internal/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	postUC usecase.PostUseCase
}

func NewHandler(postUC usecase.PostUseCase) *Handler {
	return &Handler{
		postUC: postUC,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Routes
	api := router.Group("/api")
	{
		posts := api.Group("/posts")
		{
			posts.GET("", h.GetPosts)
			posts.GET("/:id", h.GetPost)
			posts.POST("", h.CreatePost)
			posts.PUT("/:id", h.UpdatePost)
			posts.DELETE("/:id", h.DeletePost)

			// Replies
			posts.GET("/:id/replies", h.GetReplies)
			posts.POST("/:id/replies", h.CreateReply)
		}
	}

	return router
}

// @Summary Get all posts
// @Description Get all posts with pagination
// @Tags posts
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {array} entity.Post
// @Router /api/posts [get]
func (h *Handler) GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	offset := (page - 1) * limit

	posts, total, err := h.postUC.GetAll(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if posts == nil {
		posts = make([]*entity.Post, 0)
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": total,
	})
}

// @Summary Get post by ID
// @Description Get post by ID with replies
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} entity.Post
// @Router /api/posts/{id} [get]
func (h *Handler) GetPost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	post, err := h.postUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// @Summary Create post
// @Description Create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Param input body entity.CreatePostInput true "Post input"
// @Success 201 {object} entity.Post
// @Router /api/posts [post]
func (h *Handler) CreatePost(c *gin.Context) {
	var input entity.CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	post, err := h.postUC.Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"post":    post,
		"message": "Пост успешно создан",
	})
}

// @Summary Update post
// @Description Update an existing post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param input body entity.UpdatePostInput true "Post update input"
// @Success 200 "No Content"
// @Router /api/posts/{id} [put]
func (h *Handler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var input entity.UpdatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.postUC.Update(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Delete post
// @Description Delete an existing post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 204 "No Content"
// @Router /api/posts/{id} [delete]
func (h *Handler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	if err := h.postUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get post replies
// @Description Get all replies for a post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {array} entity.Reply
// @Router /api/posts/{id}/replies [get]
func (h *Handler) GetReplies(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	replies, err := h.postUC.GetReplies(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, replies)
}

// @Summary Create reply
// @Description Create a new reply for a post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param input body entity.CreateReplyInput true "Reply input"
// @Success 201 {object} entity.Reply
// @Router /api/posts/{id}/replies [post]
func (h *Handler) CreateReply(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var input entity.CreateReplyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reply, err := h.postUC.CreateReply(c.Request.Context(), postID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reply)
}
