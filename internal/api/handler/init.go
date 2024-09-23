package handler

import "github.com/bluznierca1/go-chi-swagger-testing/internal/api/handler/handlertype"

// InitApiHandlers inits all of our handlers in one place
func InitApiHandlers() *handlertype.ApiHandlers {
	return &handlertype.ApiHandlers{
		Ping: NewPingHandler(),
	}
}
