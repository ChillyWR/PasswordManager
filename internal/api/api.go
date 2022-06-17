package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
	"github.com/invopop/yaml"
	"github.com/julienschmidt/httprouter"

	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/schema/apischema"
	"github.com/okutsen/PasswordManager/schema/controllerSchema"
)

const (
	// PPN: Path Parameter Name
	// HPN: Header Parameter Name
	IDPPN                 = "id"
	CorrelationIDHPN      = "X-Request-ID"
	AuthorizationTokenHPN = "Authorization"

	RequestContextName = "rctx"
)

type Controller interface {
	AllRecords() ([]controllerSchema.Record, error)
	Record(id uuid.UUID) (*controllerSchema.Record, error)
	CreateRecord(record *controllerSchema.Record) (*controllerSchema.Record, error)
	UpdateRecord(id uuid.UUID, record *controllerSchema.Record) (*controllerSchema.Record, error)
	DeleteRecord(id uuid.UUID) (*controllerSchema.Record, error)

	AllUsers() ([]controllerSchema.User, error)
	User(id uuid.UUID) (*controllerSchema.User, error)
	CreateUser(user *controllerSchema.User) (*controllerSchema.User, error)
	UpdateUser(id uuid.UUID, user *controllerSchema.User) (*controllerSchema.User, error)
	DeleteUser(id uuid.UUID) (*controllerSchema.User, error)
}

type RequestContext struct {
	corID uuid.UUID
	ps    httprouter.Params
}

type API struct {
	config *Config
	ctx    *APIContext
	server http.Server
}

type APIContext struct {
	ctrl   Controller
	logger log.Logger
}

type HandlerFunc func(rw http.ResponseWriter, r *http.Request, ctx *RequestContext)

func New(config *Config, controller Controller, logger log.Logger) *API {
	return &API{
		config: config,
		ctx: &APIContext{
			ctrl:   controller,
			logger: logger.WithFields(log.Fields{"service": "API"}),
		},
	}
}

func (api *API) Start() error {
	api.ctx.logger.Info("API started")
	router := httprouter.New()
	api.SetFunctionalEndpoints(router)
	api.SetRecordEndpoints(router)
	api.SetUserEndpoints(router)

	api.server = http.Server{Addr: api.config.Address(), Handler: router}

	return api.server.ListenAndServe()
}

func (api *API) SetRecordEndpoints(r *httprouter.Router) {
	r.GET("/records",
		AuthorizationCheck(api.ctx.logger, ContextSetter(api.ctx.logger,
			NewListRecordsHandler(api.ctx))))
	r.POST("/records",
		AuthorizationCheck(api.ctx.logger, ContextSetter(api.ctx.logger,
			NewCreateRecordHandler(api.ctx))))
	r.GET(fmt.Sprintf("/records/:%s", IDPPN),
		AuthorizationCheck(api.ctx.logger, ContextSetter(api.ctx.logger,
			NewGetRecordHandler(api.ctx))))
	r.PUT(fmt.Sprintf("/records/:%s", IDPPN),
		AuthorizationCheck(api.ctx.logger, ContextSetter(api.ctx.logger,
			NewUpdateRecordHandler(api.ctx))))
	r.DELETE(fmt.Sprintf("/records/:%s", IDPPN),
		AuthorizationCheck(api.ctx.logger, ContextSetter(api.ctx.logger,
			NewDeleteRecordHandler(api.ctx))))
}

func (api *API) SetUserEndpoints(r *httprouter.Router) {
	r.GET("/users",
		AuthorizationCheck(api.ctx.logger, ContextSetter(api.ctx.logger,
			NewListUsersHandler(api.ctx))))
	r.POST("/users",
		AuthorizationCheck(api.ctx.logger, ContextSetter(api.ctx.logger,
			NewCreateUserHandler(api.ctx))))
	r.GET(fmt.Sprintf("/users/:%s", IDPPN),
		AuthorizationCheck(api.ctx.logger, ContextSetter(api.ctx.logger,
			NewGetUserHandler(api.ctx))))
	r.PUT(fmt.Sprintf("/users/:%s", IDPPN),
		AuthorizationCheck(api.ctx.logger, ContextSetter(api.ctx.logger,
			NewUpdateUserHandler(api.ctx))))
	r.DELETE(fmt.Sprintf("/users/:%s", IDPPN),
		AuthorizationCheck(api.ctx.logger, ContextSetter(api.ctx.logger,
			NewDeleteUserHandler(api.ctx))))
}

func (api *API) SetFunctionalEndpoints(r *httprouter.Router) {
	spec := NewOpenAPIv3(api.config, api.ctx.logger)
	r.GET("/authMePlease",
		ContextSetter(api.ctx.logger,
			NewFreeAccessHandler(api.ctx.logger)))
	r.GET("/openapi3.json",
		ContextSetter(api.ctx.logger,
			NewJSONSpecHandler(api.ctx.logger, spec)))
	r.GET("/openapi3.yaml",
		ContextSetter(api.ctx.logger,
			NewYAMLSpecHandler(api.ctx.logger, spec)))
}

func (api *API) Stop(ctx context.Context) error {
	api.ctx.logger.Infof("shutting down server")
	return api.server.Shutdown(ctx)
}

func NewJSONSpecHandler(parentLogger log.Logger, spec *openapi3.T) http.HandlerFunc {
	logger := parentLogger.WithFields(log.Fields{"handler": "SpecHandler"})
	return func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, &spec, http.StatusOK, logger)
	}
}

func NewYAMLSpecHandler(parentLogger log.Logger, spec *openapi3.T) http.HandlerFunc {
	logger := parentLogger.WithFields(log.Fields{"handler": "SpecHandler"})
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
		w.WriteHeader(http.StatusOK)
	}
}

func NewFreeAccessHandler(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := GenerateJWT(uuid.NewString())
		if err != nil {
			writeResponse(w, apischema.Error{Message: "Oops, failed to generate your token"}, http.StatusOK, logger)
		}
		t := struct {
			Message string `json:"message,omitempty"`
			Token   string `json:"token"`
		}{
			Message: "Here, use this as Authorization header",
			Token:   token,
		}
		writeResponse(w, &t, http.StatusOK, logger)
	}
}
