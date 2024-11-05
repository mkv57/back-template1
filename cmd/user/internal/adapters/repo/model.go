package repo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/ZergsLaw/back-template1/cmd/user/internal/app"
	"github.com/ZergsLaw/back-template1/internal/dom"
)

type (
	user struct {
		ID              uuid.UUID     `db:"id" json:"id"`
		Email           string        `db:"email" json:"email"`
		Name            string        `db:"name" json:"name"`
		FullName        string        `db:"full_name" json:"full_name"`
		CurrentAvatarID uuid.NullUUID `db:"current_avatar_id" json:"current_avatar_id,omitempty"`
		PassHash        []byte        `db:"pass_hash" json:"-"`
		Status          string        `db:"status" json:"status"`
		CreatedAt       time.Time     `db:"created_at" json:"created_at"`
		UpdatedAt       time.Time     `db:"updated_at" json:"updated_at"`
	}

	avatar struct {
		ID        uuid.UUID `db:"id"`
		OwnerID   uuid.UUID `db:"owner_id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	task struct {
		ID         uuid.UUID       `db:"id"`
		UserBytes  json.RawMessage `db:"user_bytes"`
		Kind       string          `db:"kind"`
		CreatedAt  time.Time       `db:"created_at"`
		UpdatedAt  time.Time       `db:"updated_at"`
		FinishedAt sql.NullTime    `db:"finished_at"`
	}

	statusUpdateRequest struct {
		ID             uuid.UUID `db:"id"`
		UserID         uuid.UUID `db:"user_id"`
		SolutionStatus string    `db:"solution_status"`
		CreatedAt      time.Time `db:"created_at"`
		UpdatedAt      time.Time `db:"updated_at"`
	}
)

func convert(u app.User) *user {
	return &user{
		ID:       u.ID,
		Email:    u.Email,
		Name:     u.Name,
		FullName: u.FullName,
		CurrentAvatarID: uuid.NullUUID{
			UUID:  u.AvatarID,
			Valid: u.AvatarID != uuid.Nil,
		},
		PassHash:  u.PassHash,
		Status:    u.Status.String(),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u user) convert() *app.User {
	return &app.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		FullName:  u.FullName,
		AvatarID:  u.CurrentAvatarID.UUID,
		PassHash:  u.PassHash,
		Status:    appUserStatus(u.Status),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func convertUserFile(f app.AvatarInfo) *avatar {
	return &avatar{
		ID:        f.FileID,
		OwnerID:   f.OwnerID,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func (f avatar) convert() *app.AvatarInfo {
	return &app.AvatarInfo{
		OwnerID:   f.OwnerID,
		FileID:    f.ID,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func convertTask(s app.Task) (*task, error) {
	userBytes, err := json.Marshal(convert(s.User))
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return &task{
		ID:        s.ID,
		UserBytes: userBytes,
		Kind:      s.Kind.String(),
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		FinishedAt: sql.NullTime{
			Time:  s.FinishedAt,
			Valid: !s.FinishedAt.IsZero(),
		},
	}, nil
}

func (t task) convert() (*app.Task, error) {
	var u user
	err := json.Unmarshal(t.UserBytes, &u)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return &app.Task{
		ID:         t.ID,
		User:       *u.convert(),
		Kind:       appTaskKind(t.Kind),
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
		FinishedAt: t.FinishedAt.Time,
	}, nil
}

func appTaskKind(txt string) app.TaskKind {
	switch txt {
	case app.TaskKindEventAdd.String():
		return app.TaskKindEventAdd
	case app.TaskKindEventDel.String():
		return app.TaskKindEventDel
	case app.TaskKindEventUpdate.String():
		return app.TaskKindEventUpdate
	default:
		panic(fmt.Sprintf("unknown txt: %s", txt))
	}
}

func convertStatusUpdateRequest(r app.StatusUpdateRequest) *statusUpdateRequest {
	return &statusUpdateRequest{
		ID:             r.ID,
		UserID:         r.UserID,
		SolutionStatus: r.SolutionStatus.String(),
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

func (r *statusUpdateRequest) convert() *app.StatusUpdateRequest {
	return &app.StatusUpdateRequest{
		ID:             r.ID,
		UserID:         r.UserID,
		SolutionStatus: appSolutionUpdate(r.SolutionStatus),
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

func appUserStatus(status string) dom.UserStatus {
	switch status {
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
		panic(fmt.Sprintf("unknown status: %s", status))
	}
}

func appSolutionUpdate(solution string) app.SolutionStatus {
	switch solution {
	case app.SolutionStatusNew.String():
		return app.SolutionStatusNew
	case app.SolutionStatusApprove.String():
		return app.SolutionStatusApprove
	case app.SolutionStatusCancel.String():
		return app.SolutionStatusCancel
	default:
		panic(fmt.Sprintf("unknown status: %s", solution))
	}
}
