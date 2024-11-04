// Package repo contains implements for app.Repo.
// Provide session info to and from repository.
package repo

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sipki-tech/database"
	"github.com/sipki-tech/database/connectors"
	"github.com/sipki-tech/database/migrations"

	"github.com/ZergsLaw/back-template1/cmd/session/internal/app"
	"github.com/ZergsLaw/back-template1/internal/dom"
)

var _ app.Repo = &Repo{}

type (
	// Config provide connection info for database.
	Config struct {
		Cockroach  connectors.CockroachDB
		MigrateDir string
		Driver     string
	}

	// Repo provided data from and to database.
	Repo struct {
		sql *database.SQL
	}

	session struct {
		ID        uuid.UUID `db:"id"`
		Token     string    `db:"token"`
		IP        string    `db:"ip"`
		UserAgent string    `db:"user_agent"`
		UserID    uuid.UUID `db:"user_id"`
		Status    string    `db:"status"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
)

const (
	requestUpdateStatus = `StatusUpdate`
)

// New build and returns session db.
func New(ctx context.Context, reg *prometheus.Registry, namespace string, cfg Config) (*Repo, error) {
	const subsystem = "repo"
	m := database.NewMetrics(reg, namespace, subsystem, new(app.Repo))

	returnErrs := []error{ // List of app.Err… returned by Repo methods.
		app.ErrNotFound,
		app.ErrDuplicate,
	}

	migrates, err := migrations.Parse(cfg.MigrateDir)
	if err != nil {
		return nil, fmt.Errorf("migrations.Parse: %w", err)
	}

	err = migrations.Run(ctx, cfg.Driver, &cfg.Cockroach, migrations.Up, migrates)
	if err != nil {
		return nil, fmt.Errorf("migrations.Run: %w", err)
	}

	conn, err := database.NewSQL(ctx, cfg.Driver, database.SQLConfig{
		Metrics:    m,
		ReturnErrs: returnErrs,
	}, &cfg.Cockroach)
	if err != nil {
		return nil, fmt.Errorf("librepo.NewCockroach: %w", err)
	}

	return &Repo{
		sql: conn,
	}, nil
}

// Close implements io.Closer.
func (r *Repo) Close() error {
	return r.sql.Close()
}

func convert(s app.Session) *session {
	return &session{
		ID:        s.ID,
		Token:     s.Token.Value,
		IP:        s.Origin.IP.String(),
		UserAgent: s.Origin.UserAgent,
		UserID:    s.UserID,
		Status:    s.Status.String(),
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func (s session) convert() *app.Session {
	return &app.Session{
		ID: s.ID,
		Origin: app.Origin{
			IP:        net.ParseIP(s.IP),
			UserAgent: s.UserAgent,
		},
		Token: app.Token{
			Value: s.Token,
		},
		UserID:    s.UserID,
		Status:    toUserStatus(s.Status),
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func toUserStatus(text string) dom.UserStatus {
	switch text {
	case dom.UserStatusFreeze.String():
		return dom.UserStatusFreeze
	case dom.UserStatusDefault.String():
		return dom.UserStatusDefault
	case dom.UserStatusPremium.String():
		return dom.UserStatusPremium
	case dom.UserStatusSupport.String():
		return dom.UserStatusSupport
	case dom.UserStatusAdmin.String():
		return dom.UserStatusAdmin
	case dom.UserStatusJedi.String():
		return dom.UserStatusJedi
	default:
		panic(fmt.Sprintf("unknown status: %s", text))
	}
}

// Save for implements app.Repo.
func (r *Repo) Save(ctx context.Context, session app.Session) error {
	return r.sql.NoTx(func(db *sqlx.DB) error {
		newSession := convert(session)

		const query = `
		insert into 
		sessions 
		    (id, token, ip, user_agent, user_id, status) 
		values 
			($1, $2, $3, $4, $5, $6)`

		_, err := db.ExecContext(ctx, query, newSession.ID, newSession.Token, newSession.IP, newSession.UserAgent, newSession.UserID, newSession.Status)
		if err != nil {
			return fmt.Errorf("db.ExecContext: %w", err)
		}

		return nil
	})
}

// ByID for implements app.Repo.
func (r *Repo) ByID(ctx context.Context, sessionID uuid.UUID) (s *app.Session, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `select * from sessions where id = $1`

		res := session{}
		err = db.GetContext(ctx, &res, query, sessionID)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}

		s = res.convert()

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Delete for implements app.Repo.
func (r *Repo) Delete(ctx context.Context, sessionID uuid.UUID) error {
	return r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `
		delete
		from sessions
		where id = $1 returning *`

		err := db.GetContext(ctx, &session{}, query, sessionID)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}

		return nil
	})
}

func (r *Repo) UpdateStatus(ctx context.Context, reqID, userID uuid.UUID, status dom.UserStatus) error {
	return r.sql.Tx(ctx, nil, func(tx *sqlx.Tx) error {
		err := insertToDeduplication(ctx, tx, reqID, requestUpdateStatus)
		if err != nil {
			return fmt.Errorf("r.insertToDeduplication: %w", convertErr(err))
		}

		const query = `update sessions set status = $1 where user_id = $2`

		_, err = tx.ExecContext(ctx, query, status.String(), userID)
		if err != nil {
			return fmt.Errorf("db.ExecContext: %w", convertErr(err))
		}

		return nil
	})
}

func insertToDeduplication(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, kind string) error {
	const query = `insert into deduplication (id, kind) values ($1, $2) returning id`

	return tx.GetContext(ctx, &uuid.UUID{}, query, id, kind)
}
