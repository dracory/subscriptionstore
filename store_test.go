package subscriptionstore

import (
	"context"
	"database/sql"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

func initDB(filepath string) *sql.DB {
	os.Remove(filepath) // remove database
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
