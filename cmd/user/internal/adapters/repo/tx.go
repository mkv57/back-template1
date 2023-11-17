package repo

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/ZergsLaw/back-template/cmd/user/internal/app"
)

var _ app.Repo = &txRepo{}

type txRepo struct {
	tx *sqlx.Tx
}

// Save for implements app.Repo.
func (t *txRepo) Save(ctx context.Context, u app.User) (id uuid.UUID, err error) {
	newUser := convert(u)
	const query = `
		insert into 
		users 
		    (email, name, full_name, pass_hash, status) 
		values 
			($1, $2, $3, $4, $5)
		returning id
		`

	err = t.tx.GetContext(ctx, &id, query, newUser.Email, newUser.Name, newUser.FullName, newUser.PassHash, newUser.Status)
	if err != nil {
		return uuid.Nil, fmt.Errorf("db.GetContext: %w", convertErr(err))
	}

	return id, nil
}

// Update for implements app.Repo.
func (t *txRepo) Update(ctx context.Context, u app.User) (upUser *app.User, err error) {
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
	err = t.tx.GetContext(ctx, &res, query, updateUser.Email, updateUser.Name, updateUser.FullName, updateUser.PassHash,
		updateUser.CurrentAvatarID, updateUser.Status, updateUser.ID)
	if err != nil {
		return nil, fmt.Errorf("t.tx.GetContext: %w", convertErr(err))
	}

	upUser = res.convert()

	return upUser, nil
}

// Delete for implements app.Repo.
func (t *txRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `delete from users where id = $1 returning *`

	err := t.tx.GetContext(ctx, &user{}, query, id)
	if err != nil {
		return fmt.Errorf("db.GetContext: %w", convertErr(err))
	}

	return nil
}

// ByID for implements app.Repo.
func (t *txRepo) ByID(ctx context.Context, id uuid.UUID) (u *app.User, err error) {
	const query = `select * from users where id = $1`

	res := user{}
	err = t.tx.GetContext(ctx, &res, query, id)
	if err != nil {
		return nil, fmt.Errorf("db.GetContext: %w", convertErr(err))
	}

	return res.convert(), nil
}

// ByEmail for implements app.Repo.
func (t *txRepo) ByEmail(ctx context.Context, email string) (u *app.User, err error) {
	const query = `select * from users where email = $1`

	res := user{}
	err = t.tx.GetContext(ctx, &res, query, email)
	if err != nil {
		return nil, fmt.Errorf("db.GetContext: %w", convertErr(err))
	}

	return res.convert(), nil
}

// ByUsername for implements app.Repo.
func (t *txRepo) ByUsername(ctx context.Context, username string) (u *app.User, err error) {
	const query = `select * from users where name = $1`

	res := user{}
	err = t.tx.GetContext(ctx, &res, query, username)
	if err != nil {
		return nil, fmt.Errorf("db.GetContext: %w", convertErr(err))
	}

	u = res.convert()

	return u, nil
}

// SearchUsers for implements app.Repo.
func (t *txRepo) SearchUsers(ctx context.Context, params app.SearchParams) (users []app.User, total int, err error) {
	query, args, err := getUsers(params)
	if err != nil {
		return nil, 0, fmt.Errorf("getUsers: %w", err)
	}

	res := make([]user, 0, params.Limit)
	err = t.tx.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("db.SelectContext: %w", convertErr(err))
	}

	totalQuery, totalArgs, err := getTotalUsers(params)
	if err != nil {
		return nil, 0, fmt.Errorf("getTotalUsers: %w", err)
	}

	err = t.tx.GetContext(ctx, &total, totalQuery, totalArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("db.GetContext: %w", err)
	}

	users = make([]app.User, len(res))
	for i := range res {
		users[i] = *res[i].convert()
	}

	return users, total, nil
}

// SaveAvatar for implements app.Repo.
func (t *txRepo) SaveAvatar(ctx context.Context, userFile app.AvatarInfo) (err error) {
	avatarCache := convertUserFile(userFile)
	const query = `insert into avatars (owner_id, id) values ($1, $2)`

	_, err = t.tx.ExecContext(ctx, query, avatarCache.OwnerID, avatarCache.ID)
	if err != nil {
		return fmt.Errorf("db.GetContext: %w", convertErr(err))
	}

	return nil
}

// DeleteAvatar for implements app.Repo.
func (t *txRepo) DeleteAvatar(ctx context.Context, userID, avatarID uuid.UUID) error {
	const query = `delete from avatars	where owner_id = $1	and id = $2`

	_, err := t.tx.ExecContext(ctx, query, userID, avatarID)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", convertErr(err))
	}

	return nil
}

// GetAvatar for implements app.Repo.
func (t *txRepo) GetAvatar(ctx context.Context, avatarID uuid.UUID) (f *app.AvatarInfo, err error) {
	const query = `select * from avatars where id = $1`

	res := avatar{}
	err = t.tx.GetContext(ctx, &res, query, avatarID)
	if err != nil {
		return nil, fmt.Errorf("db.GetContext: %w", convertErr(err))
	}

	f = res.convert()

	return f, nil
}

// ListAvatarByUserID for implements app.Repo.
func (t *txRepo) ListAvatarByUserID(ctx context.Context, userID uuid.UUID) (userAvatars []app.AvatarInfo, err error) {
	const query = `select * from avatars where owner_id = $1 order by created_at desc`

	var res []avatar
	err = t.tx.SelectContext(ctx, &res, query, userID)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext: %w", convertErr(err))
	}

	userAvatars = make([]app.AvatarInfo, len(res))
	for i := range res {
		userAvatars[i] = *res[i].convert()
	}

	return userAvatars, nil
}

// GetCountAvatars for implements app.Repo.
func (t *txRepo) GetCountAvatars(ctx context.Context, ownerID uuid.UUID) (total int, err error) {
	const getTotal = `select count(*) over() as total from avatars where owner_id = $1`

	err = t.tx.GetContext(ctx, &total, getTotal, ownerID)
	if err != nil {
		return 0, fmt.Errorf("db.GetContext: %w", convertErr(err))
	}

	return total, nil
}

// UpdateCurrentAvatar for implements app.Repo.
func (t *txRepo) UpdateCurrentAvatar(ctx context.Context, userID, fileID uuid.UUID) error {
	const query = `update users set current_avatar_id = $1 where id = $2;`

	_, err := t.tx.ExecContext(ctx, query, fileID, userID)
	if err != nil {
		return fmt.Errorf("db.GetContext: %w", convertErr(err))
	}

	return nil
}

// SaveTask implements app.Repo.
func (t *txRepo) SaveTask(ctx context.Context, task app.Task) (id uuid.UUID, err error) {
	newTask, err := convertTask(task)
	if err != nil {
		return uuid.Nil, fmt.Errorf("convertTask: %w", err)
	}

	const query = `
		insert into 
		tasks 
		    (user_bytes, kind) 
		values
			($1, $2)
		returning id
		`

	err = t.tx.GetContext(ctx, &id, query, newTask.UserBytes, newTask.Kind)
	if err != nil {
		return uuid.Nil, fmt.Errorf("t.tx.GetContext: %w", convertErr(err))
	}

	return id, nil
}

// FinishTask implements app.Repo.
func (t *txRepo) FinishTask(ctx context.Context, id uuid.UUID) error {
	const query = `
		update tasks set 
		updated_at = now(),
		finished_at = now()
		where id = $1`

	_, err := t.tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("t.tx.ExecContext: %w", convertErr(err))
	}

	return nil
}

// ListActualTask implements app.Repo.
func (t *txRepo) ListActualTask(ctx context.Context, limit int) ([]app.Task, error) {
	const query = `select * from tasks where finished_at is null order by created_at asc limit $1 for update`

	res := make([]task, 0, limit)
	err := t.tx.SelectContext(ctx, &res, query, limit)
	if err != nil {
		return nil, fmt.Errorf("t.tx.SelectContext: %w", convertErr(err))
	}

	tasks := make([]app.Task, len(res))
	for i := range res {
		t, err := res[i].convert()
		if err != nil {
			return nil, fmt.Errorf("convert: %w", err)
		}

		tasks[i] = *t
	}

	return tasks, nil
}

// SaveStatusUpdateRequest implements app.Repo.
func (t *txRepo) SaveStatusUpdateRequest(ctx context.Context, request app.StatusUpdateRequest) (id uuid.UUID, err error) {
	newRequest := convertStatusUpdateRequest(request)
	const query = `insert into status_update_request (user_id, solution_status) values ($1, $2) returning id`

	err = t.tx.GetContext(ctx, &id, query, newRequest.UserID, newRequest.SolutionStatus)
	if err != nil {
		return uuid.Nil, fmt.Errorf("t.tx.GetContext: %w", convertErr(err))
	}

	return id, nil
}

// SearchStatusUpdateRequest implements app.Repo.
func (t *txRepo) SearchStatusUpdateRequest(ctx context.Context, params app.SearchStatusUpdateRequest,
) (requests []app.StatusUpdateRequest, total int, err error) {
	const query = `select * from status_update_request where solution_status = $1 order by created_at asc limit $2 offset $3`

	res := make([]statusUpdateRequest, 0, params.Limit)
	err = t.tx.SelectContext(ctx, &res, query, params.SolutionStatus.String(), params.Limit, params.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("t.tx.SelectContext: %w", convertErr(err))
	}

	if len(res) == 0 {
		return nil, 0, nil
	}

	const getTotal = `select count(*) over() as total from status_update_request where solution_status = $1`
	err = t.tx.GetContext(ctx, &total, getTotal, params.SolutionStatus.String())
	if err != nil {
		return nil, 0, fmt.Errorf("db.GetContext: %w", err)
	}

	requests = make([]app.StatusUpdateRequest, len(res))
	for i := range res {
		requests[i] = *res[i].convert()
	}

	return requests, total, nil
}

// UpdateStatusUpdateRequest implements app.Repo.
func (t *txRepo) UpdateStatusUpdateRequest(ctx context.Context, request app.StatusUpdateRequest) (*app.StatusUpdateRequest, error) {
	const query = `update status_update_request 
						set 
						    solution_status = $1
						where id = $2
						returning *`

	var req statusUpdateRequest
	err := t.tx.GetContext(ctx, &req, query, request.SolutionStatus.String(), request.ID)
	if err != nil {
		return nil, fmt.Errorf("db.GetContext: %w", convertErr(err))
	}

	return req.convert(), nil
}

// GetStatusUpdateRequestByUserID implements app.Repo.
func (t *txRepo) GetStatusUpdateRequestByUserID(ctx context.Context, userID uuid.UUID) (*app.StatusUpdateRequest, error) {
	const query = `select * from status_update_request where user_id = $1;`

	res := statusUpdateRequest{}
	err := t.tx.GetContext(ctx, &res, query, userID)
	if err != nil {
		return nil, fmt.Errorf("db.GetContext: %w", convertErr(err))
	}

	return res.convert(), nil
}

// UsersByIDs for implements app.Repo.
func (t *txRepo) UsersByIDs(ctx context.Context, ids []uuid.UUID) (users []app.User, err error) {
	const query = `select * from users where id = any($1) for update`

	res := make([]user, 0, len(ids))

	err = t.tx.SelectContext(ctx, &res, query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("t.tx.SelectContext: %w", convertErr(err))
	}

	users = make([]app.User, len(res))
	for i := range res {
		users[i] = *res[i].convert()
	}

	return users, nil
}

// Tx implements app.Repo.
func (*txRepo) Tx(_ context.Context, _ func(app.Repo) error) error {
	panic("you can't start new transaction in current transaction")
}
