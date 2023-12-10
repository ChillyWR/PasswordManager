package api

import (
	"net/http"

	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/model"
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

		result, err := apictx.ctrl.AllRecords()
		if err != nil {
			logger.Errorf("Failed to list records: %s", err.Error())
			writeError(w, err, logger)
			return
		}

		writeResponse(w, result, http.StatusOK, logger)
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
			logger.Errorf("Invalid record id: %s", err.Error())
			writeResponse(w, Error{Message: InvalidRecordIDMessage}, http.StatusBadRequest, logger)
			return
		}

		result, err := apictx.ctrl.CredentialRecord(recordID, rctx.userID)
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

		var payload model.CredentialRecordForm
		if err := readBody(r.Body, &payload); err != nil {
			logger.Errorf("Failed to read body: %s", err.Error())
			writeResponse(w, Error{Message: InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}

		result, err := apictx.ctrl.CreateRecord(&payload, rctx.userID)
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
			logger.Errorf("Invalid record id: %s", err.Error())
			writeResponse(w, Error{Message: InvalidRecordIDMessage}, http.StatusBadRequest, logger)
			return
		}

		var payload model.CredentialRecordForm
		if err = readBody(r.Body, &payload); err != nil {
			logger.Errorf("Failed to read body: %s", err.Error())
			writeResponse(w, Error{Message: InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}

		result, err := apictx.ctrl.UpdateRecord(recordID, &payload, rctx.userID)
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
			logger.Warnf("Invalid record id: %s", err.Error())
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
