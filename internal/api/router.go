package api

import (
	"github.com/julienschmidt/httprouter"
	"telemetry.ex.gosuda.org/telemetry/internal/types"
)

// RegisterRoutes registers all API routes with the server and returns the server
func RegisterRoutes(s *httprouter.Router, is types.InternalServiceProvider) {
	s.Handle("GET", "/healthz", HealthzHandler(is))
}
