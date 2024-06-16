package model

import (
	"fmt"
	"time"

	"github.com/ChillyWR/PasswordManager/pkg/pmerror"
	"github.com/ChillyWR/PasswordManager/pkg/pmtime"
	"github.com/google/uuid"
)

type RecordType string

const (
	SecureNoteRecordType RecordType = "secure_note" // default
	LoginRecordType      RecordType = "login"
	CardRecordType       RecordType = "card"
	IdentityRecordType   RecordType = "identity"
)

func NewCredentialRecord(name string, notes *string, userID uuid.UUID) *CredentialRecord {
	return &CredentialRecord{
		ID:        uuid.New(),
		Name:      name,
		Notes:     notes,
		CreatedOn: pmtime.TruncateToMillisecond(time.Now().UTC()),
		UpdatedOn: pmtime.TruncateToMillisecond(time.Now().UTC()),
		CreatedBy: userID,
		UpdatedBy: userID,
	}
}

type CredentialRecord struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Notes     *string   `json:"notes"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	CreatedBy uuid.UUID `json:"created_by"`
}

func (r *CredentialRecord) ApplyForm(f *CredentialRecordForm) {
	if f.Name != nil {
		r.Name = *f.Name
	}

	if f.Notes != nil {
		r.Notes = f.Notes
	}
}

type LoginRecord struct {
	CredentialRecord
	Username *string `json:"username"`
	Password *string `json:"password"`
	URL      *string `json:"url"`
}

func (r *LoginRecord) ApplyForm(f *LoginRecordForm) {
	r.CredentialRecord.ApplyForm(&f.CredentialRecordForm)

	if f.Username != nil {
		r.Username = f.Username
	}

	if f.Password != nil {
		r.Password = f.Password
	}

	if f.URL != nil {
		r.URL = f.URL
	}
}

type CardRecord struct {
	CredentialRecord
	Brand           *string `json:"brand"`
	Number          *string `json:"number"`
	ExpirationMonth *string `json:"expiration_month"`
	ExpirationYear  *string `json:"expiration_year"`
	CVV             *string `json:"cvv"`
}

func (r *CardRecord) ApplyForm(f *CardRecordForm) {
	r.CredentialRecord.ApplyForm(&f.CredentialRecordForm)

	if f.Brand != nil {
		r.Brand = f.Brand
	}

	if f.Number != nil {
		r.Number = f.Number
	}

	if f.ExpirationMonth != nil {
		r.ExpirationMonth = f.ExpirationMonth
	}

	if f.ExpirationYear != nil {
		r.ExpirationYear = f.ExpirationYear
	}

	if f.CVV != nil {
		r.CVV = f.CVV
	}
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

func (r IdentityRecord) ApplyForm(f *IdentityRecordForm) {
	r.CredentialRecord.ApplyForm(&f.CredentialRecordForm)

	if f.FirstName != nil {
		r.FirstName = f.FirstName
	}

	if f.MiddleName != nil {
		r.MiddleName = f.MiddleName
	}

	if f.LastName != nil {
		r.LastName = f.LastName
	}

	if f.Address != nil {
		r.Address = f.Address
	}

	if f.Email != nil {
		r.Email = f.Email
	}

	if f.PhoneNumber != nil {
		r.PhoneNumber = f.PhoneNumber
	}

	if f.PassportNumber != nil {
		r.PassportNumber = f.PassportNumber
	}

	if f.Country != nil {
		r.Country = f.Country
	}
}

// Forms are meant to be filled by user

type CredentialRecordForm struct {
	Name  *string `json:"name"`
	Notes *string `json:"notes"`
}

func (f CredentialRecordForm) Validate() error {
	if f.Name != nil && *f.Name == "" {
		return fmt.Errorf("%w: Name is empty", pmerror.ErrInvalidInput)
	}

	if f.Notes != nil && *f.Notes == "" {
		return fmt.Errorf("%w: Notes is empty", pmerror.ErrInvalidInput)
	}

	return nil
}

func (f CredentialRecordForm) Empty() bool {
	return (f.Name == nil || *f.Name == "") &&
		(f.Notes == nil || *f.Notes == "")
}

type LoginRecordForm struct {
	CredentialRecordForm
	Username *string `json:"username"`
	Password *string `json:"password"`
	URL      *string `json:"url"`
}

func (f LoginRecordForm) Validate() error {
	if err := f.CredentialRecordForm.Validate(); err != nil {
		return fmt.Errorf("validate credential record: %w", err)
	}

	return nil
}

func (f LoginRecordForm) Empty() bool {
	return (f.Username == nil || *f.Username == "") &&
		(f.Password == nil || *f.Password == "") &&
		(f.URL == nil || *f.URL == "")
}

type CardRecordForm struct {
	CredentialRecordForm
	Brand           *string `json:"brand"`
	Number          *string `json:"number"`
	ExpirationMonth *string `json:"expiration_month"`
	ExpirationYear  *string `json:"expiration_year"`
	CVV             *string `json:"cvv"`
}

func (f CardRecordForm) Validate() error {
	if err := f.CredentialRecordForm.Validate(); err != nil {
		return fmt.Errorf("validate credential record: %w", err)
	}

	return nil
}

func (f CardRecordForm) Empty() bool {
	return (f.Brand == nil || *f.Brand == "") &&
		(f.Number == nil || *f.Number == "") &&
		(f.ExpirationMonth == nil || *f.ExpirationMonth == "") &&
		(f.ExpirationYear == nil || *f.ExpirationYear == "") &&
		(f.CVV == nil || *f.CVV == "")
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

func (f IdentityRecordForm) Validate() error {
	if err := f.CredentialRecordForm.Validate(); err != nil {
		return fmt.Errorf("validate credential record: %w", err)
	}

	return nil
}

func (f IdentityRecordForm) Empty() bool {
	return (f.FirstName == nil || *f.FirstName == "") &&
		(f.MiddleName == nil || *f.MiddleName == "") &&
		(f.LastName == nil || *f.LastName == "") &&
		(f.Address == nil || *f.Address == "") &&
		(f.Email == nil || *f.Email == "") &&
		(f.PhoneNumber == nil || *f.PhoneNumber == "") &&
		(f.PassportNumber == nil || *f.PassportNumber == "") &&
		(f.Country == nil || *f.Country == "")
}
