package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/habinkwon/chat-app/graph/generated"
	"github.com/habinkwon/chat-app/graph/model"
	"github.com/habinkwon/chat-app/pkg/util"
)

func (r *chatResolver) Members(ctx context.Context, obj *model.Chat) ([]*model.User, error) {
	memberIds, err := r.ChatSvc.GetMemberIds(ctx, obj.ID)
	if err != nil {
		return nil, err
	}
	return r.UserSvc.GetUsers(ctx, memberIds)
}

func (r *chatResolver) Messages(ctx context.Context, obj *model.Chat, first *int, after *int64, desc *bool) ([]*model.Message, error) {
	return r.ChatSvc.ListMessages(ctx, obj.ID, util.IntOr(first, 10), util.Int64Or(after, 0), util.BoolOr(desc, false))
}

func (r *chatResolver) CreatedBy(ctx context.Context, obj *model.Chat) (*model.User, error) {
	return r.UserSvc.GetUser(ctx, obj.CreatorID)
}

func (r *chatEventResolver) User(ctx context.Context, obj *model.ChatEvent) (*model.User, error) {
	if obj.User == nil {
		return nil, nil
	}
	return r.UserSvc.GetUser(ctx, obj.User.ID)
}

func (r *messageResolver) Sender(ctx context.Context, obj *model.Message) (*model.User, error) {
	return r.UserSvc.GetUser(ctx, obj.SenderID)
}

func (r *messageResolver) ReplyTo(ctx context.Context, obj *model.Message) (*model.Message, error) {
	if obj.ReplyToID == nil {
		return nil, nil
	}
	return r.ChatSvc.GetMessage(ctx, *obj.ReplyToID)
}

func (r *mutationResolver) CreateChat(ctx context.Context, userIds []int64) (*model.Chat, error) {
	id, err := r.ChatSvc.CreateChat(ctx, userIds)
	if err != nil {
		return nil, err
	}
	chat := &model.Chat{
		ID: id,
	}
	return chat, nil
}

func (r *mutationResolver) DeleteChat(ctx context.Context, id int64) (*model.Chat, error) {
	if err := r.ChatSvc.DeleteChat(ctx, id); err != nil {
		return nil, err
	}
	chat := &model.Chat{
		ID: id,
	}
	return chat, nil
}

func (r *mutationResolver) PostMessage(ctx context.Context, chatID int64, text string, replyTo *int64) (*model.Message, error) {
	id, err := r.ChatSvc.PostMessage(ctx, chatID, text, replyTo)
	if err != nil {
		return nil, err
	}
	message := &model.Message{
		ID: id,
	}
	return message, nil
}

func (r *mutationResolver) EditMessage(ctx context.Context, id int64, text string) (*model.Message, error) {
	if err := r.ChatSvc.EditMessage(ctx, id, text); err != nil {
		return nil, err
	}
	message := &model.Message{
		ID: id,
	}
	return message, nil
}

func (r *mutationResolver) DeleteMessage(ctx context.Context, id int64) (*model.Message, error) {
	if err := r.ChatSvc.DeleteMessage(ctx, id); err != nil {
		return nil, err
	}
	message := &model.Message{
		ID: id,
	}
	return message, nil
}

func (r *mutationResolver) SetAsOnline(ctx context.Context) (*model.User, error) {
	return r.UserSvc.SetAsOnline(ctx)
}

func (r *mutationResolver) UserTyping(ctx context.Context, chatID int64) (*model.User, error) {
	return r.ChatSvc.UserTyping(ctx, chatID)
}

func (r *queryResolver) Chat(ctx context.Context, id int64) (*model.Chat, error) {
	return r.ChatSvc.GetChat(ctx, id)
}

func (r *queryResolver) Chats(ctx context.Context, first *int, after *int64) ([]*model.Chat, error) {
	return r.ChatSvc.ListChats(ctx, util.IntOr(first, 10), util.Int64Or(after, 0))
}

func (r *subscriptionResolver) ChatEvent(ctx context.Context) (<-chan *model.ChatEvent, error) {
	return r.ChatSvc.ReceiveEvents(ctx)
}

func (r *userResolver) Status(ctx context.Context, obj *model.User) (model.UserStatus, error) {
	return r.UserSvc.GetStatus(ctx, obj.ID)
}

// Chat returns generated.ChatResolver implementation.
func (r *Resolver) Chat() generated.ChatResolver { return &chatResolver{r} }

// ChatEvent returns generated.ChatEventResolver implementation.
func (r *Resolver) ChatEvent() generated.ChatEventResolver { return &chatEventResolver{r} }

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return &messageResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type chatResolver struct{ *Resolver }
type chatEventResolver struct{ *Resolver }
type messageResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
