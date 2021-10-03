package model

import "time"

type Chat struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	MemberIDs    []int64   `json:"members"`
	MessageIDs   []int64   `json:"messages"`
	CreatorID    int64     `json:"createdBy"`
	CreatedAt    time.Time `json:"createdAt"`
	LastPostedAt time.Time `json:"lastPostedAt"`
}
