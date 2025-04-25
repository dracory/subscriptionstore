package subscriptionstore

import (
	"github.com/doug-martin/goqu/v9"
)

// PlanQueryInterface defines the methods for querying plans using PlanQueryOptions.
type PlanQueryInterface interface {
	ToQuery(store StoreInterface) *goqu.SelectDataset

	ID() string
	HasID() bool
	SetID(id string) PlanQueryInterface

	IDIn() []string
	HasIDIn() bool
	SetIDIn(idIn []string) PlanQueryInterface

	Status() string
	HasStatus() bool
	SetStatus(status string) PlanQueryInterface

	StatusIn() []string
	HasStatusIn() bool
	SetStatusIn(statusIn []string) PlanQueryInterface

	Type() string
	HasType() bool
	SetType(type_ string) PlanQueryInterface

	Offset() int
	HasOffset() bool
	SetOffset(offset int) PlanQueryInterface

	Limit() int
	HasLimit() bool
	SetLimit(limit int) PlanQueryInterface

	OrderBy() string
	HasOrderBy() bool
	SetOrderBy(orderBy string) PlanQueryInterface

	OrderDirection() string
	HasOrderDirection() bool
	SetOrderDirection(orderByDirection string) PlanQueryInterface

	CountOnly() bool
	HasCountOnly() bool
	SetCountOnly(countOnly bool) PlanQueryInterface

	WithDeleted() bool
	HasWithDeleted() bool
	SetWithDeleted(withDeleted bool) PlanQueryInterface
}
