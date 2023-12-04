package api

import (
	"time"

	"github.com/google/uuid"
)

const (
	InvalidJSONMessage     = "Invalid JSON"
	InvalidRecordIDMessage = "Invalid record ID"
	InvalidUserIDMessage   = "Invalid user ID"
	InternalErrorMessage   = "Oops, something went wrong"
	UnAuthorizedMessage    = "Sign in to use service"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CredentialRecord struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Login       string    `json:"login"`
	Password    string    `json:"password"`
	URL         string    `json:"url,omitempty"`
	Description string    `json:"description,omitempty"`
	UpdatedBy   string    `json:"updated_by,omitempty"`
	CreatedBy   string    `json:"created_by,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type Error struct {
	Message string `json:"message"`
}
