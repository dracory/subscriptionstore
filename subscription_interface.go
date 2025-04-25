package subscriptionstore

import "github.com/dromara/carbon/v2"

// SubscriptionInterface defines the methods for a Subscription entity
// This interface can be implemented by any Subscription struct for flexibility and testability.
type SubscriptionInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	CreatedAt() string
	SetCreatedAt(createdAt string) SubscriptionInterface
	CreatedAtCarbon() *carbon.Carbon

	ID() string
	SetID(id string) SubscriptionInterface

	Status() string
	SetStatus(status string) SubscriptionInterface

	SubscriberID() string
	SetSubscriberID(subscriberID string) SubscriptionInterface

	PlanID() string
	SetPlanID(planID string) SubscriptionInterface

	PeriodStart() string
	SetPeriodStart(periodStart string) SubscriptionInterface

	PeriodEnd() string
	SetPeriodEnd(periodEnd string) SubscriptionInterface

	CancelAtPeriodEnd() bool
	SetCancelAtPeriodEnd(cancelAtPeriodEnd bool) SubscriptionInterface

	PaymentMethodID() string
	SetPaymentMethodID(paymentMethodID string) SubscriptionInterface

	Memo() string
	SetMemo(memo string) SubscriptionInterface

	Metas() (map[string]string, error)
	SetMetas(data map[string]string) (SubscriptionInterface, error)

	HasMeta(key string) (bool, error)
	Meta(key string) (string, error)
	SetMeta(key string, value string) (SubscriptionInterface, error)
	DeleteMeta(key string) (SubscriptionInterface, error)

	SoftDeletedAt() string
	SetSoftDeletedAt(softDeletedAt string) SubscriptionInterface
	SoftDeletedAtCarbon() *carbon.Carbon

	UpdatedAt() string
	SetUpdatedAt(updatedAt string) SubscriptionInterface
	UpdatedAtCarbon() *carbon.Carbon
}
