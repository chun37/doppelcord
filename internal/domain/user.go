package domain

import "time"

type User struct {
	ID           int64
	DiscordID    string
	RegisteredAt time.Time
	UpdatedAt    time.Time
}
