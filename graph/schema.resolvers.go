package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/neomarica/undergraduate-project/graph/generated"
	"github.com/neomarica/undergraduate-project/graph/model"
)

func (r *mutationResolver) CreateChat(ctx context.Context, userIds []int64) (*model.Chat, error) {
	chatId, err := r.Chat.CreateChat(ctx, userIds)
	if err != nil {
		return nil, err
	}
	return &model.Chat{
		ID: chatId,
	}, nil
}

func (r *mutationResolver) PostMessage(ctx context.Context, chatID int64, text string, replyTo *int64) (*model.Message, error) {
	return message, nil
}

func (r *mutationResolver) EditMessage(ctx context.Context, id int64, text string) (*model.Message, error) {
	return message, nil
}

func (r *mutationResolver) DeleteMessage(ctx context.Context, id int64) (*model.Message, error) {
	return message, nil
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	return user, nil
}

func (r *queryResolver) User(ctx context.Context, id int64) (*model.User, error) {
	return user, nil
}

func (r *queryResolver) Chats(ctx context.Context, first *int, after *string) ([]*model.Chat, error) {
	return []*model.Chat{chat}, nil
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
					Type:    model.ChatEventTypeMessagePosted,
					Message: message,
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
