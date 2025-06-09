package middlewares

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spinmozgJr/note-service/internal/auth"
	"github.com/spinmozgJr/note-service/internal/dependencies"
	"github.com/spinmozgJr/note-service/internal/httpx"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

var noAuthRequired = []string{
	"/sign-in",
	"/login",
}

//var noAuthRequiredPrefixes = []string{
//	"/docs/",
//}

func NewLoggerMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middlewares/logger"),
		)

		log.Info("logger middlewares enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

func corsMiddleware() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "X-Requested-With", "Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}

func JwtAuthMiddleware(deps *dependencies.Dependencies) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middlwares.JwtAuthMiddleware"
			deps.Log = deps.Log.With(
				slog.String("op", op),
			)

			if !pathRequiresAuth(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				deps.Log.Error("отсутствует заголовок 'Authorization'", "error")
				httpx.SendErrorJSON(w, r, http.StatusUnauthorized, errors.New("отсутствует заголовок 'Authorization'"))
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				deps.Log.Error("ошибка разбора токена", "error")
				httpx.SendErrorJSON(w, r, http.StatusUnauthorized, errors.New("ошибка разбора токена"))
				return
			}

			//var userClaims auth.UserClaims
			claims := jwt.MapClaims{}
			err := deps.TokenManager.ParseToken(token, claims)
			if err != nil {
				deps.Log.Error("ошибка парсинга токена", "error", err)
				httpx.SendErrorJSON(w, r, http.StatusUnauthorized, err)
				return
			}

			ctx := auth.IntoContext(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func pathRequiresAuth(path string) bool {
	for _, p := range noAuthRequired {
		if p == path {
			return false
		}
	}
	//for _, prefix := range noAuthRequiredPrefixes {
	//	if strings.HasPrefix(path, prefix) {
	//		return false
	//	}
	//}
	return true
}
