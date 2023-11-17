package testhelper

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

// Default values for making nats test container.
const (
	MinioImage           = `docker.io/bitnami/minio`
	MinioVersion         = `2023`
	MinioDefaultEndpoint = ``
	MinioDefaultHost     = `localhost`
)

// Minio build and run test container with file storage.
//
// Notice:
//   - Starts Minio container;
//   - Make connections by credentials;
func Minio(
	ctx context.Context,
	t *testing.T,
	assert *require.Assertions,
	accessKey, secretKey, sessionToken string,
	secure bool,
	region string,
) string {
	t.Helper()

	opt := &dockertest.RunOptions{
		Repository: MinioImage,
		Tag:        MinioVersion,
		Env: []string{
			"MINIO_ROOT_USER=" + accessKey,
			"MINIO_ROOT_PASSWORD=" + secretKey,
		},
	}

	endpoint := ""
	err := runContainer(ctx, t, assert, opt, "9000/tcp", MinioDefaultEndpoint, func(port int) (err error) {
		defer func() {
			if err != nil {
				t.Logf("connection problem: %s", err)
			}
		}()

		transport, err := minio.DefaultTransport(secure)
		if err != nil {
			return fmt.Errorf("minio.DefaultTransport: %w", err)
		}

		opts := &minio.Options{
			Creds:        credentials.NewStaticV4(accessKey, secretKey, sessionToken),
			Secure:       secure,
			Transport:    transport,
			Region:       region,
			BucketLookup: minio.BucketLookupAuto,
		}
		endpoint = net.JoinHostPort(MinioDefaultHost, strconv.Itoa(port))

		client, err := minio.New(endpoint, opts)
		if err != nil {
			return fmt.Errorf("minio.New: %w", err)
		}

		_, err = client.BucketExists(ctx, "test.bucket")
		if err != nil {
			return fmt.Errorf("client.HealthCheck: %w", err)
		}

		return nil
	})
	assert.NoError(err)

	return endpoint
}
