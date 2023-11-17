// Package app contains business logic.
package app

// App manages business logic methods.
type App struct {
	session Repo
	auth    Auth
	id      ID
	queue   Queue
}

// New build and returns new App.
func New(r Repo, a Auth, id ID, q Queue) *App {
	return &App{
		session: r,
		auth:    a,
		id:      id,
		queue:   q,
	}
}
