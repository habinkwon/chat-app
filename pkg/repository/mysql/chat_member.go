package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type ChatMember struct {
	DB *sql.DB
}

func (r *ChatMember) Add(ctx context.Context, chatId, userId, inviterId int64, tm time.Time) error {
	_, err := r.DB.ExecContext(ctx, `
	INSERT INTO chat_members (chat_id, user_id, added_by, added_at)
	VALUES (?, ?, ?, ?)
	`, chatId, userId, inviterId, tm)
	if err != nil {
		return fmt.Errorf("error adding chat member: %w", err)
	}
	return nil
}
