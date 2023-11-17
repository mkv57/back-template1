package testhelper

import (
	"context"
	"testing"
	"time"
)

// MaxTimeout for tests.
const MaxTimeout = time.Second * 120

// Context build new context for tests and set cleanup cancel function.
func Context(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), MaxTimeout)
	t.Cleanup(cancel)

	return ctx
}
