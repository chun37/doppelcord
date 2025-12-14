package cached

import (
	"context"
	"sync"

	"github.com/chun37/doppelcord/internal/domain"
	"github.com/chun37/doppelcord/internal/repository"
)

type CachedUserRepository struct {
	inner      repository.UserRepository
	registered map[string]struct{}
	mu         sync.RWMutex
}

func NewCachedUserRepository(inner repository.UserRepository) *CachedUserRepository {
	return &CachedUserRepository{
		inner:      inner,
		registered: make(map[string]struct{}),
	}
}

func (r *CachedUserRepository) LoadAll(ctx context.Context) error {
	ids, err := r.inner.GetAllDiscordIDs(ctx)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.registered = make(map[string]struct{}, len(ids))
	for _, id := range ids {
		r.registered[id] = struct{}{}
	}
	return nil
}

func (r *CachedUserRepository) IsRegistered(ctx context.Context, discordID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.registered[discordID]
	return exists, nil
}

func (r *CachedUserRepository) Register(ctx context.Context, discordID string) (*domain.User, error) {
	user, err := r.inner.Register(ctx, discordID)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	r.registered[discordID] = struct{}{}
	r.mu.Unlock()

	return user, nil
}

func (r *CachedUserRepository) GetAllDiscordIDs(ctx context.Context) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.registered))
	for id := range r.registered {
		ids = append(ids, id)
	}
	return ids, nil
}
