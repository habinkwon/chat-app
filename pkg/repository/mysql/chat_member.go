package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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
		return fmt.Errorf("ChatMember.Add: %w", err)
	}
	return nil
}

func (r *ChatMember) Exists(ctx context.Context, userIds []int64) (chatId int64, err error) {
	// https://stackoverflow.com/questions/12776178/sql-select-sets-containing-exactly-given-members
	in := strings.Repeat(", ?", len(userIds)-1)
	args := make([]interface{}, len(userIds))
	for i, id := range userIds {
		args[i] = id
	}
	err = r.DB.QueryRowContext(ctx, fmt.Sprintf(`
	SELECT chat_id
	FROM chat_members A
	WHERE user_id IN (?%s)
	GROUP BY chat_id
	HAVING COUNT(*) = (
		SELECT COUNT(*)
		FROM chat_members B
		WHERE B.chat_id = A.chat_id
		GROUP BY B.chat_id
	)
	`, in), args...).Scan(&chatId)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("ChatMember.Exists: %w", err)
	}
	return
}
