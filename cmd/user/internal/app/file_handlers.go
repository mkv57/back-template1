package app

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/ZergsLaw/back-template/internal/dom"
)

const (
	maxAvatarCountInUser = 10
)

// SaveAvatar save info about avatar.
func (a *App) SaveAvatar(ctx context.Context, session dom.Session, file Avatar) (avatarID uuid.UUID, err error) {
	if err = validateFormat(file.ContentType); err != nil {
		return uuid.Nil, fmt.Errorf("validateFormat: %w", err)
	}

	err = a.repo.Tx(ctx, func(repo Repo) error {
		count, err := repo.GetCountAvatars(ctx, session.UserID)
		switch {
		default:
			return fmt.Errorf("a.repo.GetCountAvatars: %w", err)
		case err == nil || errors.Is(err, ErrNotFound):
		}

		if count >= maxAvatarCountInUser {
			return ErrMaxFiles
		}

		avatarID, err = a.file.UploadFile(ctx, file)
		if err != nil {
			return fmt.Errorf("a.file.UploadFile: %w", err)
		}

		fileCache := AvatarInfo{
			FileID:  avatarID,
			OwnerID: session.UserID,
		}
		if err = repo.SaveAvatar(ctx, fileCache); err != nil {
			return fmt.Errorf("repo.SaveAvatar: %w", err)
		}

		user, err := repo.ByID(ctx, session.UserID)
		if err != nil {
			return fmt.Errorf("repo.ByID: %w", err)
		}
		user.AvatarID = avatarID

		_, err = repo.Update(ctx, *user)
		if err != nil {
			return fmt.Errorf("repo.Update: %w", err)
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return avatarID, nil
}

// RemoveAvatar remove info about avatar.
func (a *App) RemoveAvatar(ctx context.Context, session dom.Session, fileID uuid.UUID) error {
	fileCache, err := a.repo.GetAvatar(ctx, fileID)
	if err != nil {
		return fmt.Errorf("a.user.GetAvatarCache: %w", err)
	}

	if fileCache.OwnerID != session.UserID {
		return ErrAccessDenied
	}

	return a.repo.Tx(ctx, func(repo Repo) error {
		if err = repo.DeleteAvatar(ctx, session.UserID, fileID); err != nil {
			return fmt.Errorf("a.user.DeleteAvatarCache: %w", err)
		}

		if err = a.file.DeleteFile(ctx, fileID); err != nil {
			return fmt.Errorf("a.file.RemoveObject: %w", err)
		}

		filesInCache, err := repo.ListAvatarByUserID(ctx, session.UserID)
		if err != nil {
			return fmt.Errorf("repo.ListAvatarByUserID: %w", err)
		}

		newAvatarID := uuid.Nil
		if len(filesInCache) > 0 {
			newAvatarID = filesInCache[0].FileID
		}

		user, err := repo.ByID(ctx, session.UserID)
		if err != nil {
			return fmt.Errorf("repo.ByID: %w", err)
		}
		user.AvatarID = newAvatarID

		_, err = repo.Update(ctx, *user)
		if err != nil {
			return fmt.Errorf("repo.Update: %w", err)
		}

		return nil
	})
}

// ListUserAvatars get list user avatars.
func (a *App) ListUserAvatars(ctx context.Context, session dom.Session) ([]AvatarInfo, error) {
	return a.repo.ListAvatarByUserID(ctx, session.UserID)
}

// GetFile get info about user file by file id.
func (a *App) GetFile(ctx context.Context, _ dom.Session, fileID uuid.UUID) (*Avatar, error) {
	_, err := a.repo.GetAvatar(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("a.user.GetAvatarCache: %w", err)
	}

	file, err := a.file.DownloadFile(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("a.file.GetObject: %w", err)
	}

	return file, nil
}

func validateFormat(contentType string) error {
	const contentTypeSize = 2
	splits := strings.SplitN(contentType, "/", contentTypeSize)

	if len(splits) < contentTypeSize {
		return ErrInvalidImageFormat
	}

	if err := validateFileFormat(strings.ToLower(splits[1])); err != nil {
		return fmt.Errorf("validateFileFormat: %w", err)
	}

	return nil
}
