package types

import (
	"context"
)

type PersistenceService interface {
	Ping(ctx context.Context) error

	RandflakeGC(ctx context.Context) error
	RandflakeLeaseCreate(ctx context.Context) (*RandflakeLease, error)
	RandflakeLeaseExtend(ctx context.Context, prev *RandflakeLease) (*RandflakeLease, error)
}
