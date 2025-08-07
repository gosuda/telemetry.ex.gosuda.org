package types

import (
	"context"
)

type PersistenceService interface {
	Ping(ctx context.Context) error

	RandflakeGC(ctx context.Context) error
	RandflakeLeaseCreate(ctx context.Context) (*RandflakeLease, error)
	RandflakeLeaseExtend(ctx context.Context, prev *RandflakeLease) (*RandflakeLease, error)

	ClientRegisterFingerprint(ctx context.Context, fpID int64, clientID int64, userAgent string, userAgentData string, fpversion int32, fphash string) error
	ClientLookupByID(ctx context.Context, clientID int64) (ClientIdentifier, error)
	ClientLookupByToken(ctx context.Context, token string) (ClientIdentifier, error)
	ClientVerifyToken(ctx context.Context, clientID int64, token string) (bool, error)
}
