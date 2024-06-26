package api

import (
	"net/http"

	pmlogger "github.com/ChillyWR/PasswordManager/internal/logger"
	"github.com/ChillyWR/PasswordManager/model"
)

func NewLoginHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(pmlogger.Fields{
		"handler": "ListUsers",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(pmlogger.Fields{
			"cor_id": rctx.corID.String(),
		})

		var user model.UserForm
		if err := readBody(r.Body, &user); err != nil {
			logger.Errorf("Failed to read body: %s", err.Error())
			writeResponse(w, Error{Message: InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}

		id, err := apictx.ctrl.Login(&user)
		if err != nil {
			logger.Errorf("Failed to login: %s", err.Error())
			writeError(w, err, logger)
			return
		}

		token, err := GenerateJWT(id.String())
		if err != nil {
			logger.Errorf("Failed to generate jwt: %s", err.Error())
			writeResponse(w, Error{Message: "Oops, failed to generate your token"}, http.StatusInternalServerError, logger)
			return
		}

		t := struct {
			Message string `json:"message,omitempty"`
			Token   string `json:"token"`
		}{
			Message: "Welcome, welcome, use this as Authorization header",
			Token:   token,
		}

		writeResponse(w, t, http.StatusOK, logger)
	}
}

func NewListUsersHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(pmlogger.Fields{
		"handler": "ListUsers",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(pmlogger.Fields{
			"cor_id": rctx.corID.String(),
		})

		users, err := apictx.ctrl.AllUsers()
		if err != nil {
			logger.Errorf("Failed to list users: %s", err.Error())
			writeError(w, err, logger)
			return
		}

		writeResponse(w, users, http.StatusOK, logger)
	}
}

func NewGetUserHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(pmlogger.Fields{
		"handler": "GetUser",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(pmlogger.Fields{
			"cor_id": rctx.corID.String(),
		})

		userID, err := getIDFrom(rctx.params, logger)
		if err != nil {
			logger.Errorf("Invalid user id: %s", err.Error())
			writeResponse(w, Error{Message: InvalidUserIDMessage}, http.StatusBadRequest, logger)
			return
		}

		user, err := apictx.ctrl.GetUser(userID)
		if err != nil {
			logger.Errorf("Failed to get user: %s", err.Error())
			writeError(w, err, logger)
			return
		}

		writeResponse(w, user, http.StatusOK, logger)
	}
}

func NewCreateUserHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(pmlogger.Fields{
		"handler": "CreateUser",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(pmlogger.Fields{
			"cor_id": rctx.corID.String(),
		})

		var user model.UserForm
		if err := readBody(r.Body, &user); err != nil {
			logger.Errorf("Failed to read body: %s", err.Error())
			writeResponse(w, Error{Message: InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}

		result, err := apictx.ctrl.CreateUser(&user)
		if err != nil {
			logger.Errorf("Failed to create user: %s", err.Error())
			writeError(w, err, logger)
			return
		}

		token, err := GenerateJWT(result.ID.String())
		if err != nil {
			logger.Errorf("Failed to generate jwt: %s", err.Error())
			writeResponse(w, nil, http.StatusInternalServerError, logger)
			return
		}

		response := struct {
			User  *model.User
			Token string
		}{
			User:  result,
			Token: token,
		}

		writeResponse(w, response, http.StatusCreated, logger)
	}
}

func NewUpdateUserHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(pmlogger.Fields{
		"handler": "UpdateUser",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(pmlogger.Fields{
			"cor_id": rctx.corID.String(),
		})

		userID, err := getIDFrom(rctx.params, logger)
		if err != nil {
			logger.Errorf("Invalid User id: %s", err.Error())
			writeResponse(w, Error{Message: InvalidUserIDMessage}, http.StatusBadRequest, logger)
			return
		}

		var form model.UserForm
		if err = readBody(r.Body, &form); err != nil {
			logger.Errorf("Failed to read JSON: %s", err.Error())
			writeResponse(w, Error{Message: InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}

		result, err := apictx.ctrl.UpdateUser(userID, &form)
		if err != nil {
			logger.Errorf("Failed to update user: %s", err.Error())
			writeError(w, err, logger)
			return
		}

		writeResponse(w, result, http.StatusAccepted, logger)
	}
}

func NewDeleteUserHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(pmlogger.Fields{
		"handler": "DeleteUser",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(pmlogger.Fields{
			"cor_id": rctx.corID.String(),
		})

		userID, err := getIDFrom(rctx.params, logger)
		if err != nil {
			logger.Errorf("Invalid User id: %s", err.Error())
			writeResponse(w, Error{Message: InvalidUserIDMessage}, http.StatusBadRequest, logger)
			return
		}

		result, err := apictx.ctrl.DeleteUser(userID)
		if err != nil {
			logger.Errorf("Failed to delete user: %s", err.Error())
			writeError(w, err, logger)
			return
		}

		writeResponse(w, result, http.StatusOK, logger)
	}
}
