package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ChillyWR/PasswordManager/model"
	"github.com/ChillyWR/PasswordManager/pkg/pmerror"
	"github.com/ChillyWR/PasswordManager/pkg/pmtime"
)

func (c *Controller) AllRecords(userID uuid.UUID) ([]model.CredentialRecord, []model.LoginRecord, []model.CardRecord, []model.IdentityRecord, error) {
	return c.recordRepo.GetAll(userID)
}

func (c *Controller) GetRecord(id uuid.UUID, userID uuid.UUID) (interface{}, error) {
	if err := c.authorizeRecord(id, userID); err != nil {
		return nil, fmt.Errorf("authorize: %w", err)
	}

	credentialRecord, err := c.recordRepo.GetCredentialRecord(id)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	login, err := c.recordRepo.GetLogin(id)
	if err == nil {
		if err := c.decryptLogin(login); err != nil {
			return nil, fmt.Errorf("decrypt login: %w", err)
		}

		return login, nil
	} else if !errors.Is(err, pmerror.ErrNotFound) {
		return nil, fmt.Errorf("get login: %w", err)
	}

	card, err := c.recordRepo.GetCard(id)
	if err == nil {
		if err := c.decryptCard(card); err != nil {
			return nil, fmt.Errorf("decrypt card: %w", err)
		}

		return card, nil
	} else if !errors.Is(err, pmerror.ErrNotFound) {
		return nil, fmt.Errorf("get card: %w", err)
	}

	identity, err := c.recordRepo.GetIdentity(id)
	if err == nil {
		if err := c.decryptIdentity(identity); err != nil {
			return nil, fmt.Errorf("decrypt identity: %w", err)
		}

		return identity, nil
	} else if !errors.Is(err, pmerror.ErrNotFound) {
		return nil, fmt.Errorf("get identity: %w", err)
	}

	if err := c.decryptCredentialRecord(credentialRecord); err != nil {
		return nil, fmt.Errorf("decrypt identity: %w", err)
	}

	return credentialRecord, nil
}

func (c *Controller) CreateRecord(recordType model.RecordType, rawForm json.RawMessage, userID uuid.UUID) (interface{}, error) {
	switch recordType {
	case model.SecureNoteRecordType:
		var form model.CredentialRecordForm
		if err := json.Unmarshal(rawForm, &form); err != nil {
			return nil, fmt.Errorf("%w: unmarshal: %s", pmerror.ErrInvalidInput, err.Error())
		}

		if err := form.Validate(); err != nil {
			return nil, fmt.Errorf("validate: %w", err)
		}

		record := model.NewCredentialRecord(*form.Name, form.Notes, userID)

		if err := c.encryptCredentialRecord(record); err != nil {
			return nil, fmt.Errorf("validate: %w", err)
		}

		return c.recordRepo.CreateCredentialRecord(record)
	case model.LoginRecordType:
		var form model.LoginRecordForm
		if err := json.Unmarshal(rawForm, &form); err != nil {
			return nil, fmt.Errorf("%w: unmarshal: %s", pmerror.ErrInvalidInput, err.Error())
		}

		if err := form.Validate(); err != nil {
			return nil, fmt.Errorf("validate: %w", err)
		}

		record := model.LoginRecord{
			CredentialRecord: *model.NewCredentialRecord(*form.Name, form.Notes, userID),
			Username:         form.Username,
			Password:         form.Password,
			URL:              form.URL,
		}

		if err := c.encryptLogin(&record); err != nil {
			return nil, fmt.Errorf("validate: %w", err)
		}

		return c.recordRepo.CreateLogin(&record)
	case model.CardRecordType:
		var form model.CardRecordForm
		if err := json.Unmarshal(rawForm, &form); err != nil {
			return nil, fmt.Errorf("%w: unmarshal: %s", pmerror.ErrInvalidInput, err.Error())
		}

		if err := form.Validate(); err != nil {
			return nil, fmt.Errorf("validate: %w", err)
		}

		record := model.CardRecord{
			CredentialRecord: *model.NewCredentialRecord(*form.Name, form.Notes, userID),
			Brand:            form.Brand,
			Number:           form.Number,
			ExpirationMonth:  form.ExpirationMonth,
			ExpirationYear:   form.ExpirationYear,
			CVV:              form.CVV,
		}

		if err := c.encryptCard(&record); err != nil {
			return nil, fmt.Errorf("validate: %w", err)
		}

		return c.recordRepo.CreateCard(&record)
	case model.IdentityRecordType:
		var form model.IdentityRecordForm
		if err := json.Unmarshal(rawForm, &form); err != nil {
			return nil, fmt.Errorf("%w: unmarshal: %s", pmerror.ErrInvalidInput, err.Error())
		}

		if err := form.Validate(); err != nil {
			return nil, fmt.Errorf("validate: %w", err)
		}

		record := model.IdentityRecord{
			CredentialRecord: *model.NewCredentialRecord(*form.Name, form.Notes, userID),
			FirstName:        form.FirstName,
			MiddleName:       form.MiddleName,
			LastName:         form.LastName,
			Address:          form.Address,
			Email:            form.Email,
			PhoneNumber:      form.PhoneNumber,
			PassportNumber:   form.PassportNumber,
			Country:          form.Country,
		}

		if err := c.encryptIdentity(&record); err != nil {
			return nil, fmt.Errorf("validate: %w", err)
		}

		return c.recordRepo.CreateIdentity(&record)
	default:
		return nil, fmt.Errorf("%w: unsupported record type", pmerror.ErrInvalidInput)
	}
}

func (c *Controller) UpdateRecord(id uuid.UUID, rawForm json.RawMessage, userID uuid.UUID) (interface{}, error) {
	if err := c.authorizeRecord(id, userID); err != nil {
		return nil, fmt.Errorf("authorize: %w", err)
	}

	record, err := c.GetRecord(id, userID)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	switch record.(type) {
	case *model.CredentialRecord:
		var form model.CredentialRecordForm
		if err := json.Unmarshal(rawForm, &form); err != nil {
			return nil, fmt.Errorf("%w: unmarshal: %s", pmerror.ErrInvalidInput, err.Error())
		}

		if form.Empty() {
			return nil, fmt.Errorf("%w: empty form", pmerror.ErrInvalidInput)
		}

		record := model.CredentialRecord{
			ID:        id,
			UpdatedOn: pmtime.TruncateToMillisecond(time.Now().UTC()),
			UpdatedBy: userID,
		}
		record.ApplyForm(&form)

		if err := c.encryptCredentialRecord(&record); err != nil {
			return nil, fmt.Errorf("validate: %w", err)
		}

		return c.recordRepo.UpdateCredentialRecord(&record)
	case *model.LoginRecord:
		var form model.LoginRecordForm
		if err := json.Unmarshal(rawForm, &form); err != nil {
			return nil, fmt.Errorf("%w: unmarshal: %s", pmerror.ErrInvalidInput, err.Error())
		}

		if form.Empty() {
			return nil, fmt.Errorf("%w: empty form", pmerror.ErrInvalidInput)
		}

		record := model.LoginRecord{
			CredentialRecord: model.CredentialRecord{
				ID:        id,
				UpdatedOn: pmtime.TruncateToMillisecond(time.Now().UTC()),
				UpdatedBy: userID,
			},
		}
		record.ApplyForm(&form)

		if err := c.encryptLogin(&record); err != nil {
			return nil, fmt.Errorf("validate: %w", err)
		}

		return c.recordRepo.UpdateLogin(&record)
	case *model.CardRecord:
		var form model.CardRecordForm
		if err := json.Unmarshal(rawForm, &form); err != nil {
			return nil, fmt.Errorf("%w: unmarshal: %s", pmerror.ErrInvalidInput, err.Error())
		}

		if form.Empty() {
			return nil, fmt.Errorf("%w: empty form", pmerror.ErrInvalidInput)
		}

		record := model.CardRecord{
			CredentialRecord: model.CredentialRecord{
				ID:        id,
				UpdatedOn: pmtime.TruncateToMillisecond(time.Now().UTC()),
				UpdatedBy: userID,
			},
		}
		record.ApplyForm(&form)

		if err := c.encryptCard(&record); err != nil {
			return nil, fmt.Errorf("validate: %w", err)
		}

		return c.recordRepo.UpdateCard(&record)
	case *model.IdentityRecord:
		var form model.IdentityRecordForm
		if err := json.Unmarshal(rawForm, &form); err != nil {
			return nil, fmt.Errorf("%w: unmarshal: %s", pmerror.ErrInvalidInput, err.Error())
		}

		if form.Empty() {
			return nil, fmt.Errorf("%w: empty form", pmerror.ErrInvalidInput)
		}

		record := model.IdentityRecord{
			CredentialRecord: model.CredentialRecord{
				ID:        id,
				UpdatedOn: pmtime.TruncateToMillisecond(time.Now().UTC()),
				UpdatedBy: userID,
			},
		}
		record.ApplyForm(&form)

		if err := c.encryptIdentity(&record); err != nil {
			return nil, fmt.Errorf("validate: %w", err)
		}

		return c.recordRepo.UpdateIdentity(&record)
	default:
		return nil, fmt.Errorf("%w: type assertion %T", pmerror.ErrInternal, record)
	}
}

func (c *Controller) DeleteRecord(id uuid.UUID, userID uuid.UUID) (*model.CredentialRecord, error) {
	if err := c.authorizeRecord(id, userID); err != nil {
		return nil, fmt.Errorf("authorize: %w", err)
	}

	return c.recordRepo.Delete(id)
}

func (c *Controller) authorizeRecord(id uuid.UUID, userID uuid.UUID) error {
	record, err := c.recordRepo.GetCredentialRecord(id)
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}

	if record.CreatedBy.String() != userID.String() {
		return fmt.Errorf("%w: user %s does not own record %s", pmerror.ErrForbidden, userID.String(), id.String())
	}

	return nil
}
