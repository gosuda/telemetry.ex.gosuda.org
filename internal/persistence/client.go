package persistence

import (
	"context"
	"database/sql"
	"time"

	"telemetry.ex.gosuda.org/telemetry/internal/persistence/database"
	"telemetry.ex.gosuda.org/telemetry/internal/types"
)

func (g *PersistenceClient) ClientRegister(ctx context.Context, id int64, token string) error {
	return g.db.ClientRegister(ctx, database.ClientRegisterParams{
		ID:        id,
		Token:     token,
		CreatedAt: time.Now().UnixNano(),
	})
}

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

func (g *PersistenceClient) ViewInsertWithCount(ctx context.Context, id int64, urlID int64, clientID int64, countID int64) error {
	// Start a transaction
	tx, err := g.pool.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create a new queries instance using the transaction
	txQueries := database.New(tx)
	now := time.Now().UnixNano()

	// Insert the view
	err = txQueries.ViewInsert(ctx, database.ViewInsertParams{
		ID:        id,
		UrlID:     urlID,
		ClientID:  clientID,
		CreatedAt: now,
	})
	if err != nil {
		return err
	}

	// Try to update the view count, if it doesn't exist, insert a new one
	err = txQueries.ViewCountUpdate(ctx, database.ViewCountUpdateParams{
		UpdatedAt: now,
		UrlID:     urlID,
	})
	if err != nil {
		// If update failed, try to insert a new view count
		err = txQueries.ViewCountInsert(ctx, database.ViewCountInsertParams{
			ID:        countID,
			UrlID:     urlID,
			UpdatedAt: now,
		})
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

func (g *PersistenceClient) UrlLookupByUrl(ctx context.Context, url string) (types.Url, error) {
	return g.db.UrlLookupByUrl(ctx, url)
}

func (g *PersistenceClient) UrlInsert(ctx context.Context, id int64, url string) error {
	return g.db.UrlInsert(ctx, database.UrlInsertParams{
		ID:        id,
		Url:       url,
		CreatedAt: time.Now().UnixNano(),
	})
}

func (g *PersistenceClient) ViewCountLookup(ctx context.Context, urlID int64) (types.ViewCount, error) {
	return g.db.ViewCountLookup(ctx, urlID)
}

func (g *PersistenceClient) LikeInsertWithCount(ctx context.Context, id int64, urlID int64, clientID int64, countID int64) error {
	// Start a transaction
	tx, err := g.pool.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create a new queries instance using the transaction
	txQueries := database.New(tx)
	now := time.Now().UnixNano()

	// Insert the like
	err = txQueries.LikeInsert(ctx, database.LikeInsertParams{
		ID:        id,
		UrlID:     urlID,
		ClientID:  clientID,
		CreatedAt: now,
	})
	if err != nil {
		return err
	}

	// Try to update the like count, if it doesn't exist, insert a new one
	err = txQueries.LikeCountUpdate(ctx, database.LikeCountUpdateParams{
		UpdatedAt: now,
		UrlID:     urlID,
	})
	if err != nil {
		// If update failed, try to insert a new like count
		err = txQueries.LikeCountInsert(ctx, database.LikeCountInsertParams{
			ID:        countID,
			UrlID:     urlID,
			UpdatedAt: now,
		})
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

func (g *PersistenceClient) LikeCountLookup(ctx context.Context, urlID int64) (types.LikeCount, error) {
	return g.db.LikeCountLookup(ctx, urlID)
}
