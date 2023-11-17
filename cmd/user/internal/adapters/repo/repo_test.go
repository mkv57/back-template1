//go:build integration

package repo_test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/samber/lo"

	"github.com/ZergsLaw/back-template/cmd/user/internal/app"
	"github.com/ZergsLaw/back-template/internal/dom"
)

type customRepoInterface interface {
	app.Repo
	Delete(ctx context.Context, id uuid.UUID) error
}

func TestRepo_Smoke(t *testing.T) {
	t.Parallel()

	ctx, r, assert := start(t)

	user := app.User{
		ID:        uuid.Nil,
		Email:     "email@gmail.com",
		Name:      "username",
		FullName:  "Elon Musk",
		Status:    dom.UserStatusDefault,
		PassHash:  []byte("pass"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user2 := user
	user3 := user
	user3.Name = "username3"
	user3.Email = "user3@gmail.com"

	id, err := r.Save(ctx, user)
	assert.NoError(err)
	assert.NotNil(id)
	user.ID = id

	user.Name = "new_username"
	user.FullName = "Elon Musk2"
	_, err = r.Update(ctx, user)
	assert.NoError(err)

	_, err = r.Save(ctx, user2)
	assert.ErrorIs(err, app.ErrEmailExist)

	user2.Email = "free@gmail.com"
	user2.Name = user.Name
	_, err = r.Save(ctx, user2)
	assert.ErrorIs(err, app.ErrUsernameExist)

	user3ID, err := r.Save(ctx, user3)
	assert.NoError(err)
	assert.NotNil(user3ID)

	res, err := r.ByID(ctx, user.ID)
	assert.NoError(err)
	user.CreatedAt = res.CreatedAt
	user.UpdatedAt = res.UpdatedAt
	assert.Equal(user, *res)

	res, err = r.ByEmail(ctx, user.Email)
	assert.NoError(err)
	assert.Equal(user, *res)

	res, err = r.ByUsername(ctx, user.Name)
	assert.NoError(err)
	assert.Equal(user, *res)

	users, err := r.UsersByIDs(ctx, []uuid.UUID{user.ID, user3ID})
	assert.NoError(err)
	assert.Len(users, 2)

	users, err = r.UsersByIDs(ctx, []uuid.UUID{uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())})
	assert.NoError(err)
	assert.Len(users, 0)

	listRes, total, err := r.SearchUsers(ctx, app.SearchParams{OwnerID: user3ID, Username: user.Name, FullName: user.FullName, Limit: 5})
	assert.NoError(err)
	assert.Equal(1, total)
	assert.Equal([]app.User{user}, listRes)

	err = r.Delete(ctx, id)
	assert.NoError(err)

	res, err = r.ByID(ctx, user.ID)
	assert.Nil(res)
	assert.ErrorIs(err, app.ErrNotFound)

	taskAdd := app.Task{
		User: user,
		Kind: app.TaskKindEventAdd,
	}
	taskDel := app.Task{
		User: user,
		Kind: app.TaskKindEventDel,
	}
	taskDel.User.PassHash = nil

	taskIDAdd, err := r.SaveTask(ctx, taskAdd)
	assert.NoError(err)
	assert.NotEmpty(taskIDAdd)
	taskAdd.ID = taskIDAdd

	taskIDDel, err := r.SaveTask(ctx, taskDel)
	assert.NoError(err)
	assert.NotEmpty(taskIDDel)
	taskDel.ID = taskIDDel

	tasks, err := r.ListActualTask(ctx, 5)
	assert.NoError(err)
	assert.Len(tasks, 2)

	for i := range tasks {
		assert.NotEmpty(tasks[i].CreatedAt)
		tasks[i].CreatedAt = time.Time{}
		assert.NotEmpty(tasks[i].UpdatedAt)
		tasks[i].UpdatedAt = time.Time{}
		assert.Empty(tasks[i].FinishedAt)
		assert.NotEmpty(tasks[i].User.CreatedAt)
		tasks[i].User.CreatedAt = time.Time{}
		assert.NotEmpty(tasks[i].User.UpdatedAt)
		tasks[i].User.UpdatedAt = time.Time{}
	}
	expListActualTask := lo.Map([]app.Task{taskAdd, taskDel}, func(item app.Task, _ int) app.Task {
		item.User.PassHash = nil
		item.User.CreatedAt = time.Time{}
		item.User.UpdatedAt = time.Time{}
		return item
	})
	sort.Slice(expListActualTask, func(i, j int) bool {
		return expListActualTask[i].ID.String() > expListActualTask[j].ID.String()
	})
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID.String() > tasks[j].ID.String()
	})
	assert.Equal(expListActualTask, tasks)

	err = r.FinishTask(ctx, taskAdd.ID)
	assert.NoError(err)

	tasks, err = r.ListActualTask(ctx, 5)
	assert.NoError(err)
	assert.Len(tasks, 1)

	err = r.FinishTask(ctx, taskDel.ID)
	assert.NoError(err)

	for i := range tasks {
		assert.NotEmpty(tasks[i].CreatedAt)
		tasks[i].CreatedAt = time.Time{}
		assert.NotEmpty(tasks[i].UpdatedAt)
		tasks[i].UpdatedAt = time.Time{}
		assert.Empty(tasks[i].FinishedAt)

		assert.NotEmpty(tasks[i].User.CreatedAt)
		tasks[i].User.CreatedAt = time.Time{}
		assert.NotEmpty(tasks[i].User.UpdatedAt)
		tasks[i].User.UpdatedAt = time.Time{}
	}

	taskDel.User.CreatedAt = time.Time{}
	taskDel.User.UpdatedAt = time.Time{}
	assert.Equal([]app.Task{taskDel}, tasks)

	err = r.Tx(ctx, func(r app.Repo) error {
		return nil
	})
	assert.NoError(err)

	err = r.Tx(ctx, func(r app.Repo) error {
		return app.ErrNotFound
	})

	assert.ErrorIs(err, app.ErrNotFound)
}

func TestRepo_Tx(t *testing.T) {
	t.Parallel()

	ctx, r, assert := start(t)

	err := r.Tx(ctx, func(rep app.Repo) error {
		r := rep.(customRepoInterface)

		user := app.User{
			ID:        uuid.Nil,
			Email:     "email@gmail.com",
			Name:      "username",
			FullName:  "Elon Musk",
			PassHash:  []byte("pass"),
			Status:    dom.UserStatusDefault,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		user2 := user
		user2.Name = "username2"
		user2.FullName = "Elon Sipki"
		user2.Email = "user2@gmail.com"

		id, err := r.Save(ctx, user)
		assert.NoError(err)
		assert.NotNil(id)
		user.ID = id

		user2ID, err := r.Save(ctx, user2)
		assert.NoError(err)
		assert.NotNil(user2ID)

		user.Name = "new_username"
		_, err = r.Update(ctx, user)
		assert.NoError(err)

		res, err := r.ByID(ctx, user.ID)
		assert.NoError(err)
		user.CreatedAt = res.CreatedAt
		user.UpdatedAt = res.UpdatedAt
		assert.Equal(user, *res)

		res, err = r.ByEmail(ctx, user.Email)
		assert.NoError(err)
		assert.Equal(user, *res)

		res, err = r.ByUsername(ctx, user.Name)
		assert.NoError(err)
		assert.Equal(user, *res)

		users, err := r.UsersByIDs(ctx, []uuid.UUID{user.ID, user2ID})
		assert.NoError(err)
		assert.Len(users, 2)

		users, err = r.UsersByIDs(ctx, []uuid.UUID{uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())})
		assert.NoError(err)
		assert.Len(users, 0)

		listRes, total, err := r.SearchUsers(ctx, app.SearchParams{Username: user.Name, FullName: user.FullName, Limit: 5})
		assert.NoError(err)
		assert.Equal(1, total)
		assert.Equal([]app.User{user}, listRes)

		err = r.Delete(ctx, id)
		assert.NoError(err)

		res, err = r.ByID(ctx, user.ID)
		assert.Nil(res)
		assert.ErrorIs(err, app.ErrNotFound)

		taskAdd := app.Task{
			User: user,
			Kind: app.TaskKindEventAdd,
		}
		taskDel := app.Task{
			User: user,
			Kind: app.TaskKindEventDel,
		}
		taskDel.User.PassHash = nil

		taskIDAdd, err := r.SaveTask(ctx, taskAdd)
		assert.NoError(err)
		assert.NotEmpty(taskIDAdd)
		taskAdd.ID = taskIDAdd

		taskIDDel, err := r.SaveTask(ctx, taskDel)
		assert.NoError(err)
		assert.NotEmpty(taskIDDel)
		taskDel.ID = taskIDDel

		tasks, err := r.ListActualTask(ctx, 5)
		assert.NoError(err)
		assert.Len(tasks, 2)

		for i := range tasks {
			assert.NotEmpty(tasks[i].CreatedAt)
			tasks[i].CreatedAt = time.Time{}
			assert.NotEmpty(tasks[i].UpdatedAt)
			tasks[i].UpdatedAt = time.Time{}
			assert.Empty(tasks[i].FinishedAt)
			assert.NotEmpty(tasks[i].User.CreatedAt)
			tasks[i].User.CreatedAt = time.Time{}
			assert.NotEmpty(tasks[i].User.UpdatedAt)
			tasks[i].User.UpdatedAt = time.Time{}
		}
		expListActualTask := lo.Map([]app.Task{taskAdd, taskDel}, func(item app.Task, _ int) app.Task {
			item.User.PassHash = nil
			item.User.CreatedAt = time.Time{}
			item.User.UpdatedAt = time.Time{}
			return item
		})
		sort.Slice(expListActualTask, func(i, j int) bool {
			return expListActualTask[i].ID.String() > expListActualTask[j].ID.String()
		})
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].ID.String() > tasks[j].ID.String()
		})
		assert.Equal(expListActualTask, tasks)

		err = r.FinishTask(ctx, taskAdd.ID)
		assert.NoError(err)

		tasks, err = r.ListActualTask(ctx, 5)
		assert.NoError(err)
		assert.Len(tasks, 1)

		err = r.FinishTask(ctx, taskDel.ID)
		assert.NoError(err)

		for i := range tasks {
			assert.NotEmpty(tasks[i].CreatedAt)
			tasks[i].CreatedAt = time.Time{}
			assert.NotEmpty(tasks[i].UpdatedAt)
			tasks[i].UpdatedAt = time.Time{}
			assert.Empty(tasks[i].FinishedAt)

			assert.NotEmpty(tasks[i].User.CreatedAt)
			tasks[i].User.CreatedAt = time.Time{}
			assert.NotEmpty(tasks[i].User.UpdatedAt)
			tasks[i].User.UpdatedAt = time.Time{}
		}

		taskDel.User.CreatedAt = time.Time{}
		taskDel.User.UpdatedAt = time.Time{}
		assert.Equal([]app.Task{taskDel}, tasks)

		return nil
	})
	assert.NoError(err)

	err = r.Tx(ctx, func(r app.Repo) error {
		return app.ErrNotFound
	})
	assert.ErrorIs(err, app.ErrNotFound)
}
