// {root}/openapi_test.go
package main

import (
	"fmt"
	"strings"
	"testing"

	apiRouter "github.com/bluznierca1/go-chi-swagger-testing/internal/router"
	"github.com/bluznierca1/go-chi-swagger-testing/internal/utils/testutils"
	"github.com/getkin/kin-openapi/openapi3"
)

var (
	openApiDoc *openapi3.T
)

// Load OpenAPI file via package once (it will crash if something goes wrong)
func TestMain(m *testing.M) {
	openApiDoc = testutils.LoadOpenApiConfigFile()
	m.Run()
}

// TestOpenApi_CheckRegisteredRoutes tests the alignment between
// OpenAPI documentation and Chi router.
func TestOpenApi_CheckRegisteredRoutes(t *testing.T) {
	// Extract paths and their methods from OpenAPI documentation
	// it will be obtained in our own format (wrapper)
	openApiPaths := testutils.GetOpenApiPathsData(openApiDoc)

	// Initialize Chi router and extract its paths with methods
	router := apiRouter.SetupRouter()
	routerPaths, err := testutils.GetRouterPathsData(router)
	if err != nil {
		t.Fatalf("testutils.GetRouterPathsData(); err = %v", err)
	}

	// Slice to collect discrepancies
	var discrepancies []string

	// Preprocess OpenAPI methods into a map
	openApiMethodsMap := preprocessOpenApiMethods(openApiPaths)

	// Preprocess Router methods into a map
	routerMethodsMap := preprocessRouterMethods(routerPaths)

	// begin::Check OpenAPI to Router
	for path, openApiPathDetails := range openApiPaths {
		_, exists := routerPaths[path]
		if !exists {
			discrepancies = append(discrepancies, fmt.Sprintf("Missing path [%s] in Router", path))
			continue
		}

		for _, openApiDetail := range openApiPathDetails {
			if openApiDetail.PathOperation == nil {
				continue
			}
			if !routerMethodsMap[path][openApiDetail.Method] {
				discrepancies = append(discrepancies, fmt.Sprintf("Missing method [%s] for path [%s] in Router", openApiDetail.Method, path))
			}
		}
	}
	// end::Check OpenAPI to Router

	// begin::Check Router to OpenAPI
	for path, routerPathDetails := range routerPaths {
		_, exists := openApiPaths[path]
		if !exists {
			discrepancies = append(discrepancies, fmt.Sprintf("Missing path [%s] in OpenAPI file.", path))
			continue
		}

		for _, routerDetail := range routerPathDetails {
			if !openApiMethodsMap[path][routerDetail.Method] {
				discrepancies = append(discrepancies, fmt.Sprintf("Missing method [%s] for path [%s] in OpenAPI file.", routerDetail.Method, path))
			}
		}
	}
	// end::Check Router to OpenAPI

	// Report all discrepancies
	if len(discrepancies) > 0 {
		t.Errorf("Route discrepancies found:\n%s", strings.Join(discrepancies, "\n"))
	}
}

// preprocessOpenApiMethods converts OpenAPI path details into a map for easier usage
func preprocessOpenApiMethods(paths map[string][]testutils.OpenApiPathDetails) map[string]map[string]bool {
	methodsMap := make(map[string]map[string]bool)
	for path, details := range paths {
		if _, exists := methodsMap[path]; !exists {
			methodsMap[path] = make(map[string]bool)
		}
		for _, detail := range details {
			if detail.PathOperation != nil {
				methodsMap[path][detail.Method] = true
			}
		}
	}
	return methodsMap
}

// preprocessRouterMethods converts Router path details into a map for easier usage
func preprocessRouterMethods(paths map[string][]testutils.RouterPathDetails) map[string]map[string]bool {
	methodsMap := make(map[string]map[string]bool)
	for path, details := range paths {
		if _, exists := methodsMap[path]; !exists {
			methodsMap[path] = make(map[string]bool)
		}
		for _, detail := range details {
			methodsMap[path][detail.Method] = true
		}
	}
	return methodsMap
}
