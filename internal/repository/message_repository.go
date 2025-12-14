package repository

import (
	"context"
	"time"

	"github.com/chun37/doppelcord/internal/domain"
)

type MessageRepository interface {
	Save(ctx context.Context, msg *domain.Message) error
	FindByDiscordID(ctx context.Context, discordID string, limit int, before *time.Time) ([]*domain.Message, error)
	FindByDiscordIDAndChannelID(ctx context.Context, discordID, channelID string, limit int) ([]*domain.Message, error)
}
