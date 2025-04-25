package subscriptionstore

import "github.com/dromara/carbon/v2"

// PlanInterface defines the methods for a Plan entity
// This interface can be implemented by any Plan struct for flexibility and testability.
type PlanInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	CreatedAt() string
	SetCreatedAt(createdAt string) PlanInterface
	CreatedAtCarbon() *carbon.Carbon

	Currency() string
	SetCurrency(currency string) PlanInterface

	SoftDeletedAt() string
	SetSoftDeletedAt(softDeletedAt string) PlanInterface
	SoftDeletedAtCarbon() *carbon.Carbon

	Description() string
	SetDescription(description string) PlanInterface

	ID() string
	SetID(id string) PlanInterface

	Features() string
	SetFeatures(features string) PlanInterface

	Interval() string
	SetInterval(interval string) PlanInterface

	Memo() string
	SetMemo(memo string) PlanInterface

	Metas() (map[string]string, error)
	SetMetas(data map[string]string) (PlanInterface, error)

	HasMeta(key string) (bool, error)
	Meta(key string) (string, error)
	SetMeta(key string, value string) (PlanInterface, error)
	DeleteMeta(key string) (PlanInterface, error)

	Price() string
	PriceFloat() float64
	SetPrice(price string) PlanInterface

	Status() string
	SetStatus(status string) PlanInterface

	Title() string
	SetTitle(title string) PlanInterface

	Type() string
	SetType(type_ string) PlanInterface

	StripePriceID() string
	SetStripePriceID(stripePriceID string) PlanInterface

	UpdatedAt() string
	SetUpdatedAt(updatedAt string) PlanInterface
	UpdatedAtCarbon() *carbon.Carbon
}
