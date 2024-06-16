package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/ChillyWR/PasswordManager/internal/log"
	"github.com/ChillyWR/PasswordManager/pkg/pmerror"
)

const (
	// PPN: Path Parameter Name
	// HPN: Header Parameter Name
	IDPPN                 = "id"
	CorrelationIDHPN      = "X-Request-ID"
	AuthorizationTokenHPN = "Authorization"

	RequestContextName = "rctx"
)

const (
	InvalidJSONMessage     = "Invalid JSON"
	InvalidRecordIDMessage = "Invalid record ID"
	InvalidUserIDMessage   = "Invalid user ID"
	InternalErrorMessage   = "Oops, something went wrong"
	UnAuthorizedMessage    = "Sign in to use service"
)

type Error struct {
	Message string `json:"message"`
}

// unpackRequestContext gets and validates RequestContext from ctx
func unpackRequestContext(ctx context.Context, logger log.Logger) *RequestContext {
	rctx, ok := ctx.Value(RequestContextName).(*RequestContext)
	if !ok {
		logger.Fatalf("Failed to unpack request context, got: %s", rctx)
	}

	return rctx
}

// getIDFrom checks if id is set and returns the result of uuid parsing
func getIDFrom(ps httprouter.Params, logger log.Logger) (uuid.UUID, error) {
	idStr := ps.ByName(IDPPN)
	if idStr == "" {
		logger.Fatal("Failed to get path parameter")
	}

	return uuid.Parse(idStr)
}

func readBody(body io.ReadCloser, v any) error {
	raw, err := io.ReadAll(body) // TODO: prevent potential overflow
	defer body.Close()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw, v); err != nil {
		return err
	}

	return err
}

func writeResponse(w http.ResponseWriter, body any, statusCode int, logger log.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if body != nil {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			logger.Errorf("Failed to write JSON response: %s", err.Error())
		}
	}
	// do not log private info
	// logger.Debugf("Response written: %+v", body)
}

func writeError(w http.ResponseWriter, err error, logger log.Logger) {
	writeResponse(w, nil, errorStatus(err), logger)
}

func errorStatus(err error) int {
	switch {
	case errors.Is(err, pmerror.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, pmerror.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, pmerror.ErrForbidden):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
