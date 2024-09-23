package handlertype

import "net/http"

// IPingHandler interface for our PingHandler
type IPingHandler interface {
	Ping(w http.ResponseWriter, r *http.Request)
	InternalServerError(w http.ResponseWriter, r *http.Request)
	NotFoundError(w http.ResponseWriter, r *http.Request)
}

// ApiHandlers holds all API Handlers in one place
type ApiHandlers struct {
	Ping IPingHandler
}
