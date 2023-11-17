// Package password will hash and comparable hash-pass.
package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type (
	// Manager contains method for hashing and comparable value.
	Manager struct {
		cost int
	}
	// Option for building Password struct.
	Option func(*Manager)
)

// Cost option for sets hashing cost.
func Cost(cost int) Option {
	return func(password *Manager) {
		password.cost = cost
	}
}

// New creates and returns new Hasher.
func New(options ...Option) *Manager {
	h := &Manager{cost: bcrypt.DefaultCost}

	for i := range options {
		options[i](h)
	}

	return h
}

// Hashing value and returns bytes.
func (m *Manager) Hashing(val string) ([]byte, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(val), m.cost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	return res, nil
}

// Compare comparable two hash.
func (*Manager) Compare(val1 []byte, val2 []byte) bool {
	return bcrypt.CompareHashAndPassword(val1, val2) == nil
}
