package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/model/api"
	"github.com/okutsen/PasswordManager/model/builder"
)

func NewListUsersHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(log.Fields{"handler": "GetAllUsers"})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
			"cor_id": rctx.corID.String(),
		})
		users, err := apictx.ctrl.AllUsers()
		if err != nil {
			logger.Warnf("Failed to get users from controller: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		usersAPI := builder.BuildAPIUsersFromControllerUsers(users)
		// Write JSON by stream?
		writeResponse(w, usersAPI, http.StatusOK, logger)
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
			logger.Warnf("Invalid user id: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InvalidUserIDMessage}, http.StatusBadRequest, logger)
			return
		}
		user, err := apictx.ctrl.User(userID)
		if err != nil {
			logger.Warnf("Failed to get users from controller: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		writeResponse(w, builder.BuildAPIUserFromControllerUser(user), http.StatusOK, logger)
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
		var userAPI *api.User
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Warnf("Failed to read JSON: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}
		defer r.Body.Close()
		err = json.Unmarshal(body, &userAPI)
		if err != nil {
			logger.Warnf("failed to unmarshal JSON file: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InternalErrorMessage}, http.StatusBadRequest, logger)
			return
		}
		user := builder.BuildControllerUserFromAPIUser(userAPI)
		// TODO: if exists return err (409 Conflict)
		// FIXME: return created struct
		resultUser, err := apictx.ctrl.CreateUser(&user)
		if err != nil {
			logger.Warnf("Failed to get users from controller: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		writeResponse(w, builder.BuildAPIUserFromControllerUser(resultUser), http.StatusCreated, logger)
	}
}

func NewUpdateUserHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(log.Fields{
		"handler": "UpdateUsers",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
			"cor_id": rctx.corID.String(),
		})
		userID, err := getIDFrom(rctx.params, logger)
		if err != nil {
			logger.Warnf("Invalid User id: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InvalidUserIDMessage}, http.StatusBadRequest, logger)
			return
		}
		var userAPI *api.User
		err = readJSON(r.Body, userAPI)
		defer r.Body.Close()
		if err != nil {
			logger.Warnf("Failed to read JSON: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}
		if userID != userAPI.ID {
			logger.Warn("User id from path parameter doesn't match id from new user structure")
			writeResponse(w,
				api.Error{Message: api.InvalidJSONMessage}, http.StatusBadRequest, logger)
			return
		}
		user := builder.BuildControllerUserFromAPIUser(userAPI)
		resultUser, err := apictx.ctrl.UpdateUser(userID, &user)
		if err != nil {
			logger.Warnf("Failed to get users from controller: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		writeResponse(w,
			builder.BuildAPIUserFromControllerUser(resultUser), http.StatusAccepted, logger)
	}
}

func NewDeleteUserHandler(apictx *APIContext) http.HandlerFunc {
	logger := apictx.logger.WithFields(log.Fields{
		"handler": "DeleteUsers",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := unpackRequestContext(r.Context(), logger)
		logger = logger.WithFields(log.Fields{
			"cor_id": rctx.corID.String(),
		})
		userID, err := getIDFrom(rctx.params, logger)
		if err != nil {
			logger.Warnf("Invalid User id: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InvalidUserIDMessage}, http.StatusBadRequest, logger)
			return
		}
		resultUser, err := apictx.ctrl.DeleteUser(userID)
		if err != nil {
			logger.Errorf("Failed to get Users from controller: %s", err.Error())
			writeResponse(w,
				api.Error{Message: api.InternalErrorMessage}, http.StatusInternalServerError, logger)
			return
		}
		writeResponse(w,
			builder.BuildAPIUserFromControllerUser(resultUser), http.StatusOK, logger)
	}
}
