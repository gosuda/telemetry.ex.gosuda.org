package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"telemetry.ex.gosuda.org/telemetry/internal/core"
	"telemetry.ex.gosuda.org/telemetry/internal/types"
)

// POST /counts/bulk
func BulkCountsHandler(is types.InternalServiceProvider) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		defer r.Body.Close()

		var req types.BulkCountsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg("failed to decode bulk counts request")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
			return
		}

		if len(req.Urls) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "urls list is required"})
			return
		}

		// Normalize and deduplicate URLs while preserving first-seen order
		normalized := make([]string, 0, len(req.Urls))
		seen := make(map[string]struct{})
		for _, u := range req.Urls {
			n, err := core.NormalizeURL(u)
			if err != nil {
				log.Debug().Str("url", u).Err(err).Msg("failed to normalize url")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid url in list"})
				return
			}
			if _, ok := seen[n]; !ok {
				seen[n] = struct{}{}
				normalized = append(normalized, n)
			}
		}

		rows, err := is.BulkCountsByUrls(r.Context(), normalized)
		if err != nil {
			log.Error().Err(err).Msg("failed to query bulk counts")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Build a map with default zero values so missing URLs return 0 counts
		resultsMap := make(map[string]types.BulkCountEntry, len(normalized))
		for _, u := range normalized {
			resultsMap[u] = types.BulkCountEntry{
				URL:       u,
				ViewCount: 0,
				LikeCount: 0,
			}
		}
		for _, rr := range rows {
			resultsMap[rr.URL] = rr
		}

		// Preserve normalized (first-seen) order in response
		resp := types.BulkCountsResponse{
			Results: make([]types.BulkCountEntry, 0, len(normalized)),
		}
		for _, u := range normalized {
			resp.Results = append(resp.Results, resultsMap[u])
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
