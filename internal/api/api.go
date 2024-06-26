package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
	"github.com/invopop/yaml"
	"github.com/julienschmidt/httprouter"

	pmlogger "github.com/ChillyWR/PasswordManager/internal/logger"
	"github.com/ChillyWR/PasswordManager/model"
)

type Controller interface {
	AllRecords(userID uuid.UUID) ([]model.CredentialRecord, []model.LoginRecord, []model.CardRecord, []model.IdentityRecord, error)
	GetRecord(id uuid.UUID, userID uuid.UUID) (interface{}, error)
	CreateRecord(recordType model.RecordType, record json.RawMessage, userID uuid.UUID) (interface{}, error)
	UpdateRecord(id uuid.UUID, rawForm json.RawMessage, userID uuid.UUID) (interface{}, error)
	DeleteRecord(id uuid.UUID, userID uuid.UUID) (*model.CredentialRecord, error)

	Login(form *model.UserForm) (uuid.UUID, error)
	AllUsers() ([]model.User, error)
	GetUser(id uuid.UUID) (*model.User, error)
	CreateUser(user *model.UserForm) (*model.User, error)
	UpdateUser(id uuid.UUID, form *model.UserForm) (*model.User, error)
	DeleteUser(id uuid.UUID) (*model.User, error)
}

type RequestContext struct {
	corID  uuid.UUID
	userID uuid.UUID
	params httprouter.Params
}

type API struct {
	config *Config
	ctx    *APIContext
	server http.Server
}

type APIContext struct {
	ctrl   Controller
	logger pmlogger.Logger
}

type HandlerFunc func(rw http.ResponseWriter, r *http.Request, ctx *RequestContext)

func New(config *Config, ctrl Controller, logger pmlogger.Logger) (*API, error) {
	if ctrl == nil {
		return nil, errors.New("ctrl is nil")
	}

	return &API{
		config: config,
		ctx: &APIContext{
			ctrl:   ctrl,
			logger: logger.WithFields(pmlogger.Fields{"module": "api"}),
		},
	}, nil
}

func (api *API) Start() error {
	router := httprouter.New()
	api.SetFunctionalEndpoints(router)
	api.SetUserEndpoints(router)
	api.SetRecordEndpoints(router)

	api.server = http.Server{Addr: api.config.Address(), Handler: router}

	return api.server.ListenAndServe()
}

func (api *API) SetUserEndpoints(r *httprouter.Router) {
	r.POST("/login",
		ContextSetter(api.ctx.logger,
			Dispatch(NewLoginHandler(api.ctx))))
	r.GET("/users",
		ContextSetter(api.ctx.logger, Authentication(api.ctx.logger,
			Dispatch(NewListUsersHandler(api.ctx)))))
	r.POST("/users",
		ContextSetter(api.ctx.logger,
			Dispatch(NewCreateUserHandler(api.ctx))))
	r.GET(fmt.Sprintf("/users/:%s", IDPPN),
		ContextSetter(api.ctx.logger, Authentication(api.ctx.logger,
			Dispatch(NewGetUserHandler(api.ctx)))))
	r.PUT(fmt.Sprintf("/users/:%s", IDPPN),
		ContextSetter(api.ctx.logger, Authentication(api.ctx.logger,
			Dispatch(NewUpdateUserHandler(api.ctx)))))
	r.DELETE(fmt.Sprintf("/users/:%s", IDPPN),
		ContextSetter(api.ctx.logger, Authentication(api.ctx.logger,
			Dispatch(NewDeleteUserHandler(api.ctx)))))
}

func (api *API) SetRecordEndpoints(r *httprouter.Router) {
	r.GET("/records",
		ContextSetter(api.ctx.logger, Authentication(api.ctx.logger,
			Dispatch(NewListRecordsHandler(api.ctx)))))
	r.POST("/records",
		ContextSetter(api.ctx.logger, Authentication(api.ctx.logger,
			Dispatch(NewCreateRecordHandler(api.ctx)))))
	r.GET(fmt.Sprintf("/records/:%s", IDPPN),
		ContextSetter(api.ctx.logger, Authentication(api.ctx.logger,
			Dispatch(NewGetRecordHandler(api.ctx)))))
	r.PATCH(fmt.Sprintf("/records/:%s", IDPPN),
		ContextSetter(api.ctx.logger, Authentication(api.ctx.logger,
			Dispatch(NewUpdateRecordHandler(api.ctx)))))
	r.DELETE(fmt.Sprintf("/records/:%s", IDPPN),
		ContextSetter(api.ctx.logger, Authentication(api.ctx.logger,
			Dispatch(NewDeleteRecordHandler(api.ctx)))))
}

func (api *API) SetFunctionalEndpoints(r *httprouter.Router) {
	spec := NewOpenAPIv3(api.config, api.ctx.logger)
	r.GET("/openapi3.json",
		ContextSetter(api.ctx.logger,
			Dispatch(NewJSONSpecHandler(api.ctx.logger, spec))))
	r.GET("/openapi3.yaml",
		ContextSetter(api.ctx.logger,
			Dispatch(NewYAMLSpecHandler(api.ctx.logger, spec))))
}

func (api *API) Stop(ctx context.Context) error {
	api.ctx.logger.Infof("shutting down server")
	return api.server.Shutdown(ctx)
}

func NewJSONSpecHandler(parentLogger pmlogger.Logger, spec *openapi3.T) http.HandlerFunc {
	logger := parentLogger.WithFields(pmlogger.Fields{"handler": "SpecHandler"})
	return func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, &spec, http.StatusOK, logger)
	}
}

func NewYAMLSpecHandler(parentLogger pmlogger.Logger, spec *openapi3.T) http.HandlerFunc {
	logger := parentLogger.WithFields(pmlogger.Fields{"handler": "SpecHandler"})
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		data, err := yaml.Marshal(&spec)
		if err != nil {
			logger.Errorf("Failed to marshal yaml: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(data)
		if err != nil {
			logger.Errorf("Failed to write response: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
