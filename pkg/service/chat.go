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

func (s *Chat) CreateChat(ctx context.Context, userIds []int64) (id int64, err error) {
	inviterId := auth.UserId(ctx)
	if inviterId == 0 {
		return 0, auth.ErrNoAuth
	}
	for _, id := range userIds {
		if id == 0 || id == inviterId {
			return 0, fmt.Errorf("invalid user id %v", id)
		}
	}
	if len(userIds) == 0 {
		return 0, fmt.Errorf("no user ids specified")
	}
	userIds = append(userIds, inviterId)
	if id, err = s.ChatRepo.Exists(ctx, userIds); err != nil {
		return 0, err
	} else if id != 0 {
		return
	}
	id = s.IDNode.Generate().Int64()
	now := time.Now()
	if err := s.ChatRepo.Add(ctx, id, inviterId, now); err != nil {
		return 0, err
	}
	for _, userId := range userIds {
		if err := s.ChatMemberRepo.Add(ctx, id, userId, inviterId, now); err != nil {
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

func (s *Chat) PostMessage(ctx context.Context, chatId int64, content string) (id int64, err error) {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return 0, auth.ErrNoAuth
	}
	if ok, err := s.ChatMemberRepo.Exists(ctx, chatId, userId); err != nil {
		return 0, err
	} else if !ok {
		return 0, auth.ErrPerm
	}
	id = s.IDNode.Generate().Int64()
	now := time.Now()
	if err := s.ChatMessageRepo.Add(ctx, id, chatId, content, userId, now); err != nil {
		return 0, err
	}
	return
}

func (s *Chat) DeleteMessage(ctx context.Context, id int64) error {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return auth.ErrNoAuth
	}
	if senderId, err := s.ChatMessageRepo.GetSenderId(ctx, id); err != nil {
		return err
	} else if senderId != userId {
		return auth.ErrPerm
	}
	if err := s.ChatMessageRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *Chat) EditMessage(ctx context.Context, id int64, content string) error {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return auth.ErrNoAuth
	}
	if senderId, err := s.ChatMessageRepo.GetSenderId(ctx, id); err != nil {
		return err
	} else if senderId != userId {
		return auth.ErrPerm
	}
	now := time.Now()
	if err := s.ChatMessageRepo.Edit(ctx, id, content, now); err != nil {
		return err
	}
	return nil
}

func (s *Chat) ListMessages(ctx context.Context, chatId int64, first int, after int64, desc bool) ([]*model.Message, error) {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return nil, auth.ErrNoAuth
	}
	return s.ChatMessageRepo.List(ctx, chatId, first, after, desc)
}
