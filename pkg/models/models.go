package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching records found")
	ErrInvalidCredentials = errors.New("models: invalid Credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate Email")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}
