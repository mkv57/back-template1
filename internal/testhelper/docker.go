package testhelper

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

func runContainer(
	ctx context.Context,
	t *testing.T,
	assert *require.Assertions,
	opt *dockertest.RunOptions,
	mappedPort string,
	endpoint string,
	cb func(int) error,
) error {
	t.Helper()

	pool, err := dockertest.NewPool(endpoint)
	assert.NoError(err)

	resource, err := pool.RunWithOptions(opt, func(cfg *docker.HostConfig) {
		cfg.AutoRemove = true
	})
	assert.NoError(err)

	deadline, ok := ctx.Deadline()
	expire := MaxTimeout
	if ok {
		expire = time.Since(deadline) * -1
	}

	err = resource.Expire(uint(expire.Seconds()))
	assert.NoError(err)

	resourceTCPPort := resource.GetPort(mappedPort)
	port, err := strconv.Atoi(resourceTCPPort)
	assert.NoError(err)
	t.Cleanup(func() {
		assert.NoError(pool.Purge(resource))
	})

	return pool.Retry(func() error {
		return cb(port)
	})
}
