package app

import (
	"net"
	"time"

	"github.com/gofrs/uuid"

	"github.com/ZergsLaw/back-template1/internal/dom"
)

type (
	// Token contains auth token.
	Token struct {
		// Generate by Auth contract.
		Value string
	}

	// User contains user information.
	User struct {
		ID    uuid.UUID // Generate by repository layer.
		Email string
		Name  string
	}

	// Origin information about req user.
	Origin struct {
		IP        net.IP
		UserAgent string
	}

	// Session contains session info for identify a user.
	Session struct {
		ID        uuid.UUID // Generate by repository layer.
		Origin    Origin
		Token     Token
		UserID    uuid.UUID
		Status    dom.UserStatus
		CreatedAt time.Time // Generate by repository layer.
		UpdatedAt time.Time // Generate by repository layer.
	}

	// EventUpdateStatus contains information about change user session status.
	EventUpdateStatus struct {
		Status UpdateStatus
	}

	// UpdateStatus information about updated status.
	UpdateStatus struct {
		UserID uuid.UUID
		Status dom.UserStatus
	}
)
