package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/neomarica/undergraduate-project/pkg/middleware/auth"
	"github.com/neomarica/undergraduate-project/pkg/repository/mysql"
)

type Chat struct {
	Chat        *mysql.Chat
	ChatMember  *mysql.ChatMember
	ChatMessage *mysql.ChatMessage
	IDNode      *snowflake.Node
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
	if chatId, err = s.ChatMember.Exists(ctx, userIds); err != nil {
		return 0, err
	} else if chatId != 0 {
		return
	}
	chatId = s.IDNode.Generate().Int64()
	now := time.Now()
	if err := s.Chat.Add(ctx, chatId, userId, now); err != nil {
		return 0, err
	}
	for _, id := range userIds {
		if err := s.ChatMember.Add(ctx, chatId, id, userId, now); err != nil {
			return 0, err
		}
	}
	return
}
