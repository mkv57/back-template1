// Package app contains business logic.
package app

// App manages business logic methods.
type App struct {
	repo     Repo
	hash     PasswordHash
	sessions Sessions
	file     FileStore
	queue    Queue
}

// New build and returns new App.
func New(r Repo, ph PasswordHash, a Sessions, f FileStore, q Queue) *App {
	return &App{
		repo:     r,
		hash:     ph,
		sessions: a,
		file:     f,
		queue:    q,
	}
}
