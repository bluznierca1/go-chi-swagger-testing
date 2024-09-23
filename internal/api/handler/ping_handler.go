package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bluznierca1/go-chi-swagger-testing/internal/api/handler/handlertype"
)

type PingHandler struct{}

func (p *PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("{\"ping\":\"pong\"}"))
}

// PingHandler_UnexpectedServerError this is just for testing purposes
func (p *PingHandler) InternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Content-Type", "text/plain")
}

func (p *PingHandler) NotFoundError(w http.ResponseWriter, r *http.Request) {
	notFoundResponseBody := struct {
		ErrorCode string `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	}{
		ErrorCode: "err_not_found",
		ErrorMsg:  "Entity not found.",
	}

	// let's ignore error for sake of simplicity of example...
	marshalledResponseBody, _ := json.Marshal(notFoundResponseBody)

	w.WriteHeader(http.StatusNotFound)
	w.Header().Add("Content-Type", "application/json")
	w.Write(marshalledResponseBody)
}

func NewPingHandler() handlertype.IPingHandler {
	return &PingHandler{}
}
