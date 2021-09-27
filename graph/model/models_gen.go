// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Chat struct {
	ID        string     `json:"id"`
	Cursor    *string    `json:"cursor"`
	Name      string     `json:"name"`
	Members   []*User    `json:"members"`
	Messages  []*Message `json:"messages"`
	CreatedAt time.Time  `json:"createdAt"`
	CreatedBy *User      `json:"createdBy"`
}

type ChatEvent struct {
	Type    ChatEventType `json:"type"`
	Message *Message      `json:"message"`
}

type Message struct {
	ID        string      `json:"id"`
	Cursor    *string     `json:"cursor"`
	Type      MessageType `json:"type"`
	Text      string      `json:"text"`
	Event     string      `json:"event"`
	Sender    *User       `json:"sender"`
	ReplyTo   *Message    `json:"replyTo"`
	Replies   []*Message  `json:"replies"`
	CreatedAt time.Time   `json:"createdAt"`
	EditedAt  *time.Time  `json:"editedAt"`
}

type User struct {
	ID        string    `json:"id"`
	Cursor    *string   `json:"cursor"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type ChatEventType string

const (
	ChatEventTypeMessagePosted  ChatEventType = "MESSAGE_POSTED"
	ChatEventTypeMessageEdited  ChatEventType = "MESSAGE_EDITED"
	ChatEventTypeMessageDeleted ChatEventType = "MESSAGE_DELETED"
)

var AllChatEventType = []ChatEventType{
	ChatEventTypeMessagePosted,
	ChatEventTypeMessageEdited,
	ChatEventTypeMessageDeleted,
}

func (e ChatEventType) IsValid() bool {
	switch e {
	case ChatEventTypeMessagePosted, ChatEventTypeMessageEdited, ChatEventTypeMessageDeleted:
		return true
	}
	return false
}

func (e ChatEventType) String() string {
	return string(e)
}

func (e *ChatEventType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ChatEventType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ChatEventType", str)
	}
	return nil
}

func (e ChatEventType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type MessageType string

const (
	MessageTypeMessage MessageType = "MESSAGE"
	MessageTypeEvent   MessageType = "EVENT"
)

var AllMessageType = []MessageType{
	MessageTypeMessage,
	MessageTypeEvent,
}

func (e MessageType) IsValid() bool {
	switch e {
	case MessageTypeMessage, MessageTypeEvent:
		return true
	}
	return false
}

func (e MessageType) String() string {
	return string(e)
}

func (e *MessageType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MessageType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MessageType", str)
	}
	return nil
}

func (e MessageType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
