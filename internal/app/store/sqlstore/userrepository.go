package sqlstore

import (
	"context"
	"fmt"

	"github.com/assanoff/http-crud-server/internal/app/model"
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

	// return r.store.db.QueryRow(
	// 	"INSERT INTO users (name, email VALUES ($1, $2) RETURNING id",
	// 	u.Name,
	// 	u.Email,
	// ).Scan(&u.ID)

	row := conn.QueryRow(context.Background(),
		"INSERT INTO test.users (name, email) VALUES ($1, $2) RETURNING id",
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

	row := conn.QueryRow(context.Background(),
		"SELECT id, name, email FROM test.users WHERE id = $1",
		id)

	u := &model.User{}
	err = row.Scan(&u.ID, &u.Name, &u.Email)
	if err == pgx.ErrNoRows {
		return nil, err
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
	query := fmt.Sprintf("SELECT id, name, email FROM test.users WHERE %s = $1", fieldName)

	row := conn.QueryRow(context.Background(),
		query,
		value)

	u := &model.User{}
	err = row.Scan(&u.ID, &u.Name, &u.Email)
	if err == pgx.ErrNoRows {
		return nil, err
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
