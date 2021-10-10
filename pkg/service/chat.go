package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/neomarica/undergraduate-project/graph/model"
	"github.com/neomarica/undergraduate-project/pkg/middleware/auth"
	"github.com/neomarica/undergraduate-project/pkg/repository/mysql"
	"github.com/neomarica/undergraduate-project/pkg/repository/redis"
	"github.com/neomarica/undergraduate-project/pkg/util"
)

type Chat struct {
	IDNode          *snowflake.Node
	ChatRepo        *mysql.Chat
	ChatMemberRepo  *mysql.ChatMember
	ChatMessageRepo *mysql.ChatMessage
	ChannelRepo     *redis.Channel
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

func (s *Chat) DeleteChat(ctx context.Context, id int64) error {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return auth.ErrNoAuth
	}

	if ok, err := s.ChatMemberRepo.Exists(ctx, id, userId); err != nil {
		return err
	} else if !ok {
		return auth.ErrPerm
	}

	if err := s.ChatMemberRepo.Delete(ctx, id, userId); err != nil {
		return err
	}

	if ok, err := s.ChatMemberRepo.Exists(ctx, id, 0); err != nil {
		return err
	} else if !ok {
		if err := s.ChatRepo.Delete(ctx, id); err != nil {
			log.Print(err)
		}
		if err := s.ChatMessageRepo.DeleteAll(ctx, id); err != nil {
			log.Print(err)
		}
	}
	return nil
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

func (s *Chat) GetMemberIds(ctx context.Context, chatId int64) (memberIds []int64, err error) {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return nil, auth.ErrNoAuth
	}

	memberIds, err = s.ChatMemberRepo.Get(ctx, chatId)
	if err != nil {
		return nil, err
	} else if !util.ContainsInt64(memberIds, userId) {
		return nil, auth.ErrPerm
	}
	return
}

func (s *Chat) PostMessage(ctx context.Context, chatId int64, content string, replyTo *int64) (id int64, err error) {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return 0, auth.ErrNoAuth
	}

	memberIds, err := s.ChatMemberRepo.Get(ctx, chatId)
	if err != nil {
		return 0, err
	} else if !util.ContainsInt64(memberIds, userId) {
		return 0, auth.ErrPerm
	}

	if replyTo != nil {
		if ok, err := s.ChatMessageRepo.Exists(ctx, *replyTo, chatId); err != nil {
			return 0, err
		} else if !ok {
			return 0, fmt.Errorf("message %d does not exist in chat %d", *replyTo, chatId)
		}
	}

	m := &model.Message{
		ID:        s.IDNode.Generate().Int64(),
		Content:   content,
		SenderID:  userId,
		ReplyToID: replyTo,
		CreatedAt: time.Now(),
	}
	if err := s.ChatMessageRepo.Add(ctx, chatId, m); err != nil {
		return 0, err
	}

	e := &model.ChatEvent{
		Type:    model.ChatEventTypeMessagePosted,
		ChatID:  chatId,
		Message: m,
	}
	if err := s.ChannelRepo.SendEvent(ctx, memberIds, e); err != nil {
		log.Print(err)
	}
	return m.ID, nil
}

func (s *Chat) DeleteMessage(ctx context.Context, id int64) error {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return auth.ErrNoAuth
	}

	chatId, senderId, err := s.ChatMessageRepo.GetMetadata(ctx, id)
	if err != nil {
		return err
	} else if senderId != userId {
		return auth.ErrPerm
	}

	memberIds, err := s.ChatMemberRepo.Get(ctx, chatId)
	if err != nil {
		return err
	} else if !util.ContainsInt64(memberIds, userId) {
		return auth.ErrPerm
	}

	if err := s.ChatMessageRepo.Delete(ctx, id); err != nil {
		return err
	}

	now := time.Now()
	e := &model.ChatEvent{
		Type:   model.ChatEventTypeMessageDeleted,
		ChatID: chatId,
		Message: &model.Message{
			ID:       id,
			SenderID: senderId,
			EditedAt: &now,
		},
	}
	if err := s.ChannelRepo.SendEvent(ctx, memberIds, e); err != nil {
		log.Print(err)
	}
	return nil
}

func (s *Chat) EditMessage(ctx context.Context, id int64, content string) error {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return auth.ErrNoAuth
	}

	chatId, senderId, err := s.ChatMessageRepo.GetMetadata(ctx, id)
	if err != nil {
		return err
	} else if senderId != userId {
		return auth.ErrPerm
	}

	memberIds, err := s.ChatMemberRepo.Get(ctx, chatId)
	if err != nil {
		return err
	} else if !util.ContainsInt64(memberIds, userId) {
		return auth.ErrPerm
	}

	now := time.Now()
	if err := s.ChatMessageRepo.Edit(ctx, id, content, now); err != nil {
		return err
	}

	e := &model.ChatEvent{
		Type:   model.ChatEventTypeMessageEdited,
		ChatID: chatId,
		Message: &model.Message{
			ID:       id,
			Content:  content,
			SenderID: senderId,
			EditedAt: &now,
		},
	}
	if err := s.ChannelRepo.SendEvent(ctx, memberIds, e); err != nil {
		log.Print(err)
	}
	return nil
}

func (s *Chat) GetMessage(ctx context.Context, id int64) (*model.Message, error) {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return nil, auth.ErrNoAuth
	}
	return s.ChatMessageRepo.Get(ctx, id, userId)
}

func (s *Chat) ListMessages(ctx context.Context, chatId int64, first int, after int64, desc bool) ([]*model.Message, error) {
	userId := auth.UserId(ctx)
	if userId == 0 {
		return nil, auth.ErrNoAuth
	}
	return s.ChatMessageRepo.List(ctx, chatId, first, after, desc)
}

func (s *Chat) ReceiveEvents(ctx context.Context, userID int64) (<-chan *model.ChatEvent, error) {
	c := make(chan *model.ChatEvent, 1)
	go func() {
		defer close(c)
		if err := s.ChannelRepo.StreamEvents(ctx, userID, c); err != nil {
			log.Print(err)
		}
	}()
	return c, nil
}
