package entity

// AuthStats содержит статистику по пользователям и сессиям
type AuthStats struct {
	TotalUsers    int `json:"total_users"`
	ActiveUsers   int `json:"active_users"`
	TotalSessions int `json:"total_sessions"`
}

type ForumStats struct {
	TotalUsers  int `json:"total_users"`
	TotalTopics int `json:"total_topics"`
	TotalPosts  int `json:"total_posts"`
}
