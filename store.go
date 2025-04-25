package subscriptionstore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/dracory/base/database"
)

var _ StoreInterface = (*storeImplementation)(nil) // verify it extends the interface

type storeImplementation struct {
	planTableName         string
	subscriptionTableName string
	db                    *sql.DB
	dbDriverName          string
	automigrateEnabled    bool
	debugEnabled          bool
	sqlLogger             *slog.Logger
}

// DatabaseDriverName returns the database driver name
func (store *storeImplementation) DatabaseDriverName() string {
	return store.dbDriverName
}

// PlanTableName returns the plan table name
func (store *storeImplementation) PlanTableName() string {
	return store.planTableName
}

// SubscriptionTableName returns the subscription table name
func (store *storeImplementation) SubscriptionTableName() string {
	return store.subscriptionTableName
}

// AutoMigrate auto migrates the database schema
func (store *storeImplementation) AutoMigrate(ctx context.Context) error {
	sql, err := store.sqlPlanTableCreate()

	if err != nil {
		return err
	}

	if sql == "" {
		return errors.New("subscription store: plan create sql is empty")
	}

	_, err = store.db.ExecContext(ctx, sql)

	if err != nil {
		return err
	}

	sql, err = store.sqlSubscriptionTableCreate()

	if err != nil {
		return err
	}

	if sql == "" {
		return errors.New("subscription store: subscription create sql is empty")
	}

	_, err = store.db.ExecContext(ctx, sql)

	if err != nil {
		return err
	}

	return nil
}

// EnableDebug enables the debug option
func (st *storeImplementation) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

// ============================================================================
// == PRIVATE METHODS
// ============================================================================

// logSql logs sql to the sql logger, if debug mode is enabled
func (store *storeImplementation) logSql(sqlOperationType string, sql string, params ...interface{}) {
	if !store.debugEnabled {
		return
	}

	if store.sqlLogger != nil {
		store.sqlLogger.Debug("sql: "+sqlOperationType, slog.String("sql", sql), slog.Any("params", params))
	}
}

// toQuerableContext converts the context to a QueryableContext
func (store *storeImplementation) toQuerableContext(ctx context.Context) database.QueryableContext {
	if database.IsQueryableContext(ctx) {
		return ctx.(database.QueryableContext)
	}

	return database.Context(ctx, store.db)
}
