package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/neomarica/undergraduate-project/graph/generated"
	"github.com/neomarica/undergraduate-project/graph/model"
	"github.com/neomarica/undergraduate-project/pkg/util"
)

func (r *chatResolver) Members(ctx context.Context, obj *model.Chat, first *int, after *int64) ([]*model.User, error) {
	return r.UserSvc.GetUsers(ctx, obj.MemberIDs)
}

func (r *chatResolver) Messages(ctx context.Context, obj *model.Chat, first *int, after *int64, desc *bool) ([]*model.Message, error) {
	return nil, nil
}

func (r *chatResolver) CreatedBy(ctx context.Context, obj *model.Chat) (*model.User, error) {
	return r.UserSvc.GetUser(ctx, obj.CreatorID)
}

func (r *mutationResolver) CreateChat(ctx context.Context, userIds []int64) (*model.Chat, error) {
	chatId, err := r.ChatSvc.CreateChat(ctx, userIds)
	if err != nil {
		return nil, err
	}
	return &model.Chat{
		ID: chatId,
	}, nil
}

func (r *mutationResolver) PostMessage(ctx context.Context, chatID int64, text string, replyTo *int64) (*model.Message, error) {
	return nil, nil
}

func (r *mutationResolver) EditMessage(ctx context.Context, id int64, text string) (*model.Message, error) {
	return nil, nil
}

func (r *mutationResolver) DeleteMessage(ctx context.Context, id int64) (*model.Message, error) {
	return nil, nil
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	return r.UserSvc.GetUser(ctx, 0)
}

func (r *queryResolver) User(ctx context.Context, id int64) (*model.User, error) {
	return r.UserSvc.GetUser(ctx, id)
}

func (r *queryResolver) Chat(ctx context.Context, id int64) (*model.Chat, error) {
	return r.ChatSvc.GetChat(ctx, id)
}

func (r *queryResolver) Chats(ctx context.Context, first *int, after *int64) ([]*model.Chat, error) {
	return r.ChatSvc.ListChats(ctx, util.IntOr(first, 10), util.Int64Or(after, 0))
}

func (r *subscriptionResolver) ChatEvent(ctx context.Context, chatID int64) (<-chan *model.ChatEvent, error) {
	ch := make(chan *model.ChatEvent)
	go func() {
		defer close(ch)
		t := time.NewTicker(time.Second * 5)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				ch <- &model.ChatEvent{
					Type: model.ChatEventTypeMessagePosted,
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch, nil
}

// Chat returns generated.ChatResolver implementation.
func (r *Resolver) Chat() generated.ChatResolver { return &chatResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type chatResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
