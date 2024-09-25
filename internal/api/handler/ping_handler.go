package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bluznierca1/go-chi-swagger-testing/internal/api/handler/handlertype"
	"github.com/go-chi/chi/v5"
)

type PingHandler struct{}

func (p *PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"ping":"pong"}`))
}

func (p *PingHandler) GetRecord(w http.ResponseWriter, r *http.Request) {
	var marshalledResponseBody []byte

	w.Header().Add("Content-Type", "application/json")

	// let's validate our ID (recommended to move it into some validator helper or middleware)
	invalidIdResponseBody := struct {
		ErrorCode string `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	}{
		ErrorCode: "err_invalid_id",
		ErrorMsg:  "Id must be integer greater than 0.",
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)
	if err != nil || id < 1 {
		// let's ignore error for sake of simplicity of example...
		marshalledResponseBody, _ = json.Marshal(invalidIdResponseBody)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(marshalledResponseBody)
		return
	}

	// If provided id is 5, let's return success for our case
	if id == 5 {
		successResponseBody := struct {
			Id int `json:"id"`
		}{
			Id: 5,
		}
		marshalledResponseBody, _ = json.Marshal(successResponseBody)
		w.WriteHeader(http.StatusOK)
		w.Write(marshalledResponseBody)
		return
	}

	// if id != 5, return not found
	notFoundResponseBody := struct {
		ErrorCode string `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	}{
		ErrorCode: "err_not_found",
		ErrorMsg:  "Entity not found.",
	}

	// let's ignore error for sake of simplicity of example...
	marshalledResponseBody, _ = json.Marshal(notFoundResponseBody)

	w.WriteHeader(http.StatusNotFound)
	w.Header().Add("Content-Type", "application/json")
	w.Write(marshalledResponseBody)
}

func NewPingHandler() handlertype.IPingHandler {
	return &PingHandler{}
}
