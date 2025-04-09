package metadata

import (
	"context"
	"errors"

	"github.com/Maksim-Kot/Movie-application/metadata/internal/repository"
	"github.com/Maksim-Kot/Movie-application/metadata/pkg/model"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
	Put(ctx context.Context, id string, m *model.Metadata) error
}

// Controller defines a metadata service controller.
type Controller struct {
	repo metadataRepository
}

// New creates a metadata service controller.
func New(repo metadataRepository) *Controller {
	return &Controller{repo}
}

// Get returns movie metadata by id.
func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}

// Put writes movie metadata to repository.
func (c *Controller) Put(ctx context.Context, m *model.Metadata) error {
	return c.repo.Put(ctx, m.ID, m)
}
