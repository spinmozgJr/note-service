package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
)

type contextKey struct{}

var userKey = contextKey{}

// IntoContext сохраняет пользователя в контекст
func IntoContext(ctx context.Context, user jwt.MapClaims) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// FromContext достаёт пользователя из контекста
func FromContext(ctx context.Context) (jwt.Claims, bool) {
	user, ok := ctx.Value(userKey).(jwt.Claims)
	return user, ok
}

func GetUserIDFromRequest(r *http.Request) (int, error) {
	user, ok := FromContext(r.Context())
	if !ok {
		return 0, errors.New("пользователь не авторизован")
	}
	sub, err := user.GetSubject()
	if err != nil {
		return 0, errors.New("некорректный subject в токене")
	}
	userID, err := strconv.Atoi(sub)
	if err != nil {
		return 0, errors.New("subject не является числом")
	}

	return userID, nil
}
