package testutils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/go-chi/chi/v5"
)

// Details of the endpoint defined in OpenApi file
type OpenApiPathDetails struct {
	Method        string
	PathOperation *openapi3.Operation
}

// RouterPathDetails details of the path defined on Router
type RouterPathDetails struct {
	Method string
}

// OpenApiTestRequestData struct holding all required parameters for testing endpoints
// including request, special options (f.ex. for excluding some checks), options for validation response,
// and endpoint handler callback
type OpenApiTestRequestData struct {
	// HttpReq - your *http.Request. Build it in your test and deliver. You can attach context, too.
	HttpReq *http.Request

	// Document needs to be scanned before and delivered to keep flexibility
	OpenApiDoc *openapi3.T

	// Options for Request Validation for controlling of what gets tested
	RequestValidationOptions *openapi3filter.Options

	// Options for Response Validation for controlling of what gets tested
	ResponseValidationOptions *openapi3filter.Options

	// Your Handler code to mock endpoint's behaviour
	HandlerCallback func(w http.ResponseWriter, r *http.Request)

	// String that is expected to appear in response. Empty string skips check
	ExpectedBodySubstring string
}

// LoadOpenApiYmlFile load OpenApi yaml/json file
//
// # Make sure to have already executed godotenv.Load before using it
func LoadOpenApiConfigFile() *openapi3.T {
	loader := openapi3.NewLoader()

	// I recommend updating it to absolute path with ENV vars...
	doc, err := loader.LoadFromFile("internal/utils/misc/openapi.yml")
	if err != nil {
		log.Fatalf("Failed load swagger spec: %v", err)
	}

	return doc
}

// GetOpenApiPathsData - builds a map with all routes and their details
//
// # It is loaded from config file for OpenApi
//
// If (for method) `PathOperation == nil` that means there is no route defined
//
// # Params:
//
// - openApiDoc: *openapi3.T
//
// # Returns:
//
// - map[string][]OpenApiPathDetails
func GetOpenApiPathsData(openApiDoc *openapi3.T) map[string][]OpenApiPathDetails {
	openApiPaths := map[string][]OpenApiPathDetails{}

	// build paths data from OpenAPI
	for path, pathItem := range openApiDoc.Paths.Map() {
		openApiPaths[path] = []OpenApiPathDetails{
			{Method: http.MethodGet, PathOperation: pathItem.Get},
			{Method: http.MethodPut, PathOperation: pathItem.Put},
			{Method: http.MethodDelete, PathOperation: pathItem.Delete},
			{Method: http.MethodPost, PathOperation: pathItem.Post},
			{Method: http.MethodPatch, PathOperation: pathItem.Patch},
		}
	}

	return openApiPaths
}

// GetRouterPathsData - builds a map with all routes and their methods
//
// # Params:
//
// - router: *chi.Mux
//
// # Returns:
//
// - map[string][]RouterPathDetails
//
// - error
func GetRouterPathsData(router *chi.Mux) (map[string][]RouterPathDetails, error) {
	routerPaths := map[string][]RouterPathDetails{}
	err := chi.Walk(router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		routerPaths[route] = append(routerPaths[route], RouterPathDetails{Method: method})
		return nil
	})

	return routerPaths, err
}

func NewOpenApiRouter(openApiDoc *openapi3.T) (routers.Router, error) {
	return gorillamux.NewRouter(openApiDoc)
}

// OpenApiTestRequest - perform full test of request against OpenApi Document:
//
// 1. Test request body before sending it to endpoint (possible control via testData.testData.RequestValidationOptions)
//
// 2. Mock server and serve endpoint with provided handler
//
// 3. Test received response
//
// # Options
//
// Both request and response objects take Options as parameter, so you can control behaviour of this test.
//
// # Params
//
// testData: OpenApiTestRequestData
//
// # Returns
//
// - error: in case something fails. You can use it to fail test
func OpenApiTestRequest(testData OpenApiTestRequestData) error {
	// extract context from provided httpReq and attach timeout
	reqContext, cancel := context.WithTimeout(testData.HttpReq.Context(), 10*time.Second)
	defer cancel()

	// get new instance of router for a document
	router, err := NewOpenApiRouter(testData.OpenApiDoc)
	if err != nil {
		return fmt.Errorf("NewOpenApiRouter(); err = %v", err)
	}

	// find route in schema or fail
	route, pathParams, err := router.FindRoute(testData.HttpReq)
	if err != nil {
		return fmt.Errorf("router.FindRoute(); err = %v", err)
	}

	// build request validation input for validation against schema
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    testData.HttpReq,
		PathParams: pathParams,
		Route:      route,
		Options:    testData.RequestValidationOptions,
	}

	// retrieve context from provided request object
	// and validate our request object
	// if you need to skip some pieces, add it into Options
	err = openapi3filter.ValidateRequest(reqContext, requestValidationInput)
	if err != nil {
		return fmt.Errorf("openapi3filter.ValidateRequest(); err = %v", err)
	}

	// create response recorded to cast response from mocked request
	responseRecorder := httptest.NewRecorder()

	// assign handler func
	handler := http.HandlerFunc(testData.HandlerCallback)

	// serve mocked endpoint for given request
	handler.ServeHTTP(responseRecorder, testData.HttpReq)

	// begin::validate response against OpenApi Schema

	// make sure to validate against all returned status codes
	testData.ResponseValidationOptions.IncludeResponseStatus = true
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 responseRecorder.Code,
		Header:                 responseRecorder.Header(),
		Options:                testData.ResponseValidationOptions,
	}
	responseValidationInput.SetBodyBytes(responseRecorder.Body.Bytes())
	err = openapi3filter.ValidateResponse(reqContext, responseValidationInput)
	if err != nil {
		return fmt.Errorf("openapi3filter.ValidateResponse(); err = %v", err)
	}

	// validate expected substring is in response body
	if testData.ExpectedBodySubstring != "" {
		if !strings.Contains(string(responseRecorder.Body.Bytes()), testData.ExpectedBodySubstring) {
			return fmt.Errorf("strings.Contains testData.ExpectedBodySubstring? FALSE")
		}
	}

	// end::validate response against OpenApi Schema

	return nil
}
