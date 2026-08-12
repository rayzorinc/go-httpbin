package httpbin

import (
	"bytes"
	"net/http"
	"strings"
)

const openAPIContentType = "application/vnd.oai.openapi+json;version=3.2"

var openAPIAnyMethods = []string{
	"delete",
	"get",
	"head",
	"options",
	"patch",
	"post",
	"put",
	"query",
	"trace",
}

type openAPIMux struct {
	mux      *http.ServeMux
	patterns []string
}

func newOpenAPIMux() *openAPIMux {
	return &openAPIMux{mux: http.NewServeMux()}
}

func (m *openAPIMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func (m *openAPIMux) HandleFunc(pattern string, handler http.HandlerFunc) {
	// Keep route registration as the source of truth for the generated schema.
	m.patterns = append(m.patterns, pattern)
	m.handleFunc(pattern, handler)
}

func (m *openAPIMux) handleFunc(pattern string, handler http.HandlerFunc) {
	// Documentation routes use this path so they are served by the same mux
	// without appearing as go-httpbin API operations in the schema.
	m.mux.HandleFunc(pattern, handler)
}

func (m *openAPIMux) registerOpenAPI(prefix string) {
	const pattern = "GET /openapi.json"
	var body []byte
	serveOpenAPI := func(w http.ResponseWriter, _ *http.Request) {
		writeResponse(w, http.StatusOK, openAPIContentType, body)
	}
	// Include the schema endpoint in its own document, then freeze the document
	// before adding the Swagger UI routes that consume it.
	m.HandleFunc(pattern, serveOpenAPI)
	body = buildOpenAPIDocument(prefix, m.patterns)
	m.patterns = nil
	m.registerSwaggerUI(prefix, serveOpenAPI)
}

func buildOpenAPIDocument(prefix string, patterns []string) []byte {
	serverURL := prefix
	if serverURL == "" {
		serverURL = "/"
	}

	paths := map[string]map[string]any{}
	operation := map[string]any{
		"responses": map[string]any{
			"default": map[string]string{
				"description": "Response from go-httpbin.",
			},
		},
	}
	for _, pattern := range patterns {
		method, path := parseServeMuxPattern(pattern)
		pathItem, ok := paths[path]
		if !ok {
			pathItem = map[string]any{}
			if parameters := openAPIPathParameters(path); len(parameters) > 0 {
				pathItem["parameters"] = parameters
			}
			paths[path] = pathItem
		}

		for _, method := range openAPIMethods(method) {
			pathItem[method] = operation
		}
	}

	document := map[string]any{
		"openapi": "3.2.0",
		"info": map[string]string{
			"title":       "go-httpbin",
			"version":     "1.0.0",
			"description": "An HTTP request and response testing service.",
		},
		"servers": []map[string]string{{"url": serverURL}},
		"paths":   paths,
	}

	var buf bytes.Buffer
	mustMarshalJSON(&buf, document)
	return buf.Bytes()
}

func openAPIMethods(method string) []string {
	// The preflight and autohead middleware make OPTIONS available for every
	// endpoint and HEAD available wherever GET is registered.
	switch method {
	case "":
		return openAPIAnyMethods
	case "get":
		return []string{"get", "head", "options"}
	case "options":
		return []string{"options"}
	default:
		return []string{method, "options"}
	}
}

func parseServeMuxPattern(pattern string) (method, path string) {
	method, path, ok := strings.Cut(pattern, " ")
	if !ok {
		path = method
		method = ""
	}
	// OpenAPI has no equivalents for ServeMux's end marker or variadic wildcard,
	// so express both using the corresponding OpenAPI path template.
	path = strings.TrimSuffix(path, "{$}")
	path = strings.ReplaceAll(path, "...}", "}")
	return strings.ToLower(method), path
}

func openAPIPathParameters(path string) []map[string]any {
	var parameters []map[string]any
	for rest := path; ; {
		start := strings.IndexByte(rest, '{')
		if start < 0 {
			break
		}
		rest = rest[start+1:]
		end := strings.IndexByte(rest, '}')
		if end < 0 {
			break
		}
		name := strings.TrimSuffix(rest[:end], "...")
		parameters = append(parameters, map[string]any{
			"name":     name,
			"in":       "path",
			"required": true,
			"schema":   map[string]string{"type": "string"},
		})
		rest = rest[end+1:]
	}
	return parameters
}
