package graph

import (
	"time"

	"github.com/neomarica/undergraduate-project/graph/model"
)

var (
	now  = time.Now()
	user = &model.User{
		ID:        123,
		Name:      "David",
		Username:  "david",
		Email:     "david@example.com",
		CreatedAt: now,
	}
	chat = &model.Chat{
		ID:        456,
		Name:      "A chat",
		Members:   []*model.User{user},
		Messages:  []*model.Message{message},
		CreatedAt: now,
		CreatedBy: user,
	}
	message = &model.Message{
		ID:        789,
		Type:      model.MessageTypeMessage,
		Text:      "A message",
		Sender:    user,
		ReplyTo:   message2,
		Replies:   []*model.Message{message2},
		CreatedAt: now,
		EditedAt:  &now,
	}
	message2 = &model.Message{
		ID:        1011,
		Type:      model.MessageTypeMessage,
		Text:      "A message",
		Sender:    user,
		CreatedAt: now,
		EditedAt:  &now,
	}
)
