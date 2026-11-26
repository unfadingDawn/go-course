package http

import (
	_ "embed"
	"net/http"
)

//go:embed swagger/index.html
var swaggerHTML []byte

//go:embed swagger/openapi.json
var openAPI []byte

func swaggerIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(swaggerHTML)
}

func swaggerSpec(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(openAPI)
}
