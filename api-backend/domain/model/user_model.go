package model

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       UserID
	Email    Email
	Password string
	Unit     string
}

type (
	UserID string
	Email  string
)

var (
	emailValidator  = regexp.MustCompile(`^.+@.+$`)
	digitValidator  = regexp.MustCompile(`\d`)
	letterValidator = regexp.MustCompile(`[a-zA-Z]`)
)

const minimumPasswordLength = 8

func NewUser(id UserID, email Email, pw, unit string) (*User, error) {
	if err := passwordSpecSatisfied(pw); err != nil {
		return nil, errors.Wrapf(err, "invalid password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encrypt password")
	}

	u := &User{
		ID:       id,
		Email:    email,
		Password: string(hash),
		Unit:     unit,
	}

	if err := userSpecSatisfied(*u); err != nil {
		return nil, errors.Wrapf(err, "fail to satisfy User spec")
	}

	return u, nil
}

func passwordSpecSatisfied(pw string) error {
	if !digitValidator.MatchString(pw) || !letterValidator.MatchString(pw) {
		return errors.Errorf("password must contains at least one digit and letter")
	}

	if len(pw) < minimumPasswordLength {
		return errors.Errorf("password must contains at least eight characters")
	}

	return nil
}

func userSpecSatisfied(u User) error {
	if !emailValidator.MatchString(string(u.Email)) {
		return errors.Errorf("invalid email pattern")
	}

	if strings.TrimSpace(u.Unit) == "" {
		return errors.Errorf("invalid unit input")
	}

	return nil
}

func (u *User) ValidatePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return errors.Wrapf(err, "fail to authenticate password")
	}

	return nil
}
