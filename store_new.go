package subscriptionstore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"github.com/dracory/neat"
)

// NewStoreOptions define the options for creating a new subscription store
type NewStoreOptions struct {
	PlanTableName         string
	SubscriptionTableName string
	DB                    *sql.DB
	AutomigrateEnabled    bool
	DebugEnabled          bool
}

// NewStore creates a new subscription store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.PlanTableName == "" {
		return nil, errors.New("subscription store: PlanTableName is required")
	}

	if opts.SubscriptionTableName == "" {
		return nil, errors.New("subscription store: SubscriptionTableName is required")
	}

	if opts.DB == nil {
		return nil, errors.New("subscription store: DB is required")
	}

	neatDB, err := neat.NewFromSQLDB(opts.DB)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	store := &storeImplementation{
		planTableName:         opts.PlanTableName,
		subscriptionTableName: opts.SubscriptionTableName,
		db:                    neatDB,
		automigrateEnabled:    opts.AutomigrateEnabled,
		debugEnabled:          opts.DebugEnabled,
		sqlLogger:             logger,
	}

	if store.automigrateEnabled {
		if err := store.MigrateUp(context.Background()); err != nil {
			return nil, err
		}
	}

	return store, nil
}
