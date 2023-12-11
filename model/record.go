package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/okutsen/PasswordManager/pkg/pmerror"
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
		CreatedOn: time.Now().UTC(),
		UpdatedOn: time.Now().UTC(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}
}

type CredentialRecord struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Notes     *string   `json:"notes"`
	CreatedOn time.Time `json:"created_at"`
	UpdatedOn time.Time `json:"updated_at"`
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

func (f LoginRecordForm) Validate() error {
	if err := f.CredentialRecordForm.Validate(); err != nil {
		return fmt.Errorf("validate credential record: %w", err)
	}

	return nil
}

func (f LoginRecordForm) Empty() bool {
	return f.Username == nil && f.Password == nil && f.URL == nil
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
	return f.Brand == nil && f.Number == nil && f.ExpirationMonth == nil && f.ExpirationYear == nil && f.CVV == nil
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
	return f.FirstName == nil && f.MiddleName == nil && f.LastName == nil && f.Address == nil && f.Email == nil && f.PhoneNumber == nil && f.Country == nil
}
