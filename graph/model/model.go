package model

import "time"

type Chat struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	CreatorID    int64     `json:"createdBy"`
	CreatedAt    time.Time `json:"createdAt"`
	LastPostedAt time.Time `json:"lastPostedAt"`
}

type Message struct {
	ID        int64       `json:"id"`
	Type      MessageType `json:"type"`
	Content   string      `json:"content"`
	Event     string      `json:"event"`
	SenderID  int64       `json:"sender"`
	ReplyToID *int64      `json:"replyTo"`
	CreatedAt time.Time   `json:"createdAt"`
	EditedAt  *time.Time  `json:"editedAt"`
}
