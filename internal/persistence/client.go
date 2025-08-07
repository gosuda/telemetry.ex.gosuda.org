package persistence

import (
	"context"
	"time"

	"telemetry.ex.gosuda.org/telemetry/internal/persistence/database"
	"telemetry.ex.gosuda.org/telemetry/internal/types"
)

func (g *PersistenceClient) ClientLookupByID(ctx context.Context, id int64) (types.ClientIdentifier, error) {
	return g.db.ClientLookupByID(ctx, id)
}

func (g *PersistenceClient) ClientLookupByToken(ctx context.Context, token string) (types.ClientIdentifier, error) {
	return g.db.ClientLookupByToken(ctx, token)
}

func (g *PersistenceClient) ClientVerifyToken(ctx context.Context, id int64, token string) (bool, error) {
	ret, err := g.db.ClientVerifyToken(ctx, database.ClientVerifyTokenParams{
		ID:    id,
		Token: token,
	})
	if err != nil {
		return false, err
	}
	return ret == 1, nil
}

func (g *PersistenceClient) ClientRegisterFingerprint(
	ctx context.Context,
	fpID int64,
	clientID int64,
	userAgent string,
	userAgentData string,
	fpversion int32,
	fphash string,
) error {
	return g.db.ClientRegisterFingerprint(ctx, database.ClientRegisterFingerprintParams{
		ID:            fpID,
		ClientID:      clientID,
		UserAgent:     userAgent,
		UserAgentData: userAgentData,
		Fpversion:     fpversion,
		Fphash:        fphash,
		CreatedAt:     time.Now().UnixNano(),
	})
}
