package service

import (
	"context"

	"github.com/habinkwon/chat-app/graph/model"
	"github.com/habinkwon/chat-app/pkg/middleware/auth"
)

type User struct {
}

func (s *User) GetUser(ctx context.Context, id int64) (*model.User, error) {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return nil, auth.ErrNoAuth
	}
	if id == 0 {
		id = userId
	}
	return &model.User{ID: id}, nil
}

func (s *User) GetUsers(ctx context.Context, ids []int64) ([]*model.User, error) {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return nil, auth.ErrNoAuth
	}
	users := make([]*model.User, len(ids))
	for i, id := range ids {
		users[i] = &model.User{ID: id}
	}
	return users, nil
}
