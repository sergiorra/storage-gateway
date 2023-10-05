package object_storage

import (
	"context"

	"storage-gateway/domain/models"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	// bucketName represents the default name of the MinIO bucket
	bucketName = "object-store"
)

// MinioObjectStore represents a MinIO object storage instance
type MinioObjectStore struct {
	id string
	c  *minio.Client
}

// NewMinioObjectStore creates a new MinioObjectStore instance with the provided information.
// It establishes a connection to the MinIO server, creates the storage
// bucket if it doesn't exist, and returns the initialized MinioObjectStore
func NewMinioObjectStore(ctx context.Context, id, endpoint, accessKeyID, secretAccessKey string) (*MinioObjectStore, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	mos := &MinioObjectStore{
		id: id,
		c:  client,
	}

	if err = mos.createStorage(ctx); err != nil {
		return nil, err
	}

	return mos, nil
}

// createStorage checks if the default storage bucket exists and creates it if not
func (mos *MinioObjectStore) createStorage(ctx context.Context) error {
	exists, err := mos.c.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return mos.c.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
}

// PutObject stores an object in the MinIO bucket with the provided object metadata
func (mos *MinioObjectStore) PutObject(ctx context.Context, o *models.Object) error {
	_, err := mos.c.PutObject(ctx, bucketName, o.ID.Value(), o.Content, o.Size, minio.PutObjectOptions{ContentType: o.ContentType})
	if err != nil {
		return err
	}

	return nil
}

// GetObject retrieves an object from the MinIO bucket by its name and returns the associated object metadata
func (mos *MinioObjectStore) GetObject(ctx context.Context, name string) (*models.Object, error) {
	object, err := mos.c.GetObject(ctx, bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	objStat, err := object.Stat()
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, models.ErrObjectNotFound
		}
		return nil, err
	}

	return &models.Object{
		ID:          models.ObjectID(name),
		Content:     object,
		ContentType: objStat.ContentType,
		Size:        objStat.Size,
	}, nil
}

// ID returns the unique identifier associated with the MinioObjectStore
func (mos *MinioObjectStore) ID() string {
	return mos.id
}

// IsOnline checks if the MinioObjectStore is online and available
func (mos *MinioObjectStore) IsOnline() bool {
	return mos.c.IsOnline()
}
