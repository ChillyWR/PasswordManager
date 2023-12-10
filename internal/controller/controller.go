package controller

import (
	"errors"

	"github.com/google/uuid"

	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/model"
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

func New(userRepo UserRepository, recordRepo RecordRepository, logger log.Logger) (*Controller, error) {
	if userRepo == nil {
		return nil, errors.New("userRepo is nil")
	}

	if recordRepo == nil {
		return nil, errors.New("recordRepo is nil")
	}

	return &Controller{
		userRepo:   userRepo,
		recordRepo: recordRepo,
		log:        logger.WithFields(log.Fields{"module": "Controller"}),
	}, nil
}
