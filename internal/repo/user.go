package repo

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/okutsen/PasswordManager/model"
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

func (r *UserRepository) GetAll() ([]model.User, error) {
	var user []model.User
	result := r.db.Find(&user)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("failed to get all user from db: %w", err)
	}

	return user, nil
}
func (r *UserRepository) Get(id uuid.UUID) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, id)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("failed to get user from db: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Create(user *model.User) (*model.User, error) {
	result := r.db.Create(user)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("failed to create user in db: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Update(user *model.User) (*model.User, error) {
	result := r.db.Model(user).Clauses(clause.Returning{}).Updates(user)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("failed to update user in db: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Delete(id uuid.UUID) (*model.User, error) {
	var user model.User
	result := r.db.Model(&user).Clauses(clause.Returning{}).Delete(&user, id)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("failed to remove user from db: %w", err)
	}
	return &user, nil
}
