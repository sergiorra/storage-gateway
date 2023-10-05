package ports

import (
	"context"
)

type DiscoveryService interface {
	DiscoverNodes(ctx context.Context) ([]ObjectStorage, error)
}
