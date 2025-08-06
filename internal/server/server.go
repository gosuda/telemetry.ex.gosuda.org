package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"gosuda.org/randflake"
	"telemetry.ex.gosuda.org/telemetry/internal/api"
	"telemetry.ex.gosuda.org/telemetry/internal/types"
)

const (
	_RANDFLAKE_RENEW_WINDOW = int64(time.Minute * 8)
	_RANDFLAKE_SAFE_WINDOW  = int64(time.Second * 30)
)

type Server struct {
	mux *httprouter.Router

	ps     types.PersistenceService
	stopCh chan struct{}

	lease     *types.RandflakeLease
	randflake *randflake.Generator
}

var _ types.InternalServiceProvider = (*serverServiceProvider)(nil)

type serverServiceProvider struct {
	types.PersistenceService
	s *Server
}

func (g *serverServiceProvider) GenerateID() (int64, error) {
	return g.s.randflake.Generate()
}

func (g *serverServiceProvider) GenerateIDString() (string, error) {
	return g.s.randflake.GenerateString()
}

type ServerConfig struct {
	PersistenceService types.PersistenceService
	RandflakeSecret    []byte
}

// NewServer creates a new server instance
func NewServer(c *ServerConfig) (*Server, error) {
	g := &Server{
		ps:     c.PersistenceService,
		mux:    httprouter.New(),
		stopCh: make(chan struct{}),
	}

	ctx := context.Background()

	log.Debug().Msg("pinging persistence service")
	err := g.ps.Ping(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to ping persistence service")
		return nil, err
	}
	log.Debug().Msg("persistence service ping successful")

	log.Debug().Msg("creating initial randflake lease")
	lease, err := g.ps.RandflakeLeaseCreate(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to create randflake lease")
		return nil, err
	}
	g.lease = lease
	log.Debug().Int64("expires_at", g.lease.ExpiresAt).Int64("nodeid", g.lease.NodeID).Msg("randflake lease created")

	rf, err := randflake.NewGenerator(
		lease.NodeID,
		lease.CreatedAt/int64(time.Second),
		(lease.ExpiresAt-_RANDFLAKE_SAFE_WINDOW)/int64(time.Second),
		c.RandflakeSecret,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create randflake generator")
		return nil, err
	}
	g.randflake = rf

	// Start the randflake worker
	go g.randflakeWorker()

	is := &serverServiceProvider{
		PersistenceService: g.ps,
		s:                  g,
	}

	api.RegisterRoutes(g.mux, is)

	return g, nil
}

func (g *Server) randflakeWorker() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixNano()

			if now > g.lease.ExpiresAt {
				log.Error().Msg("randflake lease expired")
				lease, err := g.ps.RandflakeLeaseCreate(context.Background())
				if err != nil {
					log.Error().Err(err).Msg("failed to create randflake lease")
					continue
				}
				g.lease = lease
				g.randflake.UpdateLease(lease.CreatedAt, lease.ExpiresAt)
				log.Debug().Int64("expires_at", g.lease.ExpiresAt).Int64("nodeid", g.lease.NodeID).Msg("randflake lease created")
			}

			if now > g.lease.ExpiresAt-_RANDFLAKE_RENEW_WINDOW {
				log.Debug().Int64("now", now).Int64("expires_at", g.lease.ExpiresAt).Msg("randflake lease about to expire, trying to extend")
				func() {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
					defer cancel()
					lease, err := g.ps.RandflakeLeaseExtend(ctx, g.lease)
					if err != nil {
						log.Error().Err(err).Msg("failed to extend randflake lease")
						return
					}
					g.lease = lease
					g.randflake.UpdateLease(lease.CreatedAt, lease.ExpiresAt)
					log.Debug().Int64("expires_at", g.lease.ExpiresAt).Msg("randflake lease extended")
				}()
			}

			func() {
				gcStart := time.Now()
				log.Debug().Msg("running randflake gc")
				err := g.ps.RandflakeGC(context.Background())
				if err != nil {
					log.Error().Err(err).Msg("failed to run randflake gc")
				}
				log.Debug().Dur("duration", time.Since(gcStart)).Msg("randflake gc completed")
			}()
		case <-g.stopCh:
			return
		}
	}
}

func (g *Server) Serve(ln net.Listener) error {
	return http.Serve(ln, g.mux)
}

func (g *Server) Shutdown() {
	// Shutdown logic can be implemented here
}
