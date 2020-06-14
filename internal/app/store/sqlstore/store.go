package sqlstore

import (
	"github.com/assanoff/http-crud-server/internal/app/store"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Store ...
type Store struct {
	db             *pgxpool.Pool
	userRepository *UserRepository
	schema         string
}

// New ...
func New(db *pgxpool.Pool, schema string) *Store {
	return &Store{
		db:     db,
		schema: schema,
	}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
