package app

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ZergsLaw/back-template1/internal/dom"
)

type (
	// Repo interface for user data repository.
	Repo interface {
		FileInfoRepo
		TaskRepo
		// Tx starts transaction in database.
		// Errors: unknown.
		Tx(ctx context.Context, f func(Repo) error) error
		// Save adds to the new user to repository.
		// Errors: ErrEmailExist, ErrUsernameExist, unknown.
		Save(context.Context, User) (uuid.UUID, error)
		// Update update user info.
		// Errors: ErrUsernameExist, ErrEmailExist, unknown.
		Update(context.Context, User) (*User, error)
		// ByID returns user info by id.
		// Errors: ErrNotFound, unknown.
		ByID(context.Context, uuid.UUID) (*User, error)
		// ByEmail returns user info by email.
		// Errors: ErrNotFound, unknown.
		ByEmail(context.Context, string) (*User, error)
		// ByUsername returns user info by username.
		// Errors: ErrNotFound, unknown.
		ByUsername(context.Context, string) (*User, error)
		// SearchUsers returns list user info.
		// Errors: unknown.
		SearchUsers(context.Context, SearchParams) ([]User, int, error)
		// UsersByIDs returns list of users.
		// Errors: ErrNotFound, unknown.
		UsersByIDs(ctx context.Context, ids []uuid.UUID) (users []User, err error)
	}
	// FileInfoRepo provides to file info repository
	FileInfoRepo interface {
		// SaveAvatar adds to the new cache about user avatar to repository.
		// Errors: ErrUserIDAndFileIDExist, ErrMaximumNumberOfStoredFilesReached, unknown.
		SaveAvatar(ctx context.Context, fileCache AvatarInfo) error
		// DeleteAvatar delete cache about user avatar info.
		// Errors: unknown.
		DeleteAvatar(ctx context.Context, userID, fileID uuid.UUID) error
		// GetAvatar returns cache about user avatar by id.
		// Errors: ErrNotFound, unknown.
		GetAvatar(ctx context.Context, fileID uuid.UUID) (*AvatarInfo, error)
		// ListAvatarByUserID returns list cache user file.
		// Errors: unknown.
		ListAvatarByUserID(ctx context.Context, userID uuid.UUID) ([]AvatarInfo, error)
		// GetCountAvatars returns count user avatars.
		// Errors: ErrNotFound, unknown.
		GetCountAvatars(ctx context.Context, ownerID uuid.UUID) (total int, err error)
	}

	// FileStore interface for saving and getting files.
	FileStore interface {
		// UploadFile save new file in database.
		// Errors: unknown.
		UploadFile(ctx context.Context, f Avatar) (uuid.UUID, error)
		// DownloadFile get file by id.
		// Errors: unknown.
		DownloadFile(ctx context.Context, id uuid.UUID) (*Avatar, error)
		// DeleteFile delete file by id.
		// Errors: unknown.
		DeleteFile(ctx context.Context, id uuid.UUID) error
	}

	// PasswordHash module responsible for hashing password.
	PasswordHash interface {
		// Hashing returns the hashed version of the password.
		// Errors: unknown.
		Hashing(password string) ([]byte, error)
		// Compare compares two passwords for matches.
		Compare(hashedPassword []byte, password []byte) bool
	}

	// TaskRepo interface for saving tasks.
	TaskRepo interface {
		// SaveTask adds new task to repository.
		// Errors: unknown.
		SaveTask(context.Context, Task) (uuid.UUID, error)
		// FinishTask set column Task.FinishedAt task.
		// Errors: unknown.
		FinishTask(context.Context, uuid.UUID) error
		// ListActualTask returns list task by limit and ordered by created_at (ask).
		// Return tasks without Task.FinishedAt.
		// Errors: unknown.
		ListActualTask(context.Context, int) ([]Task, error)
	}

	// Sessions module for manager user's session.
	Sessions interface {
		// Get returns user session by his token.
		// Errors: ErrNotFound, unknown.
		Get(context.Context, string) (*dom.Session, error)
		// Save new session for specific user.
		// Errors: unknown.
		Save(context.Context, uuid.UUID, dom.Origin, dom.UserStatus) (*dom.Token, error)
		// Delete removes session by id.
		// Errors: ErrNotFound, unknown.
		Delete(context.Context, uuid.UUID) error
	}

	// Queue sends events to queue.
	Queue interface {
		// AddUser sends event 'EventAdd' to queue.
		// Errors: unknown.
		AddUser(context.Context, uuid.UUID, User) error
		// DeleteUser sends event 'EventDel' to queue.
		// Errors: unknown.
		DeleteUser(context.Context, uuid.UUID, User) error
		// UpdateUser sends event 'EventUpdate' to queue.
		// Errors: unknown.
		UpdateUser(context.Context, uuid.UUID, User) error
	}
)
