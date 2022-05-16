package apischema

import "github.com/google/uuid"

const (
	InvalidJSONMessage   = "Invalid JSON"
	InternalErrorMessage = "Oops, something went wrong"
)

// TODO: add validator
// TODO: add uuid (request id)
type Record struct {
	ID       uuid.UUID `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Error struct {
	Message string `json:"message"`
}
