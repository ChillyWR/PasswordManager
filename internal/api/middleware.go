package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"github.com/okutsen/PasswordManager/internal/log"
)

// ContextSetter reads header, creates RequestContext and adds it to r.Context
func ContextSetter(logger log.Logger, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		corIDStr := r.Header.Get(CorrelationIDHPN)

		corID, err := uuid.Parse(corIDStr)
		if err != nil {
			logger.Debugf("Invalid corID <%s>: %s", corIDStr, err)
			corID = uuid.New()
			logger.Debugf("Setting new corID: %s", corID.String())
		}

		ctx := context.WithValue(r.Context(), RequestContextName, &RequestContext{
			corID:  corID,
			params: ps,
		})

		next(w, r.WithContext(ctx), ps)
	}
}

func AuthorizationCheck(logger log.Logger, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		tokenStr := r.Header.Get(AuthorizationTokenHPN)
		if tokenStr == "" {
			logger.Errorf("Failed to authorize: no token provided")
			writeResponse(w, Error{Message: UnAuthorizedMessage}, http.StatusUnauthorized, logger)
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("wrong signing method")
			}
			return SigningKey, nil
		})
		if err != nil {
			logger.Errorf("Failed to parse JWT token: %s", err.Error())
			return
		}

		if !token.Valid {
			logger.Warn("Received invalid JSW token")
			writeResponse(w, Error{Message: "Invalid token"}, http.StatusUnauthorized, logger)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Errorf("Failed to extract claims from JWT token")
			writeResponse(w, Error{Message: "Failed to extract claims"}, http.StatusInternalServerError, logger)
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			logger.Errorf("Failed to extract user ID from JWT token claims")
			writeResponse(w, Error{Message: "Failed to extract user ID from claims"}, http.StatusInternalServerError, logger)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			logger.Errorf("Failed to parse user ID <%s>: %s", userIDStr, err)
			writeResponse(w, Error{Message: "Failed to extract user ID from claims"}, http.StatusInternalServerError, logger)
			return
		}

		rctx := unpackRequestContext(r.Context(), logger)
		rctx.userID = userID
		ctx := context.WithValue(r.Context(), RequestContextName, rctx)

		next(w, r.WithContext(ctx), ps)
	}
}

func Dispatch(next http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		next(w, r)
	}
}
