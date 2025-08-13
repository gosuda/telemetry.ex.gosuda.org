package api

import (
	"github.com/julienschmidt/httprouter"
	"telemetry.gosuda.org/telemetry/internal/types"
)

// RegisterRoutes registers all API routes with the server and returns the server
func RegisterRoutes(s *httprouter.Router, is types.InternalServiceProvider) {
	// index
	s.Handle("GET", "/", IndexHandler(is))

	// go package
	s.Handle("GET", "/telemetry", GoPackageHandler(is))

	// z-routes
	s.Handle("GET", "/healthz", HealthzHandler(is))
	s.Handle("GET", "/idz", IDZHandler(is))

	// telemetry routes
	s.Handle("POST", "/client/status", ClientStatusHandler(is))
	s.Handle("POST", "/client/register", ClientRegisterHandler(is))
	s.Handle("POST", "/client/checkin", ClientCheckinHandler(is))
	s.Handle("POST", "/client/view", ClientViewHandler(is))
	s.Handle("POST", "/client/like", ClientLikeHandler(is))

	// bulk counts endpoint (POST body: JSON { "urls": ["https://...","..."] })
	s.Handle("POST", "/counts/bulk", BulkCountsHandler(is))

	// view & like count lookup routes
	s.Handle("GET", "/view/count", ViewCountHandler(is))
	s.Handle("GET", "/like/count", LikeCountHandler(is))

	// Use the CORS middleware for global OPTIONS handling. The middleware will
	// handle preflight requests even when next is nil.
	s.GlobalOPTIONS = CORS(nil)
}
