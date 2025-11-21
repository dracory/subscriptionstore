package subscriptionstore

import (
	"context"
	"database/sql"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

func initDB(filepath string) *sql.DB {
	if err := os.Remove(filepath); err != nil && !errors.Is(err, os.ErrNotExist) { // remove database
		panic(err)
	}
	dsn := filepath + "?parseTime=true&loc=UTC&_loc=UTC"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func initStore() (StoreInterface, error) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                    db,
		PlanTableName:         "plan_table",
		SubscriptionTableName: "subscription_table",
		AutomigrateEnabled:    true,
	})

	if err != nil {
		return nil, err
	}

	if store == nil {
		return nil, errors.New("unexpected nil store")
	}

	err = store.AutoMigrate(context.Background())
	if err != nil {
		return nil, err
	}

	return store, nil
}
