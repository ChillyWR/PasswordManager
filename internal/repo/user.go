package repo

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/okutsen/PasswordManager/model"
	"github.com/okutsen/PasswordManager/pkg/pmerror"
)

type User model.User

func (User) TableName() string {
	return "reg_user"
}

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
	if err := r.db.Find(&user).Error; err != nil {
		return nil, fmt.Errorf("find: %w", convertError(err))
	}

	return user, nil
}
func (r *UserRepository) Get(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("first: %w", convertError(err))
	}

	return &user, nil
}

func (r *UserRepository) Create(user *model.User) (*model.User, error) {
	core := User(*user)
	if err := r.db.Create(&core).Error; err != nil {
		return nil, fmt.Errorf("create: %w", convertError(err))
	}

	return (*model.User)(&core), nil
}

func (r *UserRepository) Update(user *model.User) (*model.User, error) {
	result := r.db.Model(user).Clauses(clause.Returning{}).Updates(user)
	if result.Error != nil {
		return nil, fmt.Errorf("updates: %w", convertError(result.Error))
	}

	if result.RowsAffected == 0 {
		return nil, pmerror.ErrNotFound
	}

	return user, nil
}

func (r *UserRepository) Delete(id uuid.UUID) (*model.User, error) {
	var user User
	result := r.db.Clauses(clause.Returning{}).Where("id = ?", id).Delete(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("delete: %w", convertError(result.Error))
	}

	if result.RowsAffected == 0 {
		return nil, pmerror.ErrNotFound
	}

	return (*model.User)(&user), nil
}
