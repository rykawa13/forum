package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenManager interface {
	NewJWT(userID int, ttl time.Duration) (string, error)
	Parse(accessToken string) (int, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, fmt.Errorf("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(userID int, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Subject:   fmt.Sprintf("%d", userID),
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (int, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("error get user claims from token")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, fmt.Errorf("error get subject from token")
	}

	var userID int
	if _, err := fmt.Sscanf(subject, "%d", &userID); err != nil {
		return 0, fmt.Errorf("error parsing user ID from subject")
	}

	return userID, nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	return token.SignedString([]byte(m.signingKey))
}
