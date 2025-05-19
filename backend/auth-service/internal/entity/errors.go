package entity

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrSessionNotFound     = errors.New("session not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
