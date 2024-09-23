package handler

import (
	"net/http"
	"strings"
	"testing"

	"github.com/bluznierca1/go-chi-swagger-testing/internal/utils/testutils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
)

var (
	openApiDoc *openapi3.T
)

func TestMain(m *testing.M) {
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
