package subscriptionstore

import "errors"

// PlanQueryInterface defines the interface for querying plans.
type PlanQueryInterface interface {
	Validate() error

	IsCountOnly() bool
	HasCountOnly() bool
	SetCountOnly(countOnly bool) PlanQueryInterface

	HasID() bool
	ID() string
	SetID(id string) PlanQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) PlanQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) PlanQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) PlanQueryInterface

	HasInterval() bool
	Interval() string
	SetInterval(interval string) PlanQueryInterface

	HasIntervalIn() bool
	IntervalIn() []string
	SetIntervalIn(intervalIn []string) PlanQueryInterface

	HasType() bool
	Type() string
	SetType(type_ string) PlanQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) PlanQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) PlanQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) PlanQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) PlanQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(withSoftDeleted bool) PlanQueryInterface
}

// PlanQuery is a shortcut alias for NewPlanQuery
func PlanQuery() PlanQueryInterface {
	return NewPlanQuery()
}

// NewPlanQuery creates a new plan query
func NewPlanQuery() PlanQueryInterface {
	return &planQueryImplementation{
		properties: make(map[string]interface{}),
	}
}

var _ PlanQueryInterface = (*planQueryImplementation)(nil)

type planQueryImplementation struct {
	properties map[string]interface{}
}

func (q *planQueryImplementation) Validate() error {
	if q.HasID() && q.ID() == "" {
		return errors.New("plan query. id cannot be empty")
	}
	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("plan query. id_in cannot be empty array")
	}
	if q.HasStatus() && q.Status() == "" {
		return errors.New("plan query. status cannot be empty")
	}
	if q.HasStatusIn() && len(q.StatusIn()) < 1 {
		return errors.New("plan query. status_in cannot be empty array")
	}
	if q.HasInterval() && q.Interval() == "" {
		return errors.New("plan query. interval cannot be empty")
	}
	if q.HasIntervalIn() && len(q.IntervalIn()) < 1 {
		return errors.New("plan query. interval_in cannot be empty array")
	}
	if q.HasType() && q.Type() == "" {
		return errors.New("plan query. type cannot be empty")
	}
	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("plan query. limit cannot be negative")
	}
	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("plan query. offset cannot be negative")
	}
	return nil
}

func (q *planQueryImplementation) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

func (q *planQueryImplementation) IsCountOnly() bool {
	return q.hasProperty("count_only") && q.properties["count_only"].(bool)
}

func (q *planQueryImplementation) SetCountOnly(countOnly bool) PlanQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *planQueryImplementation) HasID() bool {
	return q.hasProperty("id")
}

func (q *planQueryImplementation) ID() string {
	return q.properties["id"].(string)
}

func (q *planQueryImplementation) SetID(id string) PlanQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *planQueryImplementation) HasIDIn() bool {
	return q.hasProperty("id_in")
}

func (q *planQueryImplementation) IDIn() []string {
	return q.properties["id_in"].([]string)
}

func (q *planQueryImplementation) SetIDIn(idIn []string) PlanQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

func (q *planQueryImplementation) HasStatus() bool {
	return q.hasProperty("status")
}

func (q *planQueryImplementation) Status() string {
	return q.properties["status"].(string)
}

func (q *planQueryImplementation) SetStatus(status string) PlanQueryInterface {
	q.properties["status"] = status
	return q
}

func (q *planQueryImplementation) HasStatusIn() bool {
	return q.hasProperty("status_in")
}

func (q *planQueryImplementation) StatusIn() []string {
	return q.properties["status_in"].([]string)
}

func (q *planQueryImplementation) SetStatusIn(statusIn []string) PlanQueryInterface {
	q.properties["status_in"] = statusIn
	return q
}

func (q *planQueryImplementation) HasInterval() bool {
	return q.hasProperty("interval")
}

func (q *planQueryImplementation) Interval() string {
	return q.properties["interval"].(string)
}

func (q *planQueryImplementation) SetInterval(interval string) PlanQueryInterface {
	q.properties["interval"] = interval
	return q
}

func (q *planQueryImplementation) HasIntervalIn() bool {
	return q.hasProperty("interval_in")
}

func (q *planQueryImplementation) IntervalIn() []string {
	return q.properties["interval_in"].([]string)
}

func (q *planQueryImplementation) SetIntervalIn(intervalIn []string) PlanQueryInterface {
	q.properties["interval_in"] = intervalIn
	return q
}

func (q *planQueryImplementation) HasType() bool {
	return q.hasProperty("type")
}

func (q *planQueryImplementation) Type() string {
	return q.properties["type"].(string)
}

func (q *planQueryImplementation) SetType(type_ string) PlanQueryInterface {
	q.properties["type"] = type_
	return q
}

func (q *planQueryImplementation) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *planQueryImplementation) Offset() int {
	return q.properties["offset"].(int)
}

func (q *planQueryImplementation) SetOffset(offset int) PlanQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *planQueryImplementation) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *planQueryImplementation) Limit() int {
	return q.properties["limit"].(int)
}

func (q *planQueryImplementation) SetLimit(limit int) PlanQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *planQueryImplementation) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *planQueryImplementation) OrderBy() string {
	return q.properties["order_by"].(string)
}

func (q *planQueryImplementation) SetOrderBy(orderBy string) PlanQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *planQueryImplementation) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *planQueryImplementation) SortOrder() string {
	return q.properties["sort_order"].(string)
}

func (q *planQueryImplementation) SetSortOrder(sortOrder string) PlanQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *planQueryImplementation) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_deleted_included")
}

func (q *planQueryImplementation) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties["soft_deleted_included"].(bool)
}

func (q *planQueryImplementation) SetSoftDeletedIncluded(softDeletedIncluded bool) PlanQueryInterface {
	q.properties["soft_deleted_included"] = softDeletedIncluded
	return q
}

func (q *planQueryImplementation) hasProperty(key string) bool {
	_, ok := q.properties[key]
	return ok
}
