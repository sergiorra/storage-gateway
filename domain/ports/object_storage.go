package ports

import (
	"context"

	"storage-gateway/domain/models"
)

type ObjectStorage interface {
	GetObject(ctx context.Context, id string) (*models.Object, error)
	PutObject(ctx context.Context, o *models.Object) error
	ID() string
	IsOnline() bool
}
