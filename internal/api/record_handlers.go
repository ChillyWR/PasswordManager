package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ChillyWR/PasswordManager/internal/log"
	"github.com/ChillyWR/PasswordManager/model"
)

func NewListRecordsHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(log.Fields{
		"handler": "ListRecords",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
			"cor_id": rctx.corID.String(),
		})

		secureNotes, logins, cards, identities, err := apictx.ctrl.AllRecords(rctx.userID)
		if err != nil {
			logger.Errorf("Failed to list records: %s", err.Error())
			writeError(w, err, logger)
			return
		}

		writeResponse(w, struct {
			SecureNotes []model.CredentialRecord `json:"secure_notes"`
			Logins      []model.LoginRecord      `json:"logins"`
			Cards       []model.CardRecord       `json:"cards"`
			Identities  []model.IdentityRecord   `json:"identities"`
		}{
			SecureNotes: secureNotes,
			Logins:      logins,
			Cards:       cards,
			Identities:  identities,
		}, http.StatusOK, logger)
	}
}

func NewGetRecordHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(log.Fields{
		"handler": "GetRecord",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
			"cor_id": rctx.corID.String(),
		})

		recordID, err := getIDFrom(rctx.params, logger)
		if err != nil {
			logger.Errorf("%s: %s", InvalidRecordIDMessage, err.Error())
			writeResponse(w, Error{Message: InvalidRecordIDMessage}, http.StatusBadRequest, logger)
			return
		}

		logger.Infof("Check auth for %v", rctx.userID)
		result, err := apictx.ctrl.GetRecord(recordID, rctx.userID)
		if err != nil {
			logger.Errorf("Failed to get record: %s", err.Error())
			writeError(w, err, logger)
			return
		}

		writeResponse(w, result, http.StatusOK, logger)
	}
}

func NewCreateRecordHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(log.Fields{
		"handler": "CreateRecord",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
			"cor_id": rctx.corID.String(),
		})

		var payload struct {
			Type string          `json:"type"`
			Form json.RawMessage `json:"form"`
		}
		if err := readBody(r.Body, &payload); err != nil {
			logger.Errorf("Failed to read body: %s", err.Error())
			writeResponse(w, Error{Message: InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}

		result, err := apictx.ctrl.CreateRecord(model.RecordType(payload.Type), payload.Form, rctx.userID)
		if err != nil {
			logger.Errorf("Failed to create record: %s", err.Error())
			writeError(w, err, logger)
			return
		}

		writeResponse(w, result, http.StatusCreated, logger)
	}
}

func NewUpdateRecordHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(log.Fields{
		"handler": "UpdateRecord",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
			"cor_id": rctx.corID.String(),
		})

		recordID, err := getIDFrom(rctx.params, logger)
		if err != nil {
			logger.Errorf("%s: %s", InvalidRecordIDMessage, err.Error())
			writeResponse(w, Error{Message: InvalidRecordIDMessage}, http.StatusBadRequest, logger)
			return
		}

		raw, err := io.ReadAll(r.Body) // TODO: prevent potential overflow
		defer r.Body.Close()
		if err != nil {
			logger.Errorf("Failed to read body: %s", err.Error())
			writeResponse(w, Error{Message: InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}

		result, err := apictx.ctrl.UpdateRecord(recordID, raw, rctx.userID)
		if err != nil {
			logger.Errorf("Failed to update record: %s", err.Error())
			writeError(w, err, logger)
			return
		}

		writeResponse(w, result, http.StatusAccepted, logger)
	}
}

func NewDeleteRecordHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(log.Fields{
		"handler": "DeleteRecord",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
			"cor_id": rctx.corID.String(),
		})

		recordID, err := getIDFrom(rctx.params, logger)
		if err != nil {
			logger.Warnf("%s: %s", InvalidRecordIDMessage, err.Error())
			writeResponse(w, Error{Message: InvalidRecordIDMessage}, http.StatusBadRequest, logger)
			return
		}

		result, err := apictx.ctrl.DeleteRecord(recordID, rctx.userID)
		if err != nil {
			logger.Errorf("Failed to delete record: %s", err.Error())
			writeError(w, err, logger)
			return
		}

		writeResponse(w, result, http.StatusOK, logger)
	}
}
