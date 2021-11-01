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

func (r *ChatMember) Delete(ctx context.Context, chatId, userId int64) error {
	_, err := r.DB.ExecContext(ctx, `
	DELETE FROM chat_members
	WHERE chat_id = ? AND user_id = ?
	`, chatId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatMember) Exists(ctx context.Context, chatId, userId int64) (ok bool, err error) {
	query := `
	SELECT COUNT(*)
	FROM chat_members
	WHERE chat_id = ?
	`
	args := []interface{}{chatId}
	if userId != 0 {
		query += " AND user_id = ?"
		args = append(args, userId)
	}
	var cnt int
	err = r.DB.QueryRowContext(ctx, query, args...).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt != 0, nil
}

func (r *ChatMember) GetIds(ctx context.Context, chatId int64) (memberIds []int64, err error) {
	rows, err := r.DB.QueryContext(ctx, `
	SELECT user_id
	FROM chat_members
	WHERE chat_id = ?
	`, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var userId int64
		if err = rows.Scan(&userId); err != nil {
			return
		}
		memberIds = append(memberIds, userId)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
