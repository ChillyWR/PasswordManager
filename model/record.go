package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/okutsen/PasswordManager/pkg/pmerror"
)

const (
	SecureNote = iota
	LoginRecordType
	CardRecordType
	IdentityRecordType
)

type CredentialRecord struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Notes     *string   `json:"notes"`
	CreatedOn time.Time `json:"created_at"`
	UpdatedOn time.Time `json:"updated_at"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	CreatedBy uuid.UUID `json:"created_by"`
}

type LoginRecord struct {
	CredentialRecord
	Username *string `json:"username"`
	Password *string `json:"password"`
	URL      *string `json:"url"`
}

type CardRecord struct {
	CredentialRecord
	Brand           *string `json:"brand"`
	Number          *string `json:"number"`
	ExpirationMonth *string `json:"expiration_month"`
	ExpirationYear  *string `json:"expiration_year"`
	CVV             *string `json:"cvv"`
}

type IdentityRecord struct {
	CredentialRecord
	FirstName      *string `json:"first_name"`
	MiddleName     *string `json:"middle_name"`
	LastName       *string `json:"last_name"`
	Address        *string `json:"address"`
	Email          *string `json:"email"`
	PhoneNumber    *string `json:"phone_number"`
	PassportNumber *string `json:"passport_number"`
	Country        *string `json:"country"`
}

// Forms are meant to be filled by user

type CredentialRecordForm struct {
	Name  *string `json:"name"`
	Notes *string `json:"notes"`
}

func (f CredentialRecordForm) Validate() error {
	if f.Name == nil || *f.Name == "" {
		return fmt.Errorf("%w: Name is empty", pmerror.ErrInvalidInput)
	}

	return nil
}

func (f CredentialRecordForm) Empty() bool {
	return f.Name == nil && f.Notes == nil
}

type LoginRecordForm struct {
	CredentialRecordForm
	Username *string `json:"username"`
	Password *string `json:"password"`
	URL      *string `json:"url"`
}

type CardRecordForm struct {
	CredentialRecordForm
	Brand           *string `json:"brand"`
	Number          *string `json:"number"`
	ExpirationMonth *string `json:"expiration_month"`
	ExpirationYear  *string `json:"expiration_year"`
	CVV             *string `json:"cvv"`
}

type IdentityRecordForm struct {
	CredentialRecordForm
	FirstName      *string `json:"first_name"`
	MiddleName     *string `json:"middle_name"`
	LastName       *string `json:"last_name"`
	Address        *string `json:"address"`
	Email          *string `json:"email"`
	PhoneNumber    *string `json:"phone_number"`
	PassportNumber *string `json:"passport_number"`
	Country        *string `json:"country"`
}
