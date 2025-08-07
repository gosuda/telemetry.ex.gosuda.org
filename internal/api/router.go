package api

import (
	"github.com/julienschmidt/httprouter"
	"telemetry.ex.gosuda.org/telemetry/internal/types"
)

// RegisterRoutes registers all API routes with the server and returns the server
func RegisterRoutes(s *httprouter.Router, is types.InternalServiceProvider) {
	// index
	s.Handle("GET", "/", IndexHandler(is))

	// z-routes
	s.Handle("GET", "/healthz", HealthzHandler(is))
	s.Handle("GET", "/idz", IDZHandler(is))

	// telemetry routes
	s.Handle("POST", "/client/status", ClientStatusHandler(is))
	s.Handle("POST", "/client/register", ClientRegisterHandler(is))
	s.Handle("POST", "/client/checkin", ClientCheckinHandler(is))
}
