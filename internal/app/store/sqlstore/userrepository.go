package sqlstore

import (
	"context"
	"fmt"

	"github.com/assanoff/http-crud-server/internal/app/model"
	"github.com/assanoff/http-crud-server/internal/app/store"
	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

// Create ...
func (r *UserRepository) Create(u *model.User) error {

	pool := r.store.db
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v", err)
		return err
	}
	defer conn.Release()

	query := fmt.Sprintf("INSERT INTO %s.users (name, email) VALUES ($1, $2) RETURNING id", r.store.schema)

	row := conn.QueryRow(context.Background(),
		query,
		u.Name, u.Email)
	var id uint64
	err = row.Scan(&id)

	if err != nil {
		log.Errorf("Unable to INSERT: %v", err)
		return err
	}
	return nil
}

// GetUserByID ...
func (r *UserRepository) GetUserByID(id int) (*model.User, error) {

	pool := r.store.db

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v\n", err)
		return nil, err
	}
	defer conn.Release()
	query := fmt.Sprintf("SELECT id, name, email FROM %s.users WHERE id = $1", r.store.schema)

	row := conn.QueryRow(context.Background(),
		query,
		id)

	u := &model.User{}
	err = row.Scan(&u.ID, &u.Name, &u.Email)

	if err == pgx.ErrNoRows {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}

// GetUserByField ...
func (r *UserRepository) GetUserByField(fieldName string, value string) (*model.User, error) {

	pool := r.store.db

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v\n", err)
		return nil, err
	}
	defer conn.Release()
	query := fmt.Sprintf("SELECT id, name, email FROM %s.users WHERE %s = $1", r.store.schema, fieldName)

	row := conn.QueryRow(context.Background(),
		query,
		value)

	u := &model.User{}
	err = row.Scan(&u.ID, &u.Name, &u.Email)
	if err == pgx.ErrNoRows {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}

// GetUsers ...
func (r *UserRepository) GetUsers() ([]*model.User, error) {

	pool := r.store.db
	var users []*model.User

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v\n", err)
		return nil, err
	}
	defer conn.Release()
	query := "SELECT id, name, email FROM test.users"

	rows, err := conn.Query(context.Background(), query)
	defer rows.Close()

	if err != nil {
		log.Errorf("could not get user query")
	}

	for rows.Next() {
		u := &model.User{}
		err = rows.Scan(&u.ID, &u.Name, &u.Email)
		if err == pgx.ErrNoRows {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// UpdateUserByID ...
func (r *UserRepository) UpdateUserByID(id int, u *model.User) (*model.User, error) {

	pool := r.store.db

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v\n", err)
		return nil, err
	}
	defer conn.Release()
	query := fmt.Sprintf("UPDATE %s.users SET name = $2, email = $3 WHERE id = $1", r.store.schema)

	ct, err := conn.Exec(context.Background(), query,
		id, u.Name, u.Email)

	if ct.RowsAffected() == 0 {
		return nil, nil
	}

	return u, nil
}

// DeleteUserByID ...
func (r *UserRepository) DeleteUserByID(id int) error {

	pool := r.store.db

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Errorf("Unable to acquire a database connection: %v\n", err)
		return err
	}
	defer conn.Release()
	query := fmt.Sprintf("DELETE FROM %s.users WHERE id = $1", r.store.schema)

	ct, err := conn.Exec(context.Background(), query,
		id)

	if ct.RowsAffected() == 0 {
		return nil
	}

	return nil
}
