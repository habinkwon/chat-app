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

type userIdKey struct{}

func UserId(ctx context.Context) int64 {
	userId, _ := ctx.Value(userIdKey{}).(int64)
	return userId
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userId int64
			if id := r.Header.Get("user-id"); id != "" {
				if id, err := strconv.ParseInt(id, 10, 64); err == nil {
					userId = id
				}
			}
			if userId != 0 {
				ctx := r.Context()
				ctx = context.WithValue(ctx, userIdKey{}, userId)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}
