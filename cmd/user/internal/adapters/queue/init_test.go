package queue_test

import (
	"context"
	"github.com/ZergsLaw/back-template/cmd/user/internal/adapters/queue"
	"github.com/ZergsLaw/back-template/internal/logger"
	que "github.com/ZergsLaw/back-template/internal/queue"
	"github.com/ZergsLaw/back-template/internal/testhelper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	queueCfgPath = `testdata/nats.conf`
)

func start(t *testing.T) (context.Context, *queue.Client, *require.Assertions, *que.Queue) {
	t.Helper()

	ctx := testhelper.Context(t)
	assert := require.New(t)

	pwd, err := os.Getwd()
	assert.NoError(err)

	namespace := testhelper.Namespace(t)

	username := "test_svc"
	password := "test_pass"

	reg := prometheus.NewPedanticRegistry()

	urlPath := testhelper.NATS(ctx, t, assert, filepath.Join(pwd, queueCfgPath), username, password)

	cfg := queue.Config{
		URLs:     []string{urlPath},
		Username: username,
		Password: password,
	}

	cliQ, err := que.Connect(ctx, strings.Join(cfg.URLs, ","), namespace, cfg.Username, cfg.Password)
	require.NoError(t, err)

	c, err := queue.New(ctx, reg, namespace, queue.Config{
		URLs:     []string{urlPath},
		Username: username,
		Password: password,
	})
	assert.NoError(err)

	t.Cleanup(func() {
		assert.NoError(c.Close())
	})

	log := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			},
		),
	)

	return logger.NewContext(ctx, log), c, assert, cliQ
}
