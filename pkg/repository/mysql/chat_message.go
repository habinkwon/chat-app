package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/neomarica/undergraduate-project/graph/model"
)

type ChatMessage struct {
	DB *sql.DB
}

func (r *ChatMessage) Add(ctx context.Context, chatId int64, m *model.Message) error {
	_, err := r.DB.ExecContext(ctx, `
	INSERT INTO chat_messages (id, chat_id, content, sender_id, created_at)
	VALUES (?, ?, ?, ?, ?)
	`, m.ID, chatId, m.Content, m.SenderID, m.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatMessage) Delete(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, `
	DELETE FROM chat_messages
	WHERE id = ?
	`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatMessage) Edit(ctx context.Context, id int64, content string, tm time.Time) error {
	_, err := r.DB.ExecContext(ctx, `
	UPDATE chat_messages
	SET content = ?, edited_at = ?
	WHERE id = ?
	`, content, tm, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatMessage) GetMetadata(ctx context.Context, id int64) (chatId, senderId int64, err error) {
	err = r.DB.QueryRowContext(ctx, `
	SELECT chat_id, sender_id
	FROM chat_messages
	WHERE id = ?
	`, id).Scan(&chatId, &senderId)
	if err == sql.ErrNoRows {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}
	return
}

func (r *ChatMessage) List(ctx context.Context, chatId int64, first int, after int64, desc bool) (messages []*model.Message, err error) {
	order := "ASC"
	if desc {
		order = "DESC"
	}
	rows, err := r.DB.QueryContext(ctx, fmt.Sprintf(`
	SELECT id, content, sender_id, created_at, edited_at
	FROM chat_messages
	WHERE chat_id = ? AND id > ?
	ORDER BY id %s
	LIMIT ?
	`, order), chatId, after, first)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id        int64
			content   sql.NullString
			senderId  int64
			createdAt time.Time
			editedAt  sql.NullTime
		)
		if err = rows.Scan(&id, &content, &senderId, &createdAt, &editedAt); err != nil {
			return
		}
		messages = append(messages, &model.Message{
			ID:        id,
			Content:   content.String,
			SenderID:  senderId,
			CreatedAt: createdAt,
			EditedAt:  &editedAt.Time,
		})
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
