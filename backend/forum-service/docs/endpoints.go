package docs

// @Summary Проверка работоспособности сервиса
// @Description Возвращает статус работоспособности сервиса
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
type HealthCheckEndpoint struct{}

// @Summary Получение списка постов
// @Description Возвращает список постов с пагинацией
// @Tags Posts
// @Produce json
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param page_size query int false "Размер страницы (по умолчанию 10)"
// @Success 200 {object} ListPostsResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts [get]
type ListPostsEndpoint struct{}

// @Summary Создание нового поста
// @Description Создает новый пост
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body CreatePostRequest true "Данные поста"
// @Success 201 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /posts [post]
type CreatePostEndpoint struct{}

// @Summary Получение поста по ID
// @Description Возвращает пост по его ID
// @Tags Posts
// @Produce json
// @Param id path int true "ID поста"
// @Success 200 {object} Post
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts/{id} [get]
type GetPostEndpoint struct{}

// @Summary Обновление поста
// @Description Обновляет существующий пост
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "ID поста"
// @Param post body UpdatePostRequest true "Данные для обновления"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /posts/{id} [put]
type UpdatePostEndpoint struct{}

// @Summary Удаление поста
// @Description Удаляет существующий пост
// @Tags Posts
// @Produce json
// @Param id path int true "ID поста"
// @Success 200 {object} SuccessResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /posts/{id} [delete]
type DeletePostEndpoint struct{}

// @Summary Получение списка ответов к посту
// @Description Возвращает список ответов к посту с пагинацией
// @Tags Replies
// @Produce json
// @Param postID path int true "ID поста"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param page_size query int false "Размер страницы (по умолчанию 10)"
// @Success 200 {object} ListRepliesResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts/{postID}/replies [get]
type ListRepliesEndpoint struct{}

// @Summary Создание нового ответа
// @Description Создает новый ответ к посту
// @Tags Replies
// @Accept json
// @Produce json
// @Param postID path int true "ID поста"
// @Param reply body CreateReplyRequest true "Данные ответа"
// @Success 201 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /posts/{postID}/replies [post]
type CreateReplyEndpoint struct{}

// @Summary Удаление ответа
// @Description Удаляет существующий ответ
// @Tags Replies
// @Produce json
// @Param postID path int true "ID поста"
// @Param replyID path int true "ID ответа"
// @Success 200 {object} SuccessResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /posts/{postID}/replies/{replyID} [delete]
type DeleteReplyEndpoint struct{}
