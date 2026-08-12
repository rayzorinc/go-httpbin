package httpbin

import "net/http"

const swaggerContentSecurityPolicy = "default-src 'self'; img-src 'self' data:; script-src 'self'; style-src 'self' 'unsafe-inline'; connect-src 'self'"

type swaggerAsset struct {
	path        string
	contentType string
	body        []byte
}

var (
	swaggerIndexHTML = mustStaticAsset("swagger/index.html")
	swaggerAssets    = []swaggerAsset{
		{"swagger-ui.css", "text/css; charset=utf-8", mustStaticAsset("swagger/swagger-ui.css")},
		{"swagger-ui-bundle.js", "text/javascript; charset=utf-8", mustStaticAsset("swagger/swagger-ui-bundle.js")},
		{"swagger-ui-standalone-preset.js", "text/javascript; charset=utf-8", mustStaticAsset("swagger/swagger-ui-standalone-preset.js")},
		{"swagger-initializer.js", "text/javascript; charset=utf-8", mustStaticAsset("swagger/swagger-initializer.js")},
		{"favicon-16x16.png", "image/png", mustStaticAsset("swagger/favicon-16x16.png")},
		{"favicon-32x32.png", "image/png", mustStaticAsset("swagger/favicon-32x32.png")},
	}
)

func (m *openAPIMux) registerSwaggerUI(prefix string, serveOpenAPI http.HandlerFunc) {
	serveIndex := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Security-Policy", swaggerContentSecurityPolicy)
		writeHTML(w, swaggerIndexHTML, http.StatusOK)
	}

	m.handleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, prefix+"/swagger/", http.StatusPermanentRedirect)
	})
	m.handleFunc("GET /swagger/{$}", serveIndex)
	m.handleFunc("GET /swagger/index.html", serveIndex)
	m.handleFunc("GET /swagger/doc.json", serveOpenAPI)
	for _, asset := range swaggerAssets {
		m.handleFunc("GET /swagger/"+asset.path, func(w http.ResponseWriter, _ *http.Request) {
			writeResponse(w, http.StatusOK, asset.contentType, asset.body)
		})
	}
}
