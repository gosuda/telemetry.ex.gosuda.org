package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"gosuda.org/randflake"
	"telemetry.ex.gosuda.org/telemetry/internal/types"
)

type ClientPassport struct {
	ClientID      string `json:"client_id"`
	ClientToken   string `json:"client_token"`
	ClientVersion string `json:"version"`
	FPVersion     int    `json:"fpv"`
	Fingerprint   string `json:"fp"`
	UserAgent     string `json:"ua"`
	UserAgentData string `json:"uad"`
}

func ClientCheckinHandler(is types.InternalServiceProvider) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		defer r.Body.Close()

		passport := ClientPassport{}
		err := json.NewDecoder(r.Body).Decode(&passport)
		if err != nil {
			log.Error().Err(err).Msg("Failed to decode client passport")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if passport.FPVersion != 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Debug().
			Str("client_id", passport.ClientID).
			Str("client_token", passport.ClientToken).
			Str("client_version", passport.ClientVersion).
			Int("fp_version", passport.FPVersion).
			Str("fingerprint", passport.Fingerprint).
			Str("user_agent", passport.UserAgent).
			Str("user_agent_data", passport.UserAgentData).
			Msg("Client Checkin Request Received")

		clientID, err := randflake.DecodeString(passport.ClientID)
		if err != nil {
			log.Debug().
				Str("client_id", passport.ClientID).
				Err(err).
				Msg("Failed to decode client ID")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ok, err := is.ClientVerifyToken(context.Background(), clientID, passport.ClientToken)
		if err != nil {
			log.Error().Err(err).Msg("Failed to verify client token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !ok {
			log.Debug().
				Str("client_id", passport.ClientID).
				Str("client_token", passport.ClientToken).
				Msg("Client token verification failed")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fpid, err := is.GenerateID()
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate fingerprint ID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = is.ClientRegisterFingerprint(
			context.Background(), fpid, clientID, passport.UserAgent, passport.UserAgentData, int32(passport.FPVersion), passport.Fingerprint)
		if err != nil {
			log.Error().Err(err).Msg("Failed to register fingerprint")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"status\": \"ok\"}"))
	}
}
