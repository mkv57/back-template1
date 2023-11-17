package dom

import (
	"net"

	"github.com/gofrs/uuid"
)

type (
	// Session contains main session info.
	Session struct {
		ID     uuid.UUID
		UserID uuid.UUID
		Status UserStatus
	}
	// Origin information about req user.
	Origin struct {
		IP        net.IP
		UserAgent string
	}
	// Token contains user's authorization token.
	Token struct {
		Value string
	}
)
