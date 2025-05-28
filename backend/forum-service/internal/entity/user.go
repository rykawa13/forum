package entity

// User represents a user in the system
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}
