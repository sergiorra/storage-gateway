package discovery_service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"storage-gateway/domain/ports"
	"storage-gateway/infrastructure/object-storage"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// DockerDiscoveryService represents a service for discovering Docker containers and extracting object storage information
type DockerDiscoveryService struct {
	c *client.Client
}

const (
	EnvKeyMinioAccessKey = "MINIO_ACCESS_KEY"
	EnvKeyMinioSecretKey = "MINIO_SECRET_KEY"
)

// NewDockerDiscoveryService creates a new instance of DockerDiscoveryService
func NewDockerDiscoveryService() (*DockerDiscoveryService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("new docker client: %w", err)
	}

	return &DockerDiscoveryService{
		c: cli,
	}, nil
}

// DiscoverNodes searches for Docker containers starting with the name "amazin-object-storage"
// and returns a list of object storage nodes for each Docker container discovered
func (dds *DockerDiscoveryService) DiscoverNodes(ctx context.Context) ([]ports.ObjectStorage, error) {
	containers, err := dds.c.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.Arg("name", "amazin-object-storage")),
	})
	if err != nil {
		return nil, err
	}

	var objectStorages []ports.ObjectStorage
	for _, c := range containers {
		containerInfo, err := dds.c.ContainerInspect(ctx, c.ID)
		if err != nil {
			return nil, err
		}

		var (
			accessKey string
			secretKey string
		)

		for _, envVar := range containerInfo.Config.Env {
			parts := strings.SplitN(envVar, "=", 2)
			if len(parts) == 2 {
				name, value := parts[0], parts[1]
				if name == EnvKeyMinioAccessKey {
					accessKey = value
				}
				if name == EnvKeyMinioSecretKey {
					secretKey = value
				}
			}
		}

		if accessKey == "" || secretKey == "" {
			return nil, errors.New("keys not found")
		}

		n, ok := c.NetworkSettings.Networks["storage-gateway_object-storage"]
		if !ok {
			return nil, errors.New("network not found")
		}

		// create a new MinioObjectStore instance for the discovered container
		node, err := object_storage.NewMinioObjectStore(ctx, c.ID, net.JoinHostPort(n.IPAddress, "9000"), accessKey, secretKey)
		if err != nil {
			return nil, err
		}

		objectStorages = append(objectStorages, node)
	}

	return objectStorages, nil
}
