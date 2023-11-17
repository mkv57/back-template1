package testhelper

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"

	"github.com/ZergsLaw/back-template/internal/queue"
)

// Default values for making nats test container.
const (
	NATSImage           = `nats`
	NATSVersion         = `2.10.5`
	NATSDefaultEndpoint = ``
	NATSDefaultHost     = `localhost`
)

// NATS build and run test container with message broker.
//
// Notice:
//   - Starts NATS container;
//   - Make connections by credentials;
func NATS(
	ctx context.Context,
	t *testing.T,
	assert *require.Assertions,
	cfgFilePath string,
	username, password string,
) string {
	t.Helper()

	opt := &dockertest.RunOptions{
		Repository: NATSImage,
		Tag:        NATSVersion,
		Cmd:        []string{"-c=/srv.conf"},
		Mounts: []string{
			fmt.Sprintf("%s:/srv.conf", cfgFilePath),
		},
		//PortBindings: map[docker.Port][]docker.PortBinding{
		//	"4222/tcp": {{
		//		HostIP:   "localhost",
		//		HostPort: fmt.Sprintf("%d", UnusedTCPPort(t, assert, Host)),
		//	}},
		//},
	}

	appName := t.Name()

	addr := ""
	err := runContainer(ctx, t, assert, opt, "4222/tcp", NATSDefaultEndpoint, func(port int) (err error) {
		defer func() {
			if err != nil {
				t.Logf("connection problem: %s", err)
			}
		}()

		url := fmt.Sprintf("nats://%s", net.JoinHostPort(NATSDefaultHost, strconv.Itoa(port)))
		q, err := queue.Connect(ctx, url, appName, username, password)
		if err != nil {
			return fmt.Errorf("queue.Connect: %w", err)
		}
		t.Cleanup(func() {
			assert.NoError(q.Drain())
		})

		addr = url

		return nil
	})
	assert.NoError(err)

	return addr
}
