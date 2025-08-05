package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"telemetry.ex.gosuda.org/telemetry/internal/types"
)

var _status_ok = []byte(`{"status":"ok"}`)

// HealthzHandler returns a handler function that handles health check requests
func HealthzHandler(is types.InternalServiceProvider) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(_status_ok)
	}
}
