//go:build integration

package repo_test

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"github.com/ZergsLaw/back-template/cmd/user/internal/adapters/repo"
	"github.com/ZergsLaw/back-template/internal/logger"
	"github.com/ZergsLaw/back-template/internal/testhelper"
)

const (
	migrateDir    = `../../../migrate`
	caCrtPath     = `../../../../../certs/cockroach/ca.crt`
	caKeyPath     = `../../../../../certs/cockroach/ca.key`
	nodeCrtPath   = `../../../../../certs/cockroach/nodes/node1/node.crt`
	nodeKeyPath   = `../../../../../certs/cockroach/nodes/node1/node.key`
	clientCrtPath = `../../../../../certs/cockroach/client.root.crt`
	clientKeyPath = `../../../../../certs/cockroach/client.root.key`
)

func start(t *testing.T) (context.Context, *repo.Repo, *require.Assertions) {
	t.Helper()
	ctx := testhelper.Context(t)
	assert := require.New(t)

	pwd, err := os.Getwd()
	assert.NoError(err)

	namespace := testhelper.Namespace(t)

	cockroachCfg := testhelper.CockroachDB(
		ctx,
		t,
		assert,
		filepath.Join(pwd, caCrtPath), filepath.Join(pwd, caKeyPath),
		filepath.Join(pwd, nodeCrtPath), filepath.Join(pwd, nodeKeyPath),
		filepath.Join(pwd, clientCrtPath), filepath.Join(pwd, clientKeyPath),
	)

	reg := prometheus.NewPedanticRegistry()
	r, err := repo.New(ctx, reg, namespace, repo.Config{
		Cockroach:  *cockroachCfg,
		MigrateDir: migrateDir,
		Driver:     "postgres",
	})
	assert.NoError(err)
	t.Cleanup(func() {
		assert.NoError(r.Close())
	})

	log := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{ //nolint:exhaustruct
				AddSource: true,
				Level:     slog.LevelDebug,
			},
		),
	)

	return logger.NewContext(ctx, log), r, assert
}
