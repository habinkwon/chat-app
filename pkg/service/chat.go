package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/neomarica/undergraduate-project/graph/model"
	"github.com/neomarica/undergraduate-project/pkg/middleware/auth"
	"github.com/neomarica/undergraduate-project/pkg/repository/mysql"
)

type Chat struct {
	ChatRepo        *mysql.Chat
	ChatMemberRepo  *mysql.ChatMember
	ChatMessageRepo *mysql.ChatMessage
	IDNode          *snowflake.Node
}

func (s *Chat) CreateChat(ctx context.Context, userIds []int64) (chatId int64, err error) {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return 0, auth.ErrNoAuth
	}
	for _, id := range userIds {
		if id == 0 || id == userId {
			return 0, fmt.Errorf("invalid user id %v", id)
		}
	}
	if len(userIds) == 0 {
		return 0, fmt.Errorf("no user ids specified")
	}
	userIds = append(userIds, userId)
	if chatId, err = s.ChatRepo.Exists(ctx, userIds); err != nil {
		return 0, err
	} else if chatId != 0 {
		return
	}
	chatId = s.IDNode.Generate().Int64()
	now := time.Now()
	if err := s.ChatRepo.Add(ctx, chatId, userId, now); err != nil {
		return 0, err
	}
	for _, id := range userIds {
		if err := s.ChatMemberRepo.Add(ctx, chatId, id, userId, now); err != nil {
			return 0, err
		}
	}
	return
}

func (s *Chat) GetChat(ctx context.Context, id int64) (*model.Chat, error) {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return nil, auth.ErrNoAuth
	}
	return s.ChatRepo.Get(ctx, id, userId)
}

func (s *Chat) ListChats(ctx context.Context, first int, after int64) ([]*model.Chat, error) {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return nil, auth.ErrNoAuth
	}
	return s.ChatRepo.List(ctx, userId, first, after)
}
