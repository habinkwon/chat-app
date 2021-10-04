package mysql

import (
	"context"
	"database/sql"
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
		return err
	}
	return nil
}

func (r *ChatMember) Exists(ctx context.Context, chatId, userId int64) (ok bool, err error) {
	var cnt int
	err = r.DB.QueryRowContext(ctx, `
	SELECT COUNT(*)
	FROM chat_members
	WHERE chat_id = ? AND user_id = ?
	`, chatId, userId).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt != 0 , nil
}
