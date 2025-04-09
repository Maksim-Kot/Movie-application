package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Maksim-Kot/Movie-application/metadata/internal/repository"
	"github.com/Maksim-Kot/Movie-application/metadata/pkg/model"

	_ "github.com/go-sql-driver/mysql"
)

// Repository defines a MySQL-based movie matadata repository.
type Repository struct {
	db *sql.DB
}

// New creates a new MySQL-based repository with credentials.
func New(credentials string) (*Repository, error) {
	db, err := sql.Open("mysql", credentials)
	if err != nil {
		return nil, err
	}

	return &Repository{db}, nil
}

// Get retrieves movie metadata for by movie id.
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var title, description, director string
	row := r.db.QueryRowContext(ctx, "SELECT title, description, director FROM movies WHERE id = ?", id)
	if err := row.Scan(&title, &description, &director); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &model.Metadata{
		ID:          id,
		Title:       title,
		Description: description,
		Director:    director,
	}, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	args := []any{id, metadata.Title, metadata.Description, metadata.Director}
	_, err := r.db.ExecContext(ctx, "INSERT INTO movies (id, title, description, director) VALUES (?, ?, ?, ?)", args)
	return err
}
