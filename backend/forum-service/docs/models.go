package docs

// Post представляет собой пост на форуме
type Post struct {
	// ID поста
	ID int64 `json:"id" example:"1"`
	// Заголовок поста
	Title string `json:"title" example:"Как использовать Go channels"`
	// Содержание поста
	Content string `json:"content" example:"В этом посте я расскажу о том, как эффективно использовать каналы в Go..."`
	// ID автора поста
	UserID int `json:"user_id" example:"42"`
	// Имя автора поста
	Username string `json:"username" example:"john_doe"`
	// Дата создания поста
	CreatedAt string `json:"created_at" example:"2024-01-20T15:04:05Z"`
	// Дата последнего обновления поста
	UpdatedAt string `json:"updated_at" example:"2024-01-20T15:04:05Z"`
	// Флаг блокировки поста для новых ответов
	IsLocked bool `json:"is_locked" example:"false"`
	// Количество просмотров поста
	ViewCount int `json:"view_count" example:"123"`
	// Количество ответов на пост
	ReplyCount int `json:"reply_count" example:"5"`
}

// CreatePostRequest представляет собой запрос на создание поста
type CreatePostRequest struct {
	// Заголовок поста
	Title string `json:"title" example:"Как использовать Go channels" validate:"required,min=5,max=200"`
	// Содержание поста
	Content string `json:"content" example:"В этом посте я расскажу..." validate:"required,min=10"`
}

// UpdatePostRequest представляет собой запрос на обновление поста
type UpdatePostRequest struct {
	// Новый заголовок поста (опционально)
	Title *string `json:"title,omitempty" example:"Обновленный заголовок" validate:"omitempty,min=5,max=200"`
	// Новое содержание поста (опционально)
	Content *string `json:"content,omitempty" example:"Обновленное содержание..." validate:"omitempty,min=10"`
	// Флаг блокировки поста (опционально)
	IsLocked *bool `json:"is_locked,omitempty" example:"true"`
}

// Reply представляет собой ответ на пост
type Reply struct {
	// ID ответа
	ID int64 `json:"id" example:"1"`
	// ID поста, к которому относится ответ
	PostID int64 `json:"post_id" example:"42"`
	// Содержание ответа
	Content string `json:"content" example:"Отличный пост! Хочу добавить, что..."`
	// ID автора ответа
	UserID int `json:"user_id" example:"15"`
	// Имя автора ответа
	Username string `json:"username" example:"jane_doe"`
	// Дата создания ответа
	CreatedAt string `json:"created_at" example:"2024-01-20T15:04:05Z"`
	// Дата последнего обновления ответа
	UpdatedAt string `json:"updated_at" example:"2024-01-20T15:04:05Z"`
}

// CreateReplyRequest представляет собой запрос на создание ответа
type CreateReplyRequest struct {
	// Содержание ответа
	Content string `json:"content" example:"Спасибо за пост!" validate:"required,min=1"`
}

// ErrorResponse представляет собой ответ с ошибкой
type ErrorResponse struct {
	// Сообщение об ошибке
	Error string `json:"error" example:"post not found"`
}

// SuccessResponse представляет собой успешный ответ
type SuccessResponse struct {
	// Статус операции
	Status string `json:"status" example:"updated"`
}

// IDResponse представляет собой ответ с ID созданного ресурса
type IDResponse struct {
	// ID созданного ресурса
	ID int64 `json:"id" example:"1"`
}

// ListPostsResponse представляет собой ответ со списком постов
type ListPostsResponse struct {
	// Список постов
	Posts []*Post `json:"posts"`
	// Общее количество постов
	Total int `json:"total" example:"42"`
}

// ListRepliesResponse представляет собой ответ со списком ответов
type ListRepliesResponse struct {
	// Список ответов
	Replies []*Reply `json:"replies"`
	// Общее количество ответов
	Total int `json:"total" example:"15"`
}
