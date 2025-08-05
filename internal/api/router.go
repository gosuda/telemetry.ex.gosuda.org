package api

import (
	"telemetry.ex.gosuda.org/telemetry/internal/server"
	"telemetry.ex.gosuda.org/telemetry/internal/types"
)

// RegisterRoutes registers all API routes with the server and returns the server
func RegisterRoutes(s *server.Server, is types.InternalServiceProvider) *server.Server {
	s.RegisterHandler("GET", "/healthz", HealthzHandler(is))
	return s
}
