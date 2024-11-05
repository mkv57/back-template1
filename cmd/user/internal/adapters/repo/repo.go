// Package repo contains wrapper for database abstraction.
package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sipki-tech/database"
	"github.com/sipki-tech/database/connectors"
	"github.com/sipki-tech/database/migrations"

	"github.com/ZergsLaw/back-template1/cmd/user/internal/app"
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
)

// New build and returns user db.
func New(ctx context.Context, reg *prometheus.Registry, namespace string, cfg Config) (*Repo, error) {
	const subsystem = "repo"
	m := database.NewMetrics(reg, namespace, subsystem, new(app.Repo))

	returnErrs := []error{ // List of app.Err… returned by Repo methods.
		app.ErrNotFound,
		app.ErrUsernameExist,
		app.ErrEmailExist,
		app.ErrUserIDAndFileIDExist,
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

// Save for implements app.Repo.
func (r *Repo) Save(ctx context.Context, u app.User) (id uuid.UUID, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		newUser := convert(u)
		const query = `
		insert into 
		users 
		    (email, name, full_name, pass_hash, status) 
		values 
			($1, $2, $3, $4, $5)
		returning id
		`

		err := db.GetContext(ctx, &id, query, newUser.Email, newUser.Name, newUser.FullName, newUser.PassHash, newUser.Status)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// Update for implements app.Repo.
func (r *Repo) Update(ctx context.Context, u app.User) (upUser *app.User, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		updateUser := convert(u)
		const query = `
		update users
		set
			email 	  		  = $1,
			name  	    	  = $2,
			full_name  		  = $3,
			pass_hash   	  = $4,
			current_avatar_id = $5,
			status 			  = $6,
			updated_at = now()
		where id = $7
		returning *`

		var res user
		err := db.GetContext(ctx, &res, query, updateUser.Email, updateUser.Name, updateUser.FullName, updateUser.PassHash,
			updateUser.CurrentAvatarID, updateUser.Status, updateUser.ID)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}

		upUser = res.convert()

		return nil
	})
	if err != nil {
		return nil, err
	}

	return upUser, nil
}

// Delete for implements app.Repo.
func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `
		delete
		from users
		where id = $1 returning *`

		err := db.GetContext(ctx, &user{}, query, id)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}

		return nil
	})
}

// ByID for implements app.Repo.
func (r *Repo) ByID(ctx context.Context, id uuid.UUID) (u *app.User, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `select * from users where id = $1`

		res := user{}
		err = db.GetContext(ctx, &res, query, id)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}

		u = res.convert()

		return nil
	})
	if err != nil {
		return nil, err
	}

	return u, nil
}

// ByEmail for implements app.Repo.
func (r *Repo) ByEmail(ctx context.Context, email string) (u *app.User, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `select * from users where email = $1`

		res := user{}
		err = db.GetContext(ctx, &res, query, email)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}

		u = res.convert()

		return nil
	})
	if err != nil {
		return nil, err
	}

	return u, nil
}

// ByUsername for implements app.Repo.
func (r *Repo) ByUsername(ctx context.Context, username string) (u *app.User, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `select * from users where name = $1`

		res := user{}
		err = db.GetContext(ctx, &res, query, username)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}

		u = res.convert()

		return nil
	})
	if err != nil {
		return nil, err
	}

	return u, nil
}

// SearchUsers for implements app.Repo.
func (r *Repo) SearchUsers(ctx context.Context, params app.SearchParams) (users []app.User, total int, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		query, args, err := getUsers(params)
		if err != nil {
			return fmt.Errorf("getUsers: %w", err)
		}

		res := make([]user, 0, params.Limit)
		err = db.SelectContext(ctx, &res, query, args...)
		if err != nil {
			return fmt.Errorf("db.SelectContext: %w", convertErr(err))
		}

		if len(res) == 0 {
			return nil
		}

		totalQuery, totalArgs, err := getTotalUsers(params)
		if err != nil {
			return fmt.Errorf("getTotalUsers: %w", err)
		}

		err = db.GetContext(ctx, &total, totalQuery, totalArgs...)
		switch {
		case errors.Is(err, sql.ErrNoRows):
		case err != nil:
			return fmt.Errorf("db.GetContext: %w", err)
		}

		users = make([]app.User, len(res))
		for i := range res {
			users[i] = *res[i].convert()
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// SaveAvatar for implements app.Repo.
func (r *Repo) SaveAvatar(ctx context.Context, userFile app.AvatarInfo) (err error) {
	return r.sql.NoTx(func(db *sqlx.DB) error {
		avatarCache := convertUserFile(userFile)
		const query = `
		insert into 
		avatars 
		    (owner_id, id) 
		values 
			($1, $2)
		`

		_, err := db.ExecContext(ctx, query, avatarCache.OwnerID, avatarCache.ID)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}

		return nil
	})
}

// DeleteAvatar for implements app.Repo.
func (r *Repo) DeleteAvatar(ctx context.Context, userID, avatarID uuid.UUID) error {
	return r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `
		delete
		from avatars
		where owner_id = $1
		and id = $2`

		_, err := db.ExecContext(ctx, query, userID, avatarID)
		if err != nil {
			return fmt.Errorf("db.ExecContext: %w", convertErr(err))
		}

		return nil
	})
}

// GetAvatar for implements app.Repo.
func (r *Repo) GetAvatar(ctx context.Context, avatarID uuid.UUID) (f *app.AvatarInfo, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `select * from avatars where id = $1`

		res := avatar{}
		err = db.GetContext(ctx, &res, query, avatarID)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}

		f = res.convert()

		return nil
	})
	if err != nil {
		return nil, err
	}

	return f, nil
}

// ListAvatarByUserID for implements app.Repo.
func (r *Repo) ListAvatarByUserID(ctx context.Context, userID uuid.UUID) (userAvatars []app.AvatarInfo, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `select * from avatars where owner_id = $1 order by created_at desc`

		var res []avatar
		err = db.SelectContext(ctx, &res, query, userID)
		if err != nil {
			return fmt.Errorf("db.SelectContext: %w", convertErr(err))
		}

		userAvatars = make([]app.AvatarInfo, len(res))
		for i := range res {
			userAvatars[i] = *res[i].convert()
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return userAvatars, nil
}

// GetCountAvatars for implements app.Repo.
func (r *Repo) GetCountAvatars(ctx context.Context, ownerID uuid.UUID) (total int, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const getTotal = `select count(*) over() as total from avatars where owner_id = $1`

		err = db.GetContext(ctx, &total, getTotal, ownerID)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}

// SaveTask implements app.Repo.
func (r *Repo) SaveTask(ctx context.Context, task app.Task) (id uuid.UUID, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		newTask, err := convertTask(task)
		if err != nil {
			return fmt.Errorf("convertTask: %w", err)
		}

		const query = `
		insert into 
		tasks 
		    (user_bytes, kind) 
		values
			($1, $2)
		returning id
		`

		err = db.GetContext(ctx, &id, query, newTask.UserBytes, newTask.Kind)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// FinishTask implements app.Repo.
func (r *Repo) FinishTask(ctx context.Context, id uuid.UUID) error {
	return r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `
		update tasks set 
		updated_at = now(),
		finished_at = now()
		where id = $1`

		_, err := db.ExecContext(ctx, query, id)
		if err != nil {
			return fmt.Errorf("db.ExecContext: %w", convertErr(err))
		}

		return nil
	})
}

// ListActualTask implements app.Repo.
func (r *Repo) ListActualTask(ctx context.Context, limit int) (tasks []app.Task, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `select * from tasks where finished_at is null order by created_at asc limit $1`

		res := make([]task, 0, limit)
		err = db.SelectContext(ctx, &res, query, limit)
		if err != nil {
			return fmt.Errorf("db.SelectContext: %w", convertErr(err))
		}

		tasks = make([]app.Task, len(res))
		for i := range res {
			t, err := res[i].convert()
			if err != nil {
				return fmt.Errorf("convert: %w", err)
			}

			tasks[i] = *t
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repo) UsersByIDs(ctx context.Context, ids []uuid.UUID) (users []app.User, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `select * from users where id = any($1)`

		res := make([]user, 0, len(ids))

		err = db.SelectContext(ctx, &res, query, pq.Array(ids))
		if err != nil {
			return fmt.Errorf("db.SelectContext: %w", convertErr(err))
		}

		users = make([]app.User, len(res))
		for i := range res {
			users[i] = *res[i].convert()
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Tx implements app.Repo.
func (r *Repo) Tx(ctx context.Context, f func(app.Repo) error) error {
	opt := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	return r.sql.Tx(ctx, opt, func(tx *sqlx.Tx) error {
		return f(&txRepo{tx: tx})
	})
}
