package controller

import (
	"github.com/google/uuid"

	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/model/builder"
	"github.com/okutsen/PasswordManager/model/controller"
	"github.com/okutsen/PasswordManager/model/db"
)

type Repository interface {
	AllRecords() ([]db.Record, error)
	RecordByID(id uuid.UUID) (*db.Record, error)
	CreateRecord(record *db.Record) (*db.Record, error)
	UpdateRecord(record *db.Record) (*db.Record, error)
	DeleteRecord(id uuid.UUID) (*db.Record, error)

	AllUsers() ([]db.User, error)
	UserByID(id uuid.UUID) (*db.User, error)
	CreateUser(user *db.User) (*db.User, error)
	UpdateUser(user *db.User) (*db.User, error)
	DeleteUser(id uuid.UUID) (*db.User, error)
}

type Controller struct {
	repo Repository
	log  log.Logger
}

func New(logger log.Logger, ctrl Repository) *Controller {
	return &Controller{
		log:  logger.WithFields(log.Fields{"service": "Controller"}),
		repo: ctrl,
	}
}

func (c *Controller) AllRecords() ([]controller.Record, error) {
	getDBRecords, err := c.repo.AllRecords()
	if err != nil {
		return nil, err
	}
	records := builder.BuildControllerRecordsFromDBRecords(getDBRecords)

	return records, nil
}

func (c *Controller) Record(id uuid.UUID) (*controller.Record, error) {
	getRecord, err := c.repo.RecordByID(id)
	if err != nil {
		return nil, err
	}

	decPassword, err := Decrypt(getRecord.Password, Salt)
	getRecord.Password = decPassword

	record := builder.BuildControllerRecordFromDBRecord(getRecord)
	return &record, nil
}

// TODO: return specific errors to identify on api 404 Not found, 409 Conflict(if exists)

func (c *Controller) CreateRecord(record *controller.Record) (*controller.Record, error) {
	encPassword, err := Encrypt(record.Password, Salt)
	if err != nil {
		return nil, err
	}

	record.Password = encPassword

	dbRecord := builder.BuildDBRecordFromControllerRecord(record)
	createRecord, err := c.repo.CreateRecord(&dbRecord)
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

func (c *Controller) UpdateRecord(id uuid.UUID, record *controller.Record) (*controller.Record, error) {
	encPassword, err := Encrypt(record.Password, Salt)
	if err != nil {
		return nil, err
	}

	record.Password = encPassword

	dbRecord := builder.BuildDBRecordFromControllerRecord(record)
	dbRecord.ID = id

	updateRecord, err := c.repo.UpdateRecord(&dbRecord)
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

// 200, 404

func (c *Controller) DeleteRecord(id uuid.UUID) (*controller.Record, error) {
	dbRecord, err := c.repo.DeleteRecord(id)
	if err != nil {
		return nil, err
	}

	record := builder.BuildControllerRecordFromDBRecord(dbRecord)

	return &record, nil
}

func (c *Controller) AllUsers() ([]controller.User, error) {
	getUsers, err := c.repo.AllUsers()
	if err != nil {
		return nil, err
	}

	users := builder.BuildControllerUsersFromDBUsers(getUsers)
	return users, nil
}

func (c *Controller) User(id uuid.UUID) (*controller.User, error) {
	getUser, err := c.repo.UserByID(id) // TODO: pass uuid
	if err != nil {
		return nil, err
	}

	decPassword, err := Decrypt(getUser.Password, Salt)
	getUser.Password = decPassword

	user := builder.BuildControllerUserFromDBUser(getUser)
	return &user, nil
}

// TODO: return specific errors to identify on api 404 Not found, 409 Conflict(if exists)

func (c *Controller) CreateUser(user *controller.User) (*controller.User, error) {
	encPassword, err := Encrypt(user.Password, Salt)
	if err != nil {
		return nil, err
	}

	user.Password = encPassword

	dbUser := builder.BuildDBUserFromControllerUser(user)
	createdDBUser, err := c.repo.CreateUser(&dbUser)
	if err != nil {
		return nil, err
	}

	decPassword, err := Decrypt(createdDBUser.Password, Salt)
	if err != nil {
		return nil, err
	}

	createdDBUser.Password = decPassword

	createdUser := builder.BuildControllerUserFromDBUser(createdDBUser)
	return &createdUser, nil
}

// 200, 204(if no changes?), 404

func (c *Controller) UpdateUser(id uuid.UUID, user *controller.User) (*controller.User, error) {
	encPassword, err := Encrypt(user.Password, Salt)
	if err != nil {
		return nil, err
	}

	user.Password = encPassword

	dbUser := builder.BuildDBUserFromControllerUser(user)
	dbUser.ID = id

	updateUser, err := c.repo.UpdateUser(&dbUser)
	if err != nil {
		return nil, err
	}

	decPassword, err := Decrypt(updateUser.Password, Salt)
	if err != nil {
		return nil, err
	}

	updateUser.Password = decPassword

	updatedUser := builder.BuildControllerUserFromDBUser(updateUser)
	return &updatedUser, nil
}

// 200, 404

func (c *Controller) DeleteUser(id uuid.UUID) (*controller.User, error) {
	dbUser, err := c.repo.DeleteUser(id)
	if err != nil {
		return nil, err
	}

	user := builder.BuildControllerUserFromDBUser(dbUser)

	return &user, nil
}
