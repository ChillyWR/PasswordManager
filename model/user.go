package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/okutsen/PasswordManager/pkg/pmerror"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
}

func (User) TableName() string {
	return "reg_user"
  }

// Forms are meant to be filled by user

type UserForm struct {
	Name     *string `json:"name"`
	Password *string `json:"password"`
}

func (f UserForm) Validate() error {
	if f.Name == nil || *f.Name == "" {
		return fmt.Errorf("%w: Password is empty", pmerror.ErrInvalidInput)
	}

	if f.Password == nil || *f.Password == "" {
		return fmt.Errorf("%w: Password is empty", pmerror.ErrInvalidInput)
	}

	return nil
}

func (f UserForm) Empty() bool {
	return (f.Name == nil || *f.Name == "") && (f.Password == nil || *f.Password == "")
}
