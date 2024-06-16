package controller

import (
	"errors"

	"github.com/google/uuid"

	"github.com/ChillyWR/PasswordManager/internal/log"
	"github.com/ChillyWR/PasswordManager/model"
)

const (
	Salt = "abc&1*~#^2^#s0^=)^^7%b34"
)

type RecordRepository interface {
	GetAll(userID uuid.UUID) ([]model.CredentialRecord, []model.LoginRecord, []model.CardRecord, []model.IdentityRecord, error)
	GetCredentialRecord(id uuid.UUID) (*model.CredentialRecord, error)
	GetLogin(id uuid.UUID) (*model.LoginRecord, error)
	GetCard(id uuid.UUID) (*model.CardRecord, error)
	GetIdentity(id uuid.UUID) (*model.IdentityRecord, error)
	CreateCredentialRecord(record *model.CredentialRecord) (*model.CredentialRecord, error)
	CreateLogin(record *model.LoginRecord) (*model.LoginRecord, error)
	CreateCard(record *model.CardRecord) (*model.CardRecord, error)
	CreateIdentity(record *model.IdentityRecord) (*model.IdentityRecord, error)
	UpdateCredentialRecord(record *model.CredentialRecord) (*model.CredentialRecord, error)
	UpdateLogin(record *model.LoginRecord) (*model.LoginRecord, error)
	UpdateCard(record *model.CardRecord) (*model.CardRecord, error)
	UpdateIdentity(record *model.IdentityRecord) (*model.IdentityRecord, error)
	Delete(id uuid.UUID) (*model.CredentialRecord, error)
}

type UserRepository interface {
	GetAll() ([]model.User, error)
	Get(id uuid.UUID) (*model.User, error)
	GetByName(name string) (*model.User, error)
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
