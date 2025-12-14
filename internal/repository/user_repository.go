package repository

import (
	"context"

	"github.com/chun37/doppelcord/internal/domain"
)

type UserRepository interface {
	IsRegistered(ctx context.Context, discordID string) (bool, error)
	Register(ctx context.Context, discordID string) (*domain.User, error)
	GetAllDiscordIDs(ctx context.Context) ([]string, error)
}
