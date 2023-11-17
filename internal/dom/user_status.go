package dom

import (
	"fmt"

	user_status_pb "github.com/ZergsLaw/back-template/api/user_status/v1"
)

// UserStatus user status in service.
type UserStatus uint8

//go:generate stringer -output=stringer.UserStatus.go -type=UserStatus -trimprefix=UserStatus
const (
	_ UserStatus = iota
	UserStatusFreeze
	UserStatusDefault
	UserStatusPremium
	UserStatusSupport
	UserStatusAdmin
	UserStatusJedi
)

func (i UserStatus) IsFreeze() bool  { return i == UserStatusFreeze }
func (i UserStatus) IsDefault() bool { return i == UserStatusDefault }
func (i UserStatus) IsPremium() bool { return i == UserStatusPremium }
func (i UserStatus) IsSupport() bool { return i == UserStatusSupport }
func (i UserStatus) IsAdmin() bool   { return i == UserStatusAdmin }
func (i UserStatus) IsJedi() bool    { return i == UserStatusJedi }

func (i UserStatus) IsSpecialist() bool {
	return i == UserStatusSupport || i == UserStatusAdmin || i == UserStatusJedi
}

func (i UserStatus) IsManager() bool { return i == UserStatusAdmin || i == UserStatusJedi }

// UserStatusFromAPI converts user status from protobuf.
func UserStatusFromAPI(kind user_status_pb.StatusKind) UserStatus {
	switch kind { //nolint:exhaustive // By design.
	case user_status_pb.StatusKind_STATUS_KIND_FREEZE:
		return UserStatusFreeze
	case user_status_pb.StatusKind_STATUS_KIND_DEFAULT:
		return UserStatusDefault
	case user_status_pb.StatusKind_STATUS_KIND_PREMIUM:
		return UserStatusPremium
	case user_status_pb.StatusKind_STATUS_KIND_SUPPORT:
		return UserStatusSupport
	case user_status_pb.StatusKind_STATUS_KIND_ADMIN:
		return UserStatusAdmin
	case user_status_pb.StatusKind_STATUS_KIND_JEDI:
		return UserStatusJedi
	default:
		panic(fmt.Sprintf("unknown kind: %s", kind))
	}
}

// UserStatusToAPI converts user status from protobuf.
func UserStatusToAPI(status UserStatus) user_status_pb.StatusKind {
	switch status {
	case UserStatusFreeze:
		return user_status_pb.StatusKind_STATUS_KIND_FREEZE
	case UserStatusDefault:
		return user_status_pb.StatusKind_STATUS_KIND_DEFAULT
	case UserStatusPremium:
		return user_status_pb.StatusKind_STATUS_KIND_PREMIUM
	case UserStatusSupport:
		return user_status_pb.StatusKind_STATUS_KIND_SUPPORT
	case UserStatusAdmin:
		return user_status_pb.StatusKind_STATUS_KIND_ADMIN
	case UserStatusJedi:
		return user_status_pb.StatusKind_STATUS_KIND_JEDI
	default:
		panic(fmt.Sprintf("unknown status: %s", status))
	}
}
