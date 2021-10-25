package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/habinkwon/chat-app/graph/model"
)

type UserStatus struct {
	Redis *redis.Client
}

func (r *UserStatus) SetStatus(ctx context.Context, userId int64, status model.UserStatus) (err error) {
	key := fmt.Sprintf("user_status:%d", userId)
	if err = r.Redis.SetEX(ctx, key, string(status), time.Minute).Err(); err != nil {
		return err
	}
	return
}

func (r *UserStatus) GetStatus(ctx context.Context, userId int64) (model.UserStatus, error) {
	key := fmt.Sprintf("user_status:%d", userId)
	var status string
	if err := r.Redis.Get(ctx, key).Scan(&status); err == redis.Nil {
		return model.UserStatusOffline, nil
	} else if err != nil {
		return model.UserStatusOffline, err
	}
	return model.UserStatus(status), nil
}
