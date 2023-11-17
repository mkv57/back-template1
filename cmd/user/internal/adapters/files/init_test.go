//go:build integration

package files_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"github.com/ZergsLaw/back-template/cmd/user/internal/adapters/files"
	"github.com/ZergsLaw/back-template/internal/logger"
	"github.com/ZergsLaw/back-template/internal/testhelper"
)

func start(t *testing.T) (context.Context, *files.Client, *require.Assertions) {
	t.Helper()
	ctx := testhelper.Context(t)
	assert := require.New(t)

	namespace := testhelper.Namespace(t)

	username := "test_svc"
	pass := "test_pass"
	endpoint := testhelper.Minio(
		ctx,
		t,
		assert,
		username, pass, "", false,
		"local-1",
	)

	reg := prometheus.NewPedanticRegistry()
	fileStorage, err := files.New(ctx, reg, namespace, files.Config{
		Secure:    false,
		Endpoint:  endpoint,
		AccessKey: username,
		SecretKey: pass,
		Region:    "local-1",
	})
	assert.NoError(err)
	t.Cleanup(func() {
		assert.NoError(fileStorage.Close())
	})

	devLogger := slog.New(slog.NewJSONHandler(
		os.Stderr, &slog.HandlerOptions{ //nolint:exhaustruct
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	)
	return logger.NewContext(ctx, devLogger), fileStorage, assert
}
