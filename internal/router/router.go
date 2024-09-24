package router

import (
	"github.com/bluznierca1/go-chi-swagger-testing/internal/api/handler"
	"github.com/go-chi/chi/v5"
)

func SetupRouter() *chi.Mux {
	router := chi.NewRouter()

	apiHandlers := handler.InitApiHandlers()

	// Let's group below routes under /api
	router.Route("/api", func(r chi.Router) {
		r.Get("/ping", apiHandlers.Ping.Ping)
		r.Get("/get-record/{id}", apiHandlers.Ping.GetRecord)
	})

	return router
}
