package repo

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"

	"github.com/okutsen/PasswordManager/model/db"
)

func (r *Repo) AllUsers() ([]db.User, error) {
	var user []db.User
	result := r.db.Find(&user)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all user from db: %w", err)
	}

	return user, nil
}
func (r *Repo) UserByID(id uuid.UUID) (*db.User, error) {
	var user db.User
	result := r.db.First(&user, id)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user from db: %w", err)
	}

	return &user, nil
}

func (r *Repo) CreateUser(user *db.User) (*db.User, error) {
	user.ID = uuid.New()
	result := r.db.Create(user)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user in db: %w", err)
	}

	return user, nil
}

func (r *Repo) UpdateUser(user *db.User) (*db.User, error) {
	result := r.db.Model(user).Clauses(clause.Returning{}).Updates(user)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to update user in db: %w", err)
	}

	return user, nil
}

func (r *Repo) DeleteUser(id uuid.UUID) (*db.User, error) {
	var user db.User
	result := r.db.Model(&user).Clauses(clause.Returning{}).Delete(&user, id)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to remove user from db: %w", err)
	}
	return &user, nil
}
