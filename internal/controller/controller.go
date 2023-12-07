package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/model"
	"github.com/okutsen/PasswordManager/pkg/pmerror"
)

type RecordRepository interface {
	GetAll() ([]model.CredentialRecord, error)
	Get(id uuid.UUID) (*model.CredentialRecord, error)
	Create(record *model.CredentialRecord) (*model.CredentialRecord, error)
	Update(record *model.CredentialRecord) (*model.CredentialRecord, error)
	Delete(id uuid.UUID) (*model.CredentialRecord, error)
}

type UserRepository interface {
	GetAll() ([]model.User, error)
	Get(id uuid.UUID) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(id uuid.UUID) (*model.User, error)
}

type Controller struct {
	userRepo   UserRepository
	recordRepo RecordRepository
	log        log.Logger
}

func New(logger log.Logger, userRepo UserRepository, recordRepo RecordRepository) (*Controller, error) {
	if userRepo == nil {
		return nil, errors.New("userRepo is nil")
	}

	if recordRepo == nil {
		return nil, errors.New("recordRepo is nil")
	}

	return &Controller{
		log:        logger.WithFields(log.Fields{"service": "Controller"}),
		userRepo:   userRepo,
		recordRepo: recordRepo,
	}, nil
}

func (c *Controller) AllRecords() ([]model.CredentialRecord, error) {
	return c.recordRepo.GetAll()
}

func (c *Controller) CredentialRecord(id uuid.UUID) (*model.CredentialRecord, error) {
	record, err := c.recordRepo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	decNotes, err := Decrypt(*record.Notes, Salt)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	record.Notes = &decNotes

	return record, nil
}

func (c *Controller) CreateRecord(form *model.CredentialRecordForm) (*model.CredentialRecord, error) {
	if err := form.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	encNotes, err := Encrypt(*form.Notes, Salt)
	if err != nil {
		return nil, err
	}

	// TODO: fill user info
	record := model.CredentialRecord{
		Name:  *form.Name,
		Notes: &encNotes,
	}

	return c.recordRepo.Create(&record)
}

func (c *Controller) UpdateRecord(id uuid.UUID, form *model.CredentialRecordForm) (*model.CredentialRecord, error) {
	if form.Empty() {
		return nil, fmt.Errorf("%w: empty form", pmerror.ErrInvalidInput)
	}

	record := model.CredentialRecord{
		ID:        id,
		UpdatedOn: time.Now().UTC(),
	}
	// TODO: validate that if Name is nil repo wont update name to empty string
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

func (c *Controller) DeleteRecord(id uuid.UUID) (*model.CredentialRecord, error) {
	return c.recordRepo.Delete(id)
}

func (c *Controller) AllUsers() ([]model.User, error) {
	return c.userRepo.GetAll()
}

func (c *Controller) User(id uuid.UUID) (*model.User, error) {
	repoUser, err := c.userRepo.Get(id)
	if err != nil {
		return nil, err
	}

	decPassword, err := Decrypt(repoUser.Password, Salt)
	if err != nil {
		return nil, err
	}

	repoUser.Password = decPassword

	return repoUser, nil
}

func (c *Controller) CreateUser(form *model.UserForm) (*model.User, error) {
	if err := form.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	encPassword, err := Encrypt(*form.Password, Salt)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Name:     *form.Name,
		Password: encPassword,
	}

	result, err := c.userRepo.Create(&user)
	if err != nil {
		return nil, err
	}

	result.Password, err = Decrypt(result.Password, Salt)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Controller) UpdateUser(id uuid.UUID, form *model.UserForm) (*model.User, error) {
	if form.Empty() {
		return nil, fmt.Errorf("%w: empty form", pmerror.ErrInvalidInput)
	}

	user := model.User{
		ID:        id,
		UpdatedOn: time.Now().UTC(),
	}

	if form.Name != nil {
		user.Name = *form.Name
	}

	if form.Password != nil {
		encPassword, err := Encrypt(*form.Password, Salt)
		if err != nil {
			return nil, fmt.Errorf("encrypt: %w", err)
		}

		user.Password = encPassword
	}

	return c.userRepo.Update(&user)
}

func (c *Controller) DeleteUser(id uuid.UUID) (*model.User, error) {
	return c.userRepo.Delete(id)
}
