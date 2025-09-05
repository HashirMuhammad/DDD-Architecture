package repository

import (
	"context"
	"ddd-user-service/internal/domain"
	"strings"
	"sync"
)

type MemoryUserRepository struct {
	users map[domain.UserID]*domain.User
	mutex sync.RWMutex
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[domain.UserID]*domain.User),
	}
}

func (r *MemoryUserRepository) Save(ctx context.Context, user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepository) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	userCopy := *user
	return &userCopy, nil
}

func (r *MemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	email = strings.ToLower(strings.TrimSpace(email))
	for _, user := range r.users {
		if user.Email == email {
			userCopy := *user
			return &userCopy, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

func (r *MemoryUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	username = strings.ToLower(strings.TrimSpace(username))
	for _, user := range r.users {
		if user.Username == username {
			userCopy := *user
			return &userCopy, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

func (r *MemoryUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		userCopy := *user
		users = append(users, &userCopy)
	}

	return users, nil
}

func (r *MemoryUserRepository) Update(ctx context.Context, user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.users[user.ID]
	if !exists {
		return domain.ErrUserNotFound
	}

	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepository) Delete(ctx context.Context, id domain.UserID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.users[id]
	if !exists {
		return domain.ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

func (r *MemoryUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	email = strings.ToLower(strings.TrimSpace(email))
	for _, user := range r.users {
		if user.Email == email {
			return true, nil
		}
	}

	return false, nil
}

func (r *MemoryUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	username = strings.ToLower(strings.TrimSpace(username))
	for _, user := range r.users {
		if user.Username == username {
			return true, nil
		}
	}

	return false, nil
}
