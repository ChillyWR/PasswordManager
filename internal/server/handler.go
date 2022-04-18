package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/okutsen/PasswordManager/internal/log"
)

// TODO: rename because router is Handler
type Handler struct {
	router *httprouter.Router
	log    log.Logger
}

func NewHandler() *Handler {
	h := &Handler{
		router: httprouter.New(),
	}
	h.setupRouter()
	return h
}

func (h *Handler) setupRouter() {
	h.router.GET("/records", h.getRecords)
	h.router.POST("/records", h.createRecords)
}

func (h *Handler) getRecords(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.log.Info("DomainServer: Endpoint Hit: getRecords")
}

func (h *Handler) createRecords(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.log.Info("DomainServer: Endpoint Hit: createRecords")
}
