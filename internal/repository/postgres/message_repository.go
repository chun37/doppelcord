package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/chun37/doppelcord/internal/domain"
	"github.com/chun37/doppelcord/internal/repository"
)

type messageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) repository.MessageRepository {
	return &messageRepository{pool: pool}
}

func (r *messageRepository) Save(ctx context.Context, msg *domain.Message) error {
	query := `
		INSERT INTO messages (discord_id, channel_id, message_id, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (message_id, created_at) DO NOTHING
	`
	_, err := r.pool.Exec(ctx, query,
		msg.DiscordID, msg.ChannelID, msg.MessageID, msg.Content, msg.CreatedAt,
	)
	return err
}

func (r *messageRepository) FindByDiscordID(ctx context.Context, discordID string, limit int, before *time.Time) ([]*domain.Message, error) {
	query := `
		SELECT id, discord_id, channel_id, message_id, content, created_at, stored_at
		FROM messages
		WHERE discord_id = $1
		  AND ($2::timestamptz IS NULL OR created_at < $2)
		ORDER BY created_at DESC
		LIMIT $3
	`
	rows, err := r.pool.Query(ctx, query, discordID, before, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(
			&msg.ID, &msg.DiscordID, &msg.ChannelID, &msg.MessageID,
			&msg.Content, &msg.CreatedAt, &msg.StoredAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}
