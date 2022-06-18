package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/schema/apischema"
	"github.com/okutsen/PasswordManager/schema/schemabuilder"
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
				apischema.Error{Message: apischema.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		recordsAPI := schemabuilder.BuildAPIRecordsFromControllerRecords(records)
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
		recordID, err := getIDFrom(rctx.ps, logger)
		if err != nil {
			logger.Warnf("Invalid record id: %s", err.Error())
			writeResponse(w,
				apischema.Error{Message: apischema.InvalidRecordIDMessage}, http.StatusBadRequest, logger)
			return
		}
		record, err := apictx.ctrl.Record(recordID)
		if err != nil {
			logger.Warnf("Failed to get records from controller: %s", err.Error())
			writeResponse(w,
				apischema.Error{Message: apischema.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		// TODO: get record from db
		writeResponse(w, schemabuilder.BuildAPIRecordFromControllerRecord(record), http.StatusOK, logger)
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
		var recordAPI *apischema.Record
		err := readJSON(r.Body, &recordAPI)
		defer r.Body.Close()
		if err != nil {
			logger.Warnf("Failed to read JSON: %s", err.Error())
			writeResponse(w,
				apischema.Error{Message: apischema.InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}
		record := schemabuilder.BuildControllerRecordFromAPIRecord(recordAPI)
		// TODO: if exists return err (409 Conflict)
		resultRecord, err := apictx.ctrl.CreateRecord(&record)
		if err != nil {
			logger.Warnf("Failed to get records from controller: %s", err.Error())
			writeResponse(w,
				apischema.Error{Message: apischema.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		// TODO: get record from db
		writeResponse(w, schemabuilder.BuildAPIRecordFromControllerRecord(resultRecord), http.StatusCreated, logger)
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
		recordID, err := getIDFrom(rctx.ps, logger)
		if err != nil {
			logger.Warnf("Invalid record id: %s", err.Error())
			writeResponse(w,
				apischema.Error{Message: apischema.InvalidRecordIDMessage}, http.StatusBadRequest, logger)
			return
		}
		var recordAPI *apischema.Record
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Warnf("Failed to read JSON: %s", err.Error())
			writeResponse(w,
				apischema.Error{Message: apischema.InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}
		defer r.Body.Close()
		err = json.Unmarshal(body, &recordAPI)
		if err != nil {
			logger.Warnf("failed to unmarshal JSON file: %s", err.Error())
			writeResponse(w,
				apischema.Error{Message: apischema.InternalErrorMessage}, http.StatusBadRequest, logger)
			return
		}
		if recordID != recordAPI.ID {
			logger.Warn("Record id from path parameter doesn't match id from new record structure")
			writeResponse(w,
				apischema.Error{Message: apischema.InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}
		record := schemabuilder.BuildControllerRecordFromAPIRecord(recordAPI)
		resultRecord, err := apictx.ctrl.UpdateRecord(recordID, &record)
		if err != nil {
			logger.Warnf("Failed to get records from controller: %s", err.Error())
			writeResponse(w,
				apischema.Error{Message: apischema.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		// TODO: get record from db
		writeResponse(w,
			schemabuilder.BuildAPIRecordFromControllerRecord(resultRecord), http.StatusAccepted, logger)
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
		recordID, err := getIDFrom(rctx.ps, logger)
		if err != nil {
			logger.Warnf("Invalid record id: %s", err.Error())
			writeResponse(w,
				apischema.Error{Message: apischema.InvalidRecordIDMessage}, http.StatusBadRequest, logger)
			return
		}
		resultRecord, err := apictx.ctrl.DeleteRecord(recordID)
		if err != nil {
			logger.Errorf("Failed to get records from controller: %s", err.Error())
			writeResponse(w,
				apischema.Error{Message: apischema.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		writeResponse(w,
			schemabuilder.BuildAPIRecordFromControllerRecord(resultRecord), http.StatusOK, logger)
	}
}
