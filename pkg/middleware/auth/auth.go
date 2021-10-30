package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrNoAuth = errors.New("authentication required")
	ErrPerm   = errors.New("operation not permitted")
)

type AuthInfo struct {
	UserId int64
}

func UserId(ctx context.Context) int64 {
	return AuthInfoFrom(ctx).UserId
}

func SetUserId(ctx context.Context, userId int64) {
	AuthInfoFrom(ctx).UserId = userId
}

type authInfoKey struct{}

func AuthInfoFrom(ctx context.Context) *AuthInfo {
	ai, _ := ctx.Value(authInfoKey{}).(*AuthInfo)
	return ai
}

func Middleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			tok, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			claims, _ := tok.Claims.(jwt.MapClaims)
			userId, ok := claims["userID"].(float64)
			if !ok {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ai := &AuthInfo{
				UserId: int64(userId),
			}
			ctx := context.WithValue(r.Context(), authInfoKey{}, ai)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
