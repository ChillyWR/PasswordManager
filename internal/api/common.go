package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/okutsen/PasswordManager/internal/log"
)

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
		logger.Fatal("Failed to get path parameter: there is no id")
	}
	return uuid.Parse(idStr)
}

func readJSON(requestBody io.ReadCloser, out any) error {
	// TODO: prevent overflow (read by batches or set max size)
	recordsJSON, err := io.ReadAll(requestBody)
	defer requestBody.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(recordsJSON, out)
	if err != nil {
		return err
	}
	return err
}

func writeResponse(w http.ResponseWriter, body any, statusCode int, logger log.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		logger.Warnf("Failed to write JSON response: %s", err.Error())
	}
	// TODO: do not log private info
	logger.Debugf("Response written: %+v", body)
}
