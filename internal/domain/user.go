package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type UserID string

func NewUserID() UserID {
	return UserID(uuid.New().String())
}

func (id UserID) String() string {
	return string(id)
}

type User struct {
	ID       UserID `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

var (
	ErrInvalidName     = errors.New("name cannot be empty")
	ErrInvalidEmail    = errors.New("email format is invalid")
	ErrInvalidUsername = errors.New("username must be at least 3 characters")
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailExists     = errors.New("email already exists")
	ErrUsernameExists  = errors.New("username already exists")
)

func NewUser(name, email, username string) (*User, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if err := validateUsername(username); err != nil {
		return nil, err
	}

	return &User{
		ID:       NewUserID(),
		Name:     strings.TrimSpace(name),
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Username: strings.ToLower(strings.TrimSpace(username)),
	}, nil
}

func (u *User) UpdateName(name string) error {
	if err := validateName(name); err != nil {
		return err
	}
	u.Name = strings.TrimSpace(name)
	return nil
}

func (u *User) UpdateEmail(email string) error {
	if err := validateEmail(email); err != nil {
		return err
	}
	u.Email = strings.ToLower(strings.TrimSpace(email))
	return nil
}

func (u *User) UpdateUsername(username string) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	u.Username = strings.ToLower(strings.TrimSpace(username))
	return nil
}

func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrInvalidName
	}
	return nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return ErrInvalidEmail
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

func validateUsername(username string) error {
	username = strings.TrimSpace(username)
	if len(username) < 3 {
		return ErrInvalidUsername
	}
	return nil
}

func (u *User) String() string {
	return fmt.Sprintf("User{ID: %s, Name: %s, Email: %s, Username: %s}",
		u.ID, u.Name, u.Email, u.Username)
}
