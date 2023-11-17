package testhelper

// Mod applies given modification to any value.
func Mod[T any](v T, f func(*T)) T {
	f(&v)

	return v
}
