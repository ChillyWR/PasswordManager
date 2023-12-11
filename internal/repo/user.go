package repo

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/okutsen/PasswordManager/model"
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
		return nil, fmt.Errorf("failed to get all user: %w", err)
	}

	return user, nil
}
func (r *UserRepository) Get(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Create(user *model.User) (*model.User, error) {
	core := User(*user)
	if err := r.db.Create(&core).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = core.ID

	return user, nil
}

func (r *UserRepository) Update(user *model.User) (*model.User, error) {
	if err := r.db.Model(user).Clauses(clause.Returning{}).Updates(user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Delete(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.Model(&user).Clauses(clause.Returning{}).Delete(&user, id).Error; err != nil {
		return nil, fmt.Errorf("failed to remove user: %w", err)
	}
	return &user, nil
}
