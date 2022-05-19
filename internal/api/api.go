package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/schema/dbschema"
)

const (
	IDParamName = "id"
)

type Controller interface {
	AllRecords() ([]dbschema.Record, error)
	GetRecord(uint64) (dbschema.Record, error)
	CreateRecord(dbschema.Record) error
	UpdateRecord(dbschema.Record) error
	DeleteRecord(uint64) error
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

func New(config *Config, ctrl Controller, logger log.Logger) *API {
	return &API{
		config: config,
		ctx: &APIContext{
			ctrl:   ctrl,
			logger: logger.WithFields(log.Fields{"service": "API"}),
		},
	}
}

func (api *API) Start() error {
	api.ctx.logger.Info("API started")
	router := httprouter.New()

	router.GET("/records", NewEndpointLoggerMiddleware(api.ctx, NewAllRecordsHandler(api.ctx)))
	router.GET(fmt.Sprintf("/records/:%s", IDParamName), NewEndpointLoggerMiddleware(api.ctx, NewRecordHandler(api.ctx)))
	router.POST("/records", NewEndpointLoggerMiddleware(api.ctx, NewCreateRecordHandler(api.ctx)))
	router.PUT("/records", NewEndpointLoggerMiddleware(api.ctx, NewUpdateRecordHandler(api.ctx)))
	router.DELETE(fmt.Sprintf("/records/:%s", IDParamName), NewEndpointLoggerMiddleware(api.ctx, NewDeleteRecordHandler(api.ctx)))

	api.server = http.Server{Addr: api.config.Address(), Handler: router}

	return api.server.ListenAndServe()
}

func (api *API) Stop(ctx context.Context) error {
	api.ctx.logger.Infof("shutting down server")
	return api.server.Shutdown(ctx)
}
