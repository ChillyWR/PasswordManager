package controller

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/okutsen/PasswordManager/model"
	"github.com/okutsen/PasswordManager/pkg/pmcrypto"
	"github.com/okutsen/PasswordManager/pkg/pmerror"
)

func (c *Controller) Login(form *model.UserForm) (uuid.UUID, error) {
	if err := form.Validate(); err != nil {
		return uuid.UUID{}, fmt.Errorf("validate: %w", err)
	}

	user, err := c.userRepo.GetByName(*form.Name)
	if err != nil {
		return uuid.UUID{},fmt.Errorf("validate: %w", err)
	}

	decPassword, err := pmcrypto.Decrypt(user.Password, Salt)
	if err != nil {
		return uuid.UUID{},fmt.Errorf("decrypt: %w", err)
	}

	if *form.Password != decPassword {
		return uuid.UUID{},fmt.Errorf("%w: password doesn't match", pmerror.ErrInvalidInput)
	}

	return user.ID, nil
}

func (c *Controller) AllUsers() ([]model.User, error) {
	return c.userRepo.GetAll()
}

func (c *Controller) GetUser(id uuid.UUID) (*model.User, error) {
	repoUser, err := c.userRepo.Get(id)
	if err != nil {
		return nil, err
	}

	v, err := pmcrypto.Decrypt(repoUser.Password, Salt)
	if err != nil {
		return nil, err
	}

	repoUser.Password = v

	return repoUser, nil
}

func (c *Controller) CreateUser(form *model.UserForm) (*model.User, error) {
	if err := form.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	encPassword, err := pmcrypto.Encrypt(*form.Password, Salt)
	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:        uuid.New(),
		Name:      *form.Name,
		Password:  encPassword,
		CreatedOn: time.Now().UTC(),
		UpdatedOn: time.Now().UTC(),
	}

	result, err := c.userRepo.Create(&user)
	if err != nil {
		return nil, err
	}

	result.Password, err = pmcrypto.Decrypt(result.Password, Salt)
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
		v, err := pmcrypto.Encrypt(*form.Password, Salt)
		if err != nil {
			return nil, fmt.Errorf("encrypt: %w", err)
		}

		user.Password = v
	}

	return c.userRepo.Update(&user)
}

func (c *Controller) DeleteUser(id uuid.UUID) (*model.User, error) {
	return c.userRepo.Delete(id)
}
