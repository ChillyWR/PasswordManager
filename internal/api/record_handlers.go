package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/model/api"
	"github.com/okutsen/PasswordManager/model/builder"
)

func NewListRecordsHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(log.Fields{"handler": "GetAllRecords"})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
			"cor_id": rctx.corID.String(),
		})
		records, err := apictx.ctrl.AllRecords()
		if err != nil {
			logger.Warnf("Failed to get records from controller: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		recordsAPI := builder.BuildAPIRecordsFromControllerRecords(records)
		// Write JSON by stream?
		writeResponse(w, recordsAPI, http.StatusOK, logger)
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
			logger.Warnf("Invalid record id: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InvalidRecordIDMessage}, http.StatusBadRequest, logger)
			return
		}
		record, err := apictx.ctrl.CredentialRecord(recordID)
		if err != nil {
			logger.Warnf("Failed to get records from controller: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		// TODO: get record from db
		writeResponse(w, builder.BuildAPIRecordFromControllerRecord(record), http.StatusOK, logger)
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
		var recordAPI *api.CredentialRecord
		err := readJSON(r.Body, &recordAPI)
		defer r.Body.Close()
		if err != nil {
			logger.Warnf("Failed to read JSON: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}
		record := builder.BuildControllerRecordFromAPIRecord(recordAPI)
		// TODO: if exists return err (409 Conflict)
		resultRecord, err := apictx.ctrl.CreateRecord(&record)
		if err != nil {
			logger.Warnf("Failed to get records from controller: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		// TODO: get record from db
		writeResponse(w, builder.BuildAPIRecordFromControllerRecord(resultRecord), http.StatusCreated, logger)
	}
}

func NewUpdateRecordHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(log.Fields{
		"handler": "UpdateRecords",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
			"cor_id": rctx.corID.String(),
		})
		recordID, err := getIDFrom(rctx.params, logger)
		if err != nil {
			logger.Warnf("Invalid record id: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InvalidRecordIDMessage}, http.StatusBadRequest, logger)
			return
		}
		var recordAPI *api.CredentialRecord
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Warnf("Failed to read JSON: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}
		defer r.Body.Close()
		err = json.Unmarshal(body, &recordAPI)
		if err != nil {
			logger.Warnf("failed to unmarshal JSON file: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InternalErrorMessage}, http.StatusBadRequest, logger)
			return
		}
		if recordID != recordAPI.ID {
			logger.Warn("Record id from path parameter doesn't match id from new record structure")
			writeResponse(w,
				api.Error{Message: api.InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}
		record := builder.BuildControllerRecordFromAPIRecord(recordAPI)
		resultRecord, err := apictx.ctrl.UpdateRecord(recordID, &record)
		if err != nil {
			logger.Warnf("Failed to get records from controller: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		// TODO: get record from db
		writeResponse(w,
			builder.BuildAPIRecordFromControllerRecord(resultRecord), http.StatusAccepted, logger)
	}
}

func NewDeleteRecordHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(log.Fields{
		"handler": "DeleteRecords",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
			"cor_id": rctx.corID.String(),
		})
		recordID, err := getIDFrom(rctx.params, logger)
		if err != nil {
			logger.Warnf("Invalid record id: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InvalidRecordIDMessage}, http.StatusBadRequest, logger)
			return
		}
		resultRecord, err := apictx.ctrl.DeleteRecord(recordID)
		if err != nil {
			logger.Errorf("Failed to get records from controller: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		writeResponse(w,
			builder.BuildAPIRecordFromControllerRecord(resultRecord), http.StatusOK, logger)
	}
}
