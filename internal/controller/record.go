package controller

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/okutsen/PasswordManager/model"
	"github.com/okutsen/PasswordManager/pkg/pmerror"
)

func (c *Controller) AllRecords() ([]model.CredentialRecord, error) {
	return c.recordRepo.GetAll()
}

func (c *Controller) CredentialRecord(id uuid.UUID, userID uuid.UUID) (*model.CredentialRecord, error) {
	record, err := c.recordRepo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	if record.Notes != nil && record.CreatedBy == userID {
		decNotes, err := Decrypt(*record.Notes, Salt)
		if err != nil {
			return nil, fmt.Errorf("decrypt: %w", err)
		}
	
		record.Notes = &decNotes
	}

	return record, nil
}

func (c *Controller) CreateRecord(form *model.CredentialRecordForm, userID uuid.UUID) (*model.CredentialRecord, error) {
	if err := form.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	encNotes, err := Encrypt(*form.Notes, Salt)
	if err != nil {
		return nil, err
	}

	record := model.CredentialRecord{
		ID:        uuid.New(),
		Name:      *form.Name,
		Notes:     &encNotes,
		CreatedOn: time.Now().UTC(),
		UpdatedOn: time.Now().UTC(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	return c.recordRepo.Create(&record)
}

func (c *Controller) UpdateRecord(id uuid.UUID, form *model.CredentialRecordForm, userID uuid.UUID) (*model.CredentialRecord, error) {
	if form.Empty() {
		return nil, fmt.Errorf("%w: empty form", pmerror.ErrInvalidInput)
	}

	if err := c.authorizeRecord(id, userID); err != nil {
		return nil, fmt.Errorf("authorize: %w", err)
	}

	record := model.CredentialRecord{
		ID:        id,
		UpdatedOn: time.Now().UTC(),
		UpdatedBy: userID,
	}

	if form.Name != nil {
		record.Name = *form.Name
	}

	if form.Notes != nil {
		encNotes, err := Encrypt(*form.Notes, Salt)
		if err != nil {
			return nil, err
		}

		record.Notes = &encNotes
	}

	return c.recordRepo.Update(&record)
}

func (c *Controller) DeleteRecord(id uuid.UUID, userID uuid.UUID) (*model.CredentialRecord, error) {
	if err := c.authorizeRecord(id, userID); err != nil {
		return nil, fmt.Errorf("authorize: %w", err)
	}

	return c.recordRepo.Delete(id)
}

func (c *Controller) authorizeRecord(id uuid.UUID, userID uuid.UUID) error {
	record, err := c.recordRepo.Get(id)
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}

	if record.CreatedBy != userID {
		return fmt.Errorf("%w: user %s does not own record %s", pmerror.ErrForbidden, userID.String(), id.String())
	}

	return nil
}
