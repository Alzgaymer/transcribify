package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
	"transcribify/internal/models"
)

const (
	Access  = 15 * time.Minute
	Refresh = 24 * time.Hour
)

type TokenManager interface {
	NewJWT(user *models.User, ttl time.Duration) (models.Token, error)
	Parse(accessToken string) (int, error)
}

type Manager struct {
	signingKey string
}

func (m *Manager) NewJWT(user *models.User, ttl time.Duration) (models.Token, error) {

	expires := time.Now().Add(ttl)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.ID
	claims["exp"] = expires.Unix()
	claims["role"] = user.Role

	signed, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return models.Token{}, err
	}

	t := models.Token{
		T:         signed,
		ExpiresAt: expires,
	}

	switch ttl {
	case Access:
		t.Key = "access"
	case Refresh:
		t.Key = "refresh"
	}

	return t, nil
}

func (m *Manager) Parse(accessToken string) (int, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return -1, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return -1, fmt.Errorf("error get user claims from token")
	}

	id, ok := claims["sub"].(float64)
	if !ok {
		return -1, err
	}
	return int(id), nil
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}
