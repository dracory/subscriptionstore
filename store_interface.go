package subscriptionstore

import "context"

// StoreInterface defines the methods for the Subscription store
// This interface allows for easier testing and separation of concerns
// between the store implementation and its consumers.
type StoreInterface interface {
	AutoMigrate(ctx context.Context) error
	EnableDebug(debug bool)

	DatabaseDriverName() string

	PlanCount(ctx context.Context, query PlanQueryInterface) (int64, error)
	PlanCreate(ctx context.Context, plan PlanInterface) error
	PlanDelete(ctx context.Context, plan PlanInterface) error
	PlanDeleteByID(ctx context.Context, id string) error
	PlanExists(ctx context.Context, planID string) (bool, error)
	PlanFindByID(ctx context.Context, id string) (PlanInterface, error)
	PlanList(ctx context.Context, query PlanQueryInterface) ([]PlanInterface, error)
	PlanSoftDelete(ctx context.Context, plan PlanInterface) error
	PlanSoftDeleteByID(ctx context.Context, id string) error
	PlanTableName() string
	PlanUpdate(ctx context.Context, plan PlanInterface) error

	SubscriptionCount(ctx context.Context, query SubscriptionQueryInterface) (int64, error)
	SubscriptionCreate(ctx context.Context, subscription SubscriptionInterface) error
	SubscriptionDelete(ctx context.Context, subscription SubscriptionInterface) error
	SubscriptionDeleteByID(ctx context.Context, id string) error
	SubscriptionExists(ctx context.Context, subscriptionID string) (bool, error)
	SubscriptionFindByID(ctx context.Context, id string) (SubscriptionInterface, error)
	SubscriptionList(ctx context.Context, query SubscriptionQueryInterface) ([]SubscriptionInterface, error)
	SubscriptionSoftDelete(ctx context.Context, subscription SubscriptionInterface) error
	SubscriptionSoftDeleteByID(ctx context.Context, id string) error
	SubscriptionTableName() string
	SubscriptionUpdate(ctx context.Context, subscription SubscriptionInterface) error
}
