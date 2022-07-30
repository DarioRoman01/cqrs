package database

import (
	"context"
	"database/sql"

	"github.com/DarioRoman01/cqrs/models"
	_ "github.com/lib/pq"
)

// UserRepository is a struct that implements the Repository interface
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(url string) (*UserRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &UserRepository{db: db}, nil
}

func (r *UserRepository) Close() error {
	return r.db.Close()
}

func (r *UserRepository) Insert(ctx context.Context, user *models.User) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO users (id, name, password, created_at) VALUES ($1, $2, $3, $4)",
		user.ID, user.Name, user.Password, user.CreatedAt,
	)

	return err
}

func (r *UserRepository) List(ctx context.Context) ([]*models.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Password, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepository) Get(ctx context.Context, id string) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Password, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET name = $1, password = $2, created_at = $3 WHERE id = $4",
		user.Name, user.Password, user.CreatedAt, user.ID,
	)

	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}

func (r *UserRepository) Login(ctx context.Context, name, password string) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE name = $1", name)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Password, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}
