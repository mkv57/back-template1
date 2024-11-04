package files

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/minio/minio-go/v7"

	"github.com/ZergsLaw/back-template1/cmd/user/internal/app"
)

// UploadFile implements app.FileStore.
func (c *Client) UploadFile(ctx context.Context, f app.Avatar) (uuid.UUID, error) {
	id := uuid.Must(uuid.NewV4())

	const partSize = 1024 * 1024 / 2
	_, err := c.store.PutObject(ctx, bucketName, id.String(), f, f.Size, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			headerSrcName: f.Name,
		},
		ContentType: f.ContentType,
		PartSize:    partSize,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("c.store.PutObject: %w", err)
	}

	return id, nil
}

// DownloadFile implements app.FileStore.
func (c *Client) DownloadFile(ctx context.Context, id uuid.UUID) (*app.Avatar, error) {
	file, err := c.store.GetObject(ctx, bucketName, id.String(), minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("c.store.GetObject: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("file.stat: %w", err)
	}

	if stat.IsDeleteMarker {
		return nil, app.ErrNotFound
	}

	f := &app.Avatar{
		ReadSeekCloser: file,
		ID:             id,
		Name:           stat.Metadata.Get(headerSrcName),
		Size:           stat.Size,
		ModTime:        stat.LastModified,
		ContentType:    stat.ContentType,
	}

	return f, nil
}

// DeleteFile implements app.FileStore.
func (c *Client) DeleteFile(ctx context.Context, id uuid.UUID) error {
	err := c.store.RemoveObject(ctx, bucketName, id.String(), minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("c.store.RemoveObject: %w", err)
	}

	return nil
}

// Close implements io.Closer.
func (*Client) Close() error {
	return nil
}
