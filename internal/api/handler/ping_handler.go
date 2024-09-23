package handler

import (
	"net/http"

	"github.com/bluznierca1/go-chi-swagger-testing/internal/api/handler/handlertype"
)

type PingHandler struct{}

func (p *PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("{\"ping\":\"pong\"}"))
}

func NewPingHandler() handlertype.IPingHandler {
	return &PingHandler{}
}
