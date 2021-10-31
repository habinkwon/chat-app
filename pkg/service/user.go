package service

import (
	"context"

	"github.com/habinkwon/chat-app/graph/model"
	"github.com/habinkwon/chat-app/pkg/middleware/auth"
	"github.com/habinkwon/chat-app/pkg/repository/mysql"
	"github.com/habinkwon/chat-app/pkg/repository/redis"
)

type User struct {
	UserRepo       *mysql.User
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
	return s.UserRepo.Get(ctx, id)
}

func (s *User) GetUsers(ctx context.Context, ids []int64) ([]*model.User, error) {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return nil, auth.ErrNoAuth
	}
	return s.UserRepo.GetAll(ctx, ids)
}

func (s *User) ListUsers(ctx context.Context, first int, after int64) ([]*model.User, error) {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return nil, auth.ErrNoAuth
	}
	return s.UserRepo.List(ctx, first, after)
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
