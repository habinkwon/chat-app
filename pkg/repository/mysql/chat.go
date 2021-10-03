package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Chat struct {
	DB *sql.DB
}

func (r *Chat) Add(ctx context.Context, id, ownerId int64, tm time.Time) error {
	_, err := r.DB.ExecContext(ctx, `
	INSERT INTO chats (id, created_by, created_at)
	VALUES (?, ?, ?)
	`, id, ownerId, tm)
	if err != nil {
		return fmt.Errorf("error adding chat: %w", err)
	}
	return nil
}
