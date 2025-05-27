package entity

// ForumStats представляет статистику форума
type ForumStats struct {
	TotalUsers  int `json:"totalUsers"`
	TotalTopics int `json:"totalTopics"`
	TotalPosts  int `json:"totalPosts"`
}
