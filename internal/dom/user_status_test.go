package dom_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ZergsLaw/back-template/internal/dom"
)

func TestUserStatus_IsFreeze(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		want   bool
		status dom.UserStatus
	}{
		"freeze":  {true, dom.UserStatusFreeze},
		"default": {false, dom.UserStatusDefault},
		"premium": {false, dom.UserStatusPremium},
		"support": {false, dom.UserStatusSupport},
		"admin":   {false, dom.UserStatusAdmin},
		"jedi":    {false, dom.UserStatusJedi},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert := require.New(t)
			assert.Equal(tc.want, tc.status.IsFreeze())
		})
	}
}

func TestUserStatus_IsDefault(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		want   bool
		status dom.UserStatus
	}{
		"freeze":  {false, dom.UserStatusFreeze},
		"default": {true, dom.UserStatusDefault},
		"premium": {false, dom.UserStatusPremium},
		"support": {false, dom.UserStatusSupport},
		"admin":   {false, dom.UserStatusAdmin},
		"jedi":    {false, dom.UserStatusJedi},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert := require.New(t)
			assert.Equal(tc.want, tc.status.IsDefault())
		})
	}
}

func TestUserStatus_IsPremium(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		want   bool
		status dom.UserStatus
	}{
		"freeze":  {false, dom.UserStatusFreeze},
		"default": {false, dom.UserStatusDefault},
		"premium": {true, dom.UserStatusPremium},
		"support": {false, dom.UserStatusSupport},
		"admin":   {false, dom.UserStatusAdmin},
		"jedi":    {false, dom.UserStatusJedi},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert := require.New(t)
			assert.Equal(tc.want, tc.status.IsPremium())
		})
	}
}

func TestUserStatus_IsSupport(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		want   bool
		status dom.UserStatus
	}{
		"freeze":  {false, dom.UserStatusFreeze},
		"default": {false, dom.UserStatusDefault},
		"premium": {false, dom.UserStatusPremium},
		"support": {true, dom.UserStatusSupport},
		"admin":   {false, dom.UserStatusAdmin},
		"jedi":    {false, dom.UserStatusJedi},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert := require.New(t)
			assert.Equal(tc.want, tc.status.IsSupport())
		})
	}
}

func TestUserStatus_IsAdmin(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		want   bool
		status dom.UserStatus
	}{
		"freeze":  {false, dom.UserStatusFreeze},
		"default": {false, dom.UserStatusDefault},
		"premium": {false, dom.UserStatusPremium},
		"support": {false, dom.UserStatusSupport},
		"admin":   {true, dom.UserStatusAdmin},
		"jedi":    {false, dom.UserStatusJedi},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert := require.New(t)
			assert.Equal(tc.want, tc.status.IsAdmin())
		})
	}
}

func TestUserStatus_IsJedi(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		want   bool
		status dom.UserStatus
	}{
		"freeze":  {false, dom.UserStatusFreeze},
		"default": {false, dom.UserStatusDefault},
		"premium": {false, dom.UserStatusPremium},
		"support": {false, dom.UserStatusSupport},
		"admin":   {false, dom.UserStatusAdmin},
		"jedi":    {true, dom.UserStatusJedi},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert := require.New(t)
			assert.Equal(tc.want, tc.status.IsJedi())
		})
	}
}

func TestUserStatus_IsSpecialist(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		want   bool
		status dom.UserStatus
	}{
		"freeze":  {false, dom.UserStatusFreeze},
		"default": {false, dom.UserStatusDefault},
		"premium": {false, dom.UserStatusPremium},
		"support": {true, dom.UserStatusSupport},
		"admin":   {true, dom.UserStatusAdmin},
		"jedi":    {true, dom.UserStatusJedi},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert := require.New(t)
			assert.Equal(tc.want, tc.status.IsSpecialist())
		})
	}
}

func TestUserStatus_IsManager(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		want   bool
		status dom.UserStatus
	}{
		"freeze":  {false, dom.UserStatusFreeze},
		"default": {false, dom.UserStatusDefault},
		"premium": {false, dom.UserStatusPremium},
		"support": {false, dom.UserStatusSupport},
		"admin":   {true, dom.UserStatusAdmin},
		"jedi":    {true, dom.UserStatusJedi},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert := require.New(t)
			assert.Equal(tc.want, tc.status.IsManager())
		})
	}
}
