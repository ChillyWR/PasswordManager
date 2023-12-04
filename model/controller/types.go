package controller

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Login     string
	Password  string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CredentialRecord struct {
	ID          uuid.UUID
	Name        string
	Login       string
	Password    string
	URL         string
	Description string
	CreatedBy   string
	UpdatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
