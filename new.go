package subscriptionstore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"github.com/gouniverse/sb"
)

// NewStoreOptions define the options for creating a new block store
type NewStoreOptions struct {
	PlanTableName         string
	SubscriptionTableName string
	DB                    *sql.DB
	DbDriverName          string
	AutomigrateEnabled    bool
	DebugEnabled          bool
	SqlLogger             *slog.Logger
}

// NewStore creates a new block store
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

	if opts.DbDriverName == "" {
		opts.DbDriverName = sb.DatabaseDriverName(opts.DB)
	}

	if opts.SqlLogger == nil {
		opts.SqlLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}

	store := &storeImplementation{
		planTableName:         opts.PlanTableName,
		subscriptionTableName: opts.SubscriptionTableName,
		automigrateEnabled:    opts.AutomigrateEnabled,
		db:                    opts.DB,
		dbDriverName:          opts.DbDriverName,
		debugEnabled:          opts.DebugEnabled,
		sqlLogger:             opts.SqlLogger,
	}

	if store.automigrateEnabled {
		err := store.AutoMigrate(context.Background())

		if err != nil {
			return nil, err
		}
	}

	return store, nil
}
