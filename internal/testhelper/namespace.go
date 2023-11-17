package testhelper

import (
	"strings"
	"testing"
)

func Namespace(t *testing.T) string {
	t.Helper()

	namespace := t.Name()
	namespace = strings.ReplaceAll(namespace, "/", "_")
	namespace = strings.ReplaceAll(namespace, "#", "_")
	namespace = strings.ReplaceAll(namespace, ".", "_")

	return namespace
}
