package controller

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/model/builder"
	"github.com/okutsen/PasswordManager/model/controller"
	"github.com/okutsen/PasswordManager/model/db"
)

type CredentialRecordRepository interface {
	GetAll() ([]db.CredentialRecord, error)
	Get(id uuid.UUID) (*db.CredentialRecord, error)
	Create(record *db.CredentialRecord) (*db.CredentialRecord, error)
	Update(record *db.CredentialRecord) (*db.CredentialRecord, error)
	Delete(id uuid.UUID) (*db.CredentialRecord, error)
}

type UserRepository interface {
	GetAll() ([]db.User, error)
	Get(id uuid.UUID) (*db.User, error)
	Create(user *db.User) (*db.User, error)
	Update(user *db.User) (*db.User, error)
	Delete(id uuid.UUID) (*db.User, error)
}

type Controller struct {
	userRepo         UserRepository
	credentialRecord CredentialRecordRepository
	log              log.Logger
}

func New(logger log.Logger, userRepo UserRepository, credentialRecord CredentialRecordRepository) (*Controller, error) {
	if userRepo == nil {
		return nil, errors.New("userRepo is nil")
	}

	if credentialRecord == nil {
		return nil, errors.New("credentialRecord is nil")
	}
	
	return &Controller{
		log:              logger.WithFields(log.Fields{"service": "Controller"}),
		userRepo:         userRepo,
		credentialRecord: credentialRecord,
	}, nil
}

func (c *Controller) AllRecords() ([]controller.CredentialRecord, error) {
	getDBRecords, err := c.credentialRecord.GetAll()
	if err != nil {
		return nil, err
	}
	records := builder.BuildControllerRecordsFromDBRecords(getDBRecords)

	return records, nil
}

func (c *Controller) CredentialRecord(id uuid.UUID) (*controller.CredentialRecord, error) {
	repoRecord, err := c.credentialRecord.Get(id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}

	decPassword, err := Decrypt(repoRecord.Password, Salt)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	repoRecord.Password = decPassword

	record := builder.BuildControllerRecordFromDBRecord(repoRecord)
	return &record, nil
}

func (c *Controller) CreateRecord(record *controller.CredentialRecord) (*controller.CredentialRecord, error) {
	encPassword, err := Encrypt(record.Password, Salt)
	if err != nil {
		return nil, err
	}

	record.Password = encPassword

	dbRecord := builder.BuildDBRecordFromControllerRecord(record)
	createRecord, err := c.credentialRecord.Create(&dbRecord)
	if err != nil {
		return nil, err
	}

	decPassword, err := Decrypt(createRecord.Password, Salt)
	if err != nil {
		return nil, err
	}

	createRecord.Password = decPassword

	createdRecord := builder.BuildControllerRecordFromDBRecord(createRecord)
	return &createdRecord, nil
}

// 200, 204(if no changes?), 404

func (c *Controller) UpdateRecord(id uuid.UUID, record *controller.CredentialRecord) (*controller.CredentialRecord, error) {
	encPassword, err := Encrypt(record.Password, Salt)
	if err != nil {
		return nil, err
	}

	record.Password = encPassword

	dbRecord := builder.BuildDBRecordFromControllerRecord(record)
	dbRecord.ID = id

	updateRecord, err := c.credentialRecord.Update(&dbRecord)
	if err != nil {
		return nil, err
	}

	decPassword, err := Decrypt(updateRecord.Password, Salt)
	if err != nil {
		return nil, err
	}

	updateRecord.Password = decPassword

	updatedRecord := builder.BuildControllerRecordFromDBRecord(updateRecord)
	return &updatedRecord, nil
}

func (c *Controller) DeleteRecord(id uuid.UUID) (*controller.CredentialRecord, error) {
	dbRecord, err := c.credentialRecord.Delete(id)
	if err != nil {
		return nil, err
	}

	record := builder.BuildControllerRecordFromDBRecord(dbRecord)

	return &record, nil
}

func (c *Controller) AllUsers() ([]controller.User, error) {
	repoUsers, err := c.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	users := builder.BuildControllerUsersFromRepoUsers(repoUsers)
	return users, nil
}

func (c *Controller) User(id uuid.UUID) (*controller.User, error) {
	repoUser, err := c.userRepo.Get(id)
	if err != nil {
		return nil, err
	}

	decPassword, err := Decrypt(repoUser.Password, Salt)
	if err != nil {
		return nil, err
	}

	repoUser.Password = decPassword

	user := builder.BuildControllerUserFromRepoUser(repoUser)
	return &user, nil
}

func (c *Controller) CreateUser(user *controller.User) (*controller.User, error) {
	encPassword, err := Encrypt(user.Password, Salt)
	if err != nil {
		return nil, err
	}

	user.Password = encPassword

	dbUser := builder.BuildDBUserFromControllerUser(user)
	createdDBUser, err := c.userRepo.Create(&dbUser)
	if err != nil {
		return nil, err
	}

	decPassword, err := Decrypt(createdDBUser.Password, Salt)
	if err != nil {
		return nil, err
	}

	createdDBUser.Password = decPassword

	createdUser := builder.BuildControllerUserFromRepoUser(createdDBUser)
	return &createdUser, nil
}

func (c *Controller) UpdateUser(id uuid.UUID, user *controller.User) (*controller.User, error) {
	encPassword, err := Encrypt(user.Password, Salt)
	if err != nil {
		return nil, err
	}

	user.Password = encPassword

	dbUser := builder.BuildDBUserFromControllerUser(user)
	dbUser.ID = id

	updateUser, err := c.userRepo.Update(&dbUser)
	if err != nil {
		return nil, err
	}

	decPassword, err := Decrypt(updateUser.Password, Salt)
	if err != nil {
		return nil, err
	}

	updateUser.Password = decPassword

	updatedUser := builder.BuildControllerUserFromRepoUser(updateUser)
	return &updatedUser, nil
}

func (c *Controller) DeleteUser(id uuid.UUID) (*controller.User, error) {
	dbUser, err := c.userRepo.Delete(id)
	if err != nil {
		return nil, err
	}

	user := builder.BuildControllerUserFromRepoUser(dbUser)

	return &user, nil
}
