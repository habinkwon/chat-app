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

type Middleware struct {
	Secret []byte
}

func (m *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ai := &AuthInfo{}
		ctx := context.WithValue(r.Context(), authInfoKey{}, ai)
		if auth := r.Header.Get("Authorization"); auth != "" {
			if err := m.Authenticate(ctx, auth); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) Authenticate(ctx context.Context, authorization string) error {
	token := strings.TrimPrefix(authorization, "Bearer ")
	p := &jwt.Parser{SkipClaimsValidation: true}
	tok, err := p.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return m.Secret, nil
	})
	if err != nil {
		return err
	}
	claims, _ := tok.Claims.(jwt.MapClaims)
	userId, ok := claims["userID"].(float64)
	if !ok {
		return errors.New("invalid token")
	}
	ai := AuthInfoFrom(ctx)
	ai.UserId = int64(userId)
	return nil
}
