package testhelper

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/sipki-tech/database/connectors"
	"github.com/stretchr/testify/require"
)

// Default values for making cockroachDB test container.
const (
	CockroachDBImage           = `cockroachdb/cockroach`
	CockroachDBVersion         = `v23.1.12`
	CockroachDBDefaultEndpoint = ``
	CockroachDBDefaultHost     = `localhost`
)

// CockroachDB build and run test container with database.
//
// Notice:
//   - Starts CockroachDB container;
//   - Make connections by root certs;
//   - Make new database/user and grant access.
func CockroachDB(
	ctx context.Context,
	t *testing.T,
	assert *require.Assertions,
	caCrtPath, caKeyPath string,
	nodeCrtPath, nodeKeyPath string,
	clientCrtPath, clientKeyPath string,
) *connectors.CockroachDB {
	t.Helper()

	name := t.Name()
	appName := t.Name()
	dbName := fmt.Sprintf("%s_db", strings.ToLower(name))
	dbUser := fmt.Sprintf("%s_svc", strings.ToLower(name))
	dbPass := fmt.Sprintf("%s_pass", strings.ToLower(name))

	opt := &dockertest.RunOptions{
		Repository: CockroachDBImage,
		Tag:        CockroachDBVersion,
		Cmd:        []string{"start-single-node", "--certs-dir=/cockroach/certs/"},
		Mounts: []string{
			fmt.Sprintf("%s:/cockroach/certs/ca.crt", caCrtPath),
			fmt.Sprintf("%s:/cockroach/certs/ca.key", caKeyPath),
			fmt.Sprintf("%s:/cockroach/certs/client.root.crt", clientCrtPath),
			fmt.Sprintf("%s:/cockroach/certs/client.root.key", clientKeyPath),
			fmt.Sprintf("%s:/cockroach/certs/node.crt", nodeCrtPath),
			fmt.Sprintf("%s:/cockroach/certs/node.key", nodeKeyPath),
		},
		Env: []string{
			fmt.Sprintf("COCKROACH_DATABASE=%s", dbName),
			fmt.Sprintf("COCKROACH_USER=%s", dbUser),
			fmt.Sprintf("COCKROACH_PASSWORD=%s", dbPass),
		},
	}

	servicePort := 0
	err := runContainer(ctx, t, assert, opt, "26257/tcp", CockroachDBDefaultEndpoint, func(port int) (err error) {
		defer func() {
			if err != nil {
				t.Logf("connection problem: %s", err)
			}
		}()

		//cfg := connectors.CockroachDB{
		//	User:     "root1",
		//	Password: "root1",
		//	Host:     CockroachDBDefaultHost,
		//	Port:     port,
		//	Database: "defaultdb",
		//	Parameters: &connectors.CockroachDBParameters{
		//		Mode: connectors.CockroachSSLRequire,
		//		//SSLRootCert: caCrtPath,
		//		//SSLCert:     clientCrtPath,
		//		//SSLKey:      clientKeyPath,
		//	},
		//}
		//
		//db, err := database.NewSQL(ctx, "postgres", database.SQLConfig{}, &cfg)
		//if err != nil {
		//	return fmt.Errorf("database.NewSQL: %w", err)
		//}
		//t.Cleanup(func() {
		//	assert.NoError(db.Close())
		//})
		//
		//err = db.Tx(ctx, nil, func(tx *sqlx.Tx) error {
		//	return initDB(ctx, dbName, dbUser, dbPass, tx)
		//})
		//if err != nil {
		//	return fmt.Errorf("db.Tx: %w", err)
		//}
		servicePort = port

		return nil
	})
	assert.NoError(err)

	time.Sleep(time.Second * 20)

	return &connectors.CockroachDB{
		User:     dbUser,
		Password: dbPass,
		Host:     CockroachDBDefaultHost,
		Port:     servicePort,
		Database: dbName,
		Parameters: &connectors.CockroachDBParameters{
			ApplicationName: appName,
			Mode:            connectors.CockroachSSLRequire,
		},
	}
}

func initDB(ctx context.Context, dbName, dbUser, dbPass string, tx *sqlx.Tx) error {
	queryCreateDB := fmt.Sprintf("create database %s", dbName)
	queryCreateUser := fmt.Sprintf("create user %s with password '%s'", dbUser, dbPass)
	queryGrantAccess := fmt.Sprintf("grant all on database %s to %s", dbName, dbUser)

	_, err := tx.ExecContext(ctx, queryCreateDB)
	if err != nil {
		return fmt.Errorf("tx.ExecContext.CreateDB: %w", err)
	}
	_, err = tx.ExecContext(ctx, queryCreateUser)
	if err != nil {
		return fmt.Errorf("tx.ExecContext.CreateUser: %w", err)
	}
	_, err = tx.ExecContext(ctx, queryGrantAccess)
	if err != nil {
		return fmt.Errorf("tx.ExecContext.GrantAccess: %w", err)
	}

	return nil
}
