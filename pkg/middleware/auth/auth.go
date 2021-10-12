package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"
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

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ai := &AuthInfo{}
			if id := r.Header.Get("user-id"); id != "" {
				if id, err := strconv.ParseInt(id, 10, 64); err == nil {
					ai.UserId = id
				}
			}
			ctx := context.WithValue(r.Context(), authInfoKey{}, ai)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
