package subscriptionstore

import (
	"github.com/doug-martin/goqu/v9"
)

// SubscriptionQueryInterface defines the methods for querying subscriptions using SubscriptionQueryOptions.
type SubscriptionQueryInterface interface {
	ToQuery(store StoreInterface) *goqu.SelectDataset

	ID() string
	HasID() bool
	SetID(id string) SubscriptionQueryInterface

	IDIn() []string
	HasIDIn() bool
	SetIDIn(idIn []string) SubscriptionQueryInterface

	Status() string
	HasStatus() bool
	SetStatus(status string) SubscriptionQueryInterface

	StatusIn() []string
	HasStatusIn() bool
	SetStatusIn(statusIn []string) SubscriptionQueryInterface

	SubscriberID() string
	HasSubscriberID() bool
	SetSubscriberID(subscriberID string) SubscriptionQueryInterface

	PlanID() string
	HasPlanID() bool
	SetPlanID(planID string) SubscriptionQueryInterface

	Offset() int
	HasOffset() bool
	SetOffset(offset int) SubscriptionQueryInterface

	Limit() int
	HasLimit() bool
	SetLimit(limit int) SubscriptionQueryInterface

	OrderBy() string
	HasOrderBy() bool
	SetOrderBy(orderBy string) SubscriptionQueryInterface

	OrderDirection() string
	HasOrderDirection() bool
	SetOrderDirection(orderByDirection string) SubscriptionQueryInterface

	CountOnly() bool
	HasCountOnly() bool
	SetCountOnly(countOnly bool) SubscriptionQueryInterface

	WithDeleted() bool
	HasWithDeleted() bool
	SetWithDeleted(withDeleted bool) SubscriptionQueryInterface
}
