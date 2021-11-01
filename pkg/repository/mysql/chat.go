package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/habinkwon/chat-app/graph/model"
)

type Chat struct {
	DB *sql.DB
}

func (r *Chat) Add(ctx context.Context, id, ownerId int64, tm time.Time) error {
	_, err := r.DB.ExecContext(ctx, `
	INSERT INTO chats (id, created_by, created_at, last_posted_at)
	VALUES (?, ?, ?, ?)
	`, id, ownerId, tm, tm)
	if err != nil {
		return err
	}
	return nil
}

func (r *Chat) Delete(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, `
	DELETE FROM chats
	WHERE id = ?
	`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Chat) Exists(ctx context.Context, userIds []int64) (chatId int64, err error) {
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
	HAVING COUNT(*) = %d AND COUNT(*) = (
		SELECT COUNT(*)
		FROM chat_members B
		WHERE B.chat_id = A.chat_id
		GROUP BY B.chat_id
	)
	`, in, len(userIds)), args...).Scan(&chatId)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return
}

func (r *Chat) Get(ctx context.Context, id, userId int64) (chat *model.Chat, err error) {
	var (
		createdBy    int64
		createdAt    time.Time
		lastPostedAt time.Time
	)
	err = r.DB.QueryRowContext(ctx, `
	SELECT C.created_by, C.created_at, C.last_posted_at
	FROM chats C
	INNER JOIN chat_members M
	ON M.chat_id = C.id
	WHERE C.id = ? AND M.user_id = ?
	`, id, userId).Scan(&createdBy, &createdAt, &lastPostedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	chat = &model.Chat{
		ID:           id,
		CreatorID:    createdBy,
		CreatedAt:    createdAt,
		LastPostedAt: lastPostedAt,
	}
	return
}

func (r *Chat) List(ctx context.Context, userId int64, first int, after int64) (chats []*model.Chat, err error) {
	rows, err := r.DB.QueryContext(ctx, `
	SELECT C.id, C.created_by, C.created_at, C.last_posted_at
	FROM chats C
	INNER JOIN chat_members M
	ON M.chat_id = C.id
	WHERE M.user_id = ? AND C.id > ?
	ORDER BY C.last_posted_at
	LIMIT ?
	`, userId, after, first)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id           int64
			createdBy    int64
			createdAt    time.Time
			lastPostedAt time.Time
		)
		if err = rows.Scan(&id, &createdBy, &createdAt, &lastPostedAt); err != nil {
			return
		}
		chats = append(chats, &model.Chat{
			ID:           id,
			CreatorID:    createdBy,
			CreatedAt:    createdAt,
			LastPostedAt: lastPostedAt,
		})
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
