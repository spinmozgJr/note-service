package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

type TokenManager interface {
	NewJWT(userId string, ttl time.Duration) (string, error)
	ParseToken(tokenStr string, dest jwt.Claims) error
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(userId string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		Subject:   userId,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) ParseToken(tokenStr string, dest jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenStr, dest, func(token *jwt.Token) (any, error) {
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("токен недействителен")
	}
	return nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
