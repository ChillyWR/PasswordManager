package repo

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/okutsen/PasswordManager/model/db"
)

func NewUserRepository(db *gorm.DB) (*UserRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &UserRepository{db: db}, nil
}

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) GetAll() ([]db.User, error) {
	var user []db.User
	result := r.db.Find(&user)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all user from db: %w", err)
	}

	return user, nil
}
func (r *UserRepository) Get(id uuid.UUID) (*db.User, error) {
	var user db.User
	result := r.db.First(&user, id)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user from db: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Create(user *db.User) (*db.User, error) {
	user.ID = uuid.New()
	result := r.db.Create(user)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user in db: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Update(user *db.User) (*db.User, error) {
	result := r.db.Model(user).Clauses(clause.Returning{}).Updates(user)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to update user in db: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Delete(id uuid.UUID) (*db.User, error) {
	var user db.User
	result := r.db.Model(&user).Clauses(clause.Returning{}).Delete(&user, id)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to remove user from db: %w", err)
	}
	return &user, nil
}
