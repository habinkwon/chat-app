package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/neomarica/undergraduate-project/graph/model"
)

type Channel struct {
	Redis *redis.Client
}

func (r *Channel) SendEvent(ctx context.Context, userIds []int64, e *model.ChatEvent) (err error) {
	val, err := json.Marshal(e)
	if err != nil {
		return err
	}
	for _, userId := range userIds {
		channel := fmt.Sprintf("chat_events:%d", userId)
		if err = r.Redis.Publish(ctx, channel, val).Err(); err != nil {
			return err
		}
	}
	return
}

func (r *Channel) StreamEvents(ctx context.Context, userId int64, c chan *model.ChatEvent) (err error) {
	channel := fmt.Sprintf("chat_events:%d", userId)
	sub := r.Redis.Subscribe(ctx, channel)
	defer sub.Close()
	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			return err
		}
		e := new(model.ChatEvent)
		if err := json.Unmarshal([]byte(msg.Payload), e); err != nil {
			log.Print(err)
			continue
		}
		select {
		case c <- e:
		case <-ctx.Done():
			return nil
		}
	}
}
