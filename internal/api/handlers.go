package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/schema/apischema"
	"github.com/okutsen/PasswordManager/schema/schemabuilder"
)

const (
	RecordCreatedMessage = "Record created"
)

func NewEndpointLoggerMiddleware(ctx *APIContext, handler httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx.logger.Infof("API: Endpoint Hit: %s %s%s\n", r.Host, r.URL.Path, r.Method)
		handler(rw, r, ps)
	}
}

func NewGetAllRecordsHandler(ctx *APIContext) httprouter.Handle {
	logger := ctx.logger.WithFields(log.Fields{"handler": "getAllRecords"})
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		records, err := ctx.ctrl.GetAllRecords()
		if err != nil {
			logger.Warnf("failed to get records from controller: %s", err.Error())
			writeJSONResponse(w, logger, apischema.Error{Message: "Failed to receive data from controller"}, http.StatusInternalServerError)
			return
		}
		recordsAPI := schemabuilder.BuildRecordsAPIFrom(records)
		writeJSONResponse(w, logger, recordsAPI, http.StatusOK)
	}
}

func NewGetRecordHandler(ctx *APIContext) httprouter.Handle {
	logger := ctx.logger.WithFields(log.Fields{"handler": "getRecord"})
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		idStr := ps.ByName(IDParamName)
		idInt, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			logger.Warnf("failed to convert path parameter id: %s", err.Error())
			writeJSONResponse(w, logger, apischema.Error{Message: "Ivalid ID"}, http.StatusBadRequest)
			return
		}
		record, err := ctx.ctrl.GetRecord(idInt)
		if err != nil {
			logger.Warnf("failed to get records from controller: %s", err.Error())
			writeJSONResponse(w, logger, apischema.Error{Message: "Failed to receive data from controller"}, http.StatusInternalServerError)
			return
		}
		recordsAPI := schemabuilder.BuildRecordAPIFrom(record)
		writeJSONResponse(w, logger, recordsAPI, http.StatusOK)
	}
}

func NewCreateRecordsHandler(ctx *APIContext) httprouter.Handle {
	logger := ctx.logger.WithFields(log.Fields{"handler": "createRecords"})
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// TODO: check content type
		var recordAPI *apischema.Record
		err := readJSON(r.Body, recordAPI)
		defer r.Body.Close()
		if err != nil {
			logger.Warnf("failed to read JSON: %s", err.Error())
			writeJSONResponse(w, logger, apischema.Error{Message: "Ivalid JSON"}, http.StatusBadRequest)
			return
		}
		record := schemabuilder.BuildRecordFrom(recordAPI)
		err = ctx.ctrl.CreateRecord(record)
		if err != nil {
			logger.Warnf("failed to get records from controller: %s", err.Error())
			writeJSONResponse(w, logger, apischema.Error{Message: "Ivalid JSON"}, http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func readJSON(requestBody io.Reader, out any) error {
	// TODO: prevent overflow (read by batches or set max size)
	recordsJSON, err := io.ReadAll(requestBody)
	if err != nil {
		return err
	}
	err = json.Unmarshal(recordsJSON, &out)
	if err != nil {
		return err
	}
	return err
}

func writeJSONResponse(w http.ResponseWriter, logger log.Logger, body any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		logger.Warnf("failed to write JSON response: %s", err.Error())
	}
	w.WriteHeader(statusCode)
	// TODO: do not log private info
	logger.Debugf("response written: %+v", body)
}
