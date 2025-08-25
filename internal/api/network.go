package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	"telemetry.gosuda.org/telemetry/internal/types"
)

func GetzHandler(is types.InternalServiceProvider) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		type NetworkHandlerResponse struct {
			Args    url.Values        `json:"args"`
			Headers map[string]string `json:"headers"`
			Origin  string            `json:"origin"`
			URL     string            `json:"url"`
		}

		response := NetworkHandlerResponse{
			Args:    r.URL.Query(),
			Headers: make(map[string]string, len(r.Header)),
			Origin:  r.RemoteAddr,
			URL:     r.URL.String(),
		}
		for k := range r.Header {
			response.Headers[k] = r.Header.Get(k)
		}

		if os.Getenv("IP_HEADER") != "" {
			realip := strings.Split(r.Header.Get(os.Getenv("IP_HEADER")), ",")
			response.Origin = strings.TrimSpace(realip[0])
		}

		jsonResponse, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	}
}
