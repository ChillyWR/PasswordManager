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

func AuthorizationCheck(log log.Logger, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		tokenStr := r.Header.Get(AuthorizationTokenHPN)
		if tokenStr == "" {
			writeResponse(w, Error{Message: UnAuthorizedMessage}, http.StatusUnauthorized, log)
			return
		}
		// TODO: use Bearer format
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("wrong signing method")
			}
			return SigningKey, nil
		})
		if err != nil {
			log.Warnf("Failed to parse JWT token: %s", err.Error())
		}

		if !token.Valid {
			log.Warn("Received invalid JSW token")
			writeResponse(w, Error{Message: "Invalid token"}, http.StatusUnauthorized, log)
			return
		}
		next(w, r, ps)
	}
}

// ContextSetter reads header, creates RequestContext and adds it to r.Context
func ContextSetter(logger log.Logger, next http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		corIDStr := r.Header.Get(CorrelationIDHPN)
		corID, err := uuid.Parse(corIDStr)
		if err != nil {
			logger.Warnf("Invalid corID <%s>: %s", corIDStr, err)
			corID = uuid.New()
			logger.Debugf("Setting new corID: %s", corID.String())
		}
		ctx := context.WithValue(r.Context(), RequestContextName, &RequestContext{
			corID:  corID,
			params: ps,
		})
		r = r.WithContext(ctx)
		next(w, r)
	}
}
