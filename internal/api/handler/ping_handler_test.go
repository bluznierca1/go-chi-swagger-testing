package handler

import (
	"context"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/bluznierca1/go-chi-swagger-testing/internal/utils/testutils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

var (
	openApiDoc *openapi3.T
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatalf("Could not initialize .env file %v", err)
	}

	openApiDoc = testutils.LoadOpenApiConfigFile()

	m.Run()
}

func TestPingHandler_Ping(t *testing.T) {
	testCases := []struct {
		name                  string
		body                  string
		expectedStatusCode    int
		expectedBodySubString string
	}{
		{
			name:                  "ping returns pong",
			body:                  "",
			expectedStatusCode:    http.StatusOK,
			expectedBodySubString: "pong",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpReq, err := http.NewRequest(http.MethodGet, "https://to-be-defined.whoknows/api/ping", strings.NewReader(tc.body))
			if err != nil {
				t.Fatalf("http.NewRequest(); err = %v", err)
			}

			// skip auth check on our schema as route not protected
			requestOptions := &openapi3filter.Options{
				AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
			}

			responseOptions := &openapi3filter.Options{}

			// Let's put our testutils openapi helper to use
			testData := testutils.OpenApiTestRequestData{
				HttpReq:                   httpReq,
				OpenApiDoc:                openApiDoc,
				RequestValidationOptions:  requestOptions,
				ResponseValidationOptions: responseOptions,
				HandlerCallback: func(w http.ResponseWriter, r *http.Request) {
					pingHandler := &PingHandler{}
					pingHandler.Ping(w, r)
				},
				ExpectedBodySubstring: tc.expectedBodySubString,
			}

			// Here we verify if response is what we expect
			err = testutils.OpenApiTestRequest(testData)
			if err != nil {
				t.Fatalf("testutils.OpenApiTestRequest(); err = %v", err)
			}

		})
	}
}

func TestPingHandler_GetRecord(t *testing.T) {
	urlBase := "https://to-be-defined.whoknows/api/get-record/"

	testCases := []struct {
		name                  string
		body                  string
		expectedStatusCode    int
		expectedBodySubString string
		id                    string
	}{
		{
			name:                  "invalid_id_param",
			body:                  "",
			expectedStatusCode:    http.StatusUnprocessableEntity,
			expectedBodySubString: "err_invalid_id",
			id:                    "-1",
		},
		{
			name:                  "not_found",
			body:                  "",
			expectedStatusCode:    http.StatusNotFound,
			expectedBodySubString: "err_not_found",
			id:                    "1",
		},
		{
			name:                  "record_found",
			body:                  "",
			expectedStatusCode:    http.StatusOK,
			expectedBodySubString: "5",
			id:                    "5",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			url := urlBase + tc.id

			httpReq, err := http.NewRequest(http.MethodGet, url, strings.NewReader(tc.body))
			if err != nil {
				t.Fatalf("http.NewRequest(); err = %v", err)
			}

			// Create a new chi.RouteContext and set the URL parameter
			// due to no context carried over in testing env
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)

			// Attach the chi.RouteContext to the request's context
			ctx := context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx)
			httpReq = httpReq.WithContext(ctx)

			// skip auth check on our schema as route not protected
			requestOptions := &openapi3filter.Options{
				AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
			}

			responseOptions := &openapi3filter.Options{}

			// Let's put our testutils openapi helper to use
			testData := testutils.OpenApiTestRequestData{
				HttpReq:                   httpReq,
				OpenApiDoc:                openApiDoc,
				RequestValidationOptions:  requestOptions,
				ResponseValidationOptions: responseOptions,
				HandlerCallback: func(w http.ResponseWriter, r *http.Request) {
					pingHandler := &PingHandler{}
					pingHandler.GetRecord(w, r)
				},
				ExpectedBodySubstring: tc.expectedBodySubString,
			}

			// Here we verify if response is what we expect
			err = testutils.OpenApiTestRequest(testData)
			if err != nil {
				t.Fatalf("testutils.OpenApiTestRequest(); err = %v", err)
			}
		})
	}
}
