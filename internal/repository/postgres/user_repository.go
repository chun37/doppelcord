package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/chun37/doppelcord/internal/domain"
	"github.com/chun37/doppelcord/internal/repository"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) IsRegistered(ctx context.Context, discordID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE discord_id = $1)`
	err := r.pool.QueryRow(ctx, query, discordID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *userRepository) Register(ctx context.Context, discordID string) (*domain.User, error) {
	query := `
		INSERT INTO users (discord_id)
		VALUES ($1)
		RETURNING id, discord_id, registered_at, updated_at
	`
	var user domain.User
	err := r.pool.QueryRow(ctx, query, discordID).Scan(
		&user.ID, &user.DiscordID, &user.RegisteredAt, &user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrUserAlreadyExists
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}
	return &user, nil
}
