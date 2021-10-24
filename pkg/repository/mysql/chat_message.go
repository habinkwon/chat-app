package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/habinkwon/chat-app/graph/model"
)

type ChatMessage struct {
	DB *sql.DB
}

func (r *ChatMessage) Add(ctx context.Context, chatId int64, m *model.Message) error {
	_, err := r.DB.ExecContext(ctx, `
	INSERT INTO chat_messages (id, chat_id, content, sender_id, reply_to, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	`, m.ID, chatId, m.Content, m.SenderID, m.ReplyToID, m.CreatedAt)
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

func (r *ChatMessage) DeleteAll(ctx context.Context, chatId int64) error {
	_, err := r.DB.ExecContext(ctx, `
	DELETE FROM chat_messages
	WHERE chat_id = ?
	`, chatId)
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

func (r *ChatMessage) Exists(ctx context.Context, id, chatId int64) (ok bool, err error) {
	var cnt int
	err = r.DB.QueryRowContext(ctx, `
	SELECT COUNT(*)
	FROM chat_messages
	WHERE id = ? AND chat_id = ?
	`, id, chatId).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt != 0, nil
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

func (r *ChatMessage) Get(ctx context.Context, id, userId int64) (message *model.Message, err error) {
	var (
		content   sql.NullString
		senderId  int64
		replyTo   *int64
		createdAt time.Time
		editedAt  sql.NullTime
	)
	err = r.DB.QueryRowContext(ctx, `
	SELECT T.content, T.sender_id, T.reply_to, T.created_at, T.edited_at
	FROM chat_messages T
	INNER JOIN chat_members U
	ON U.chat_id = T.chat_id
	WHERE T.id = ? AND U.user_id = ?
	`, id, userId).Scan(&content, &senderId, &replyTo, &createdAt, &editedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	message = &model.Message{
		ID:        id,
		Content:   content.String,
		SenderID:  senderId,
		ReplyToID: replyTo,
		CreatedAt: createdAt,
		EditedAt:  &editedAt.Time,
	}
	return
}

func (r *ChatMessage) List(ctx context.Context, chatId int64, first int, after int64, desc bool) (messages []*model.Message, err error) {
	order := "ASC"
	if desc {
		order = "DESC"
	}
	rows, err := r.DB.QueryContext(ctx, fmt.Sprintf(`
	SELECT id, content, sender_id, reply_to, created_at, edited_at
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
			replyTo   *int64
			createdAt time.Time
			editedAt  sql.NullTime
		)
		if err = rows.Scan(&id, &content, &senderId, &replyTo, &createdAt, &editedAt); err != nil {
			return
		}
		messages = append(messages, &model.Message{
			ID:        id,
			Content:   content.String,
			SenderID:  senderId,
			ReplyToID: replyTo,
			CreatedAt: createdAt,
			EditedAt:  &editedAt.Time,
		})
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
