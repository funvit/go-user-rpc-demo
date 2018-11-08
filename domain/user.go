package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// User model
type User struct {
	ID          uuid.UUID
	Login       string
	RegistredAt time.Time
}

// CreateNewUser implements creating new entity of User
func CreateNewUser(login string) *User {
	return &User{
		ID:          uuid.New(),
		Login:       login,
		RegistredAt: time.Now(),
	}
}

// Validate implements validation of User model
func (u *User) Validate() error {
	if u.Login == "" || strings.TrimSpace(u.Login) == "" {
		return errors.New("login property cant be empty")
	}
	return nil
}
