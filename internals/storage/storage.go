package storage

import "time"

type Message struct {
	From      string    `json:"from"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	ID        int64     `json:"id"`
}
