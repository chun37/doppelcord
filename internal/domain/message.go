package domain

import "time"

type Message struct {
	ID        int64
	DiscordID string
	ChannelID string
	MessageID string
	Content   string
	CreatedAt time.Time
	StoredAt  time.Time
}
