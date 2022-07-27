package database

import (
	"context"
	"database/sql"

	"github.com/DarioRoman01/cqrs/models"
	_ "github.com/lib/pq"
)

// PostgresRepository is a struct that implements the Repository interface
type PostgresRepository struct {
	// db is the database connection
	db *sql.DB
}

// NewPostgresRepository creates a new PostgresRepository
func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

// Close closes the repository connection to the database
func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

// Insert inserts a new feed into the repository
func (r *PostgresRepository) Insert(ctx context.Context, feed *models.Feed) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO feeds (id, title, description, created_at) VALUES ($1, $2, $3)",
		feed.ID, feed.Title, feed.Description, feed.CreatedAt,
	)

	return err
}

// List returns all feeds in the repository
func (r *PostgresRepository) List(ctx context.Context) ([]*models.Feed, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM feeds")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feeds []*models.Feed
	for rows.Next() {
		var feed models.Feed
		if err := rows.Scan(&feed.ID, &feed.Title, &feed.Description, &feed.CreatedAt); err != nil {
			return nil, err
		}

		feeds = append(feeds, &feed)
	}

	return feeds, nil
}

// Get returns a feed from the repository by id
func (r *PostgresRepository) Get(ctx context.Context, id string) (*models.Feed, error) {
	var feed models.Feed
	err := r.db.QueryRowContext(ctx, "SELECT * FROM feeds WHERE id = $1", id).Scan(&feed.ID, &feed.Title, &feed.Description, &feed.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

// Delete deletes a feed from the repository by id
func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM feeds WHERE id = $1", id)
	return err
}

// Update updates a feed in the repository
func (r *PostgresRepository) Update(ctx context.Context, feed *models.Feed) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE feeds SET title = $1, description = $2 WHERE id = $3",
		feed.Title, feed.Description, feed.ID,
	)

	return err
}
