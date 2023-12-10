package api

import (
	"net/http"

	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/model"
)

func NewListUsersHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(log.Fields{
		"handler": "ListUsers",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
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
	logger := apictx.logger.WithFields(log.Fields{
		"handler": "GetUser",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
			"cor_id": rctx.corID.String(),
		})

		userID, err := getIDFrom(rctx.params, logger)
		if err != nil {
			logger.Errorf("Invalid user id: %s", err.Error())
			writeResponse(w, Error{Message: InvalidUserIDMessage}, http.StatusBadRequest, logger)
			return
		}

		user, err := apictx.ctrl.User(userID)
		if err != nil {
			logger.Errorf("Failed to get user: %s", err.Error())
			writeError(w, err, logger)
			return
		}

		writeResponse(w, user, http.StatusOK, logger)
	}
}

func NewCreateUserHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(log.Fields{
		"handler": "CreateUser",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
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

		response := struct{
			User *model.User
			Token string
		} {
			User: result,
			Token: token,
		}
		
		writeResponse(w, response, http.StatusCreated, logger)
	}
}

func NewUpdateUserHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(log.Fields{
		"handler": "UpdateUser",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
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
	logger := apictx.logger.WithFields(log.Fields{
		"handler": "DeleteUser",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
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
