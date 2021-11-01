package service

import (
	"context"

	"github.com/habinkwon/chat-app/graph/model"
	"github.com/habinkwon/chat-app/pkg/middleware/auth"
	"github.com/habinkwon/chat-app/pkg/repository/redis"
)

type User struct {
	UserStatusRepo *redis.UserStatus
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

func (s *User) SetAsOnline(ctx context.Context) (*model.User, error) {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return nil, auth.ErrNoAuth
	}
	status := model.UserStatusOnline
	if err := s.UserStatusRepo.SetStatus(ctx, userId, status); err != nil {
		return nil, err
	}
	user := &model.User{
		ID:     userId,
		Status: status,
	}
	return user, nil
}

func (s *User) GetStatus(ctx context.Context, userId int64) (model.UserStatus, error) {
	status, err := s.UserStatusRepo.GetStatus(ctx, userId)
	if err != nil {
		return model.UserStatusOffline, err
	}
	return status, nil
}
