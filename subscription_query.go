package subscriptionstore

import "errors"

// SubscriptionQueryInterface defines the interface for querying subscriptions.
type SubscriptionQueryInterface interface {
	Validate() error

	IsCountOnly() bool
	HasCountOnly() bool
	SetCountOnly(countOnly bool) SubscriptionQueryInterface

	HasID() bool
	ID() string
	SetID(id string) SubscriptionQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) SubscriptionQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) SubscriptionQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) SubscriptionQueryInterface

	HasSubscriberID() bool
	SubscriberID() string
	SetSubscriberID(subscriberID string) SubscriptionQueryInterface

	HasPlanID() bool
	PlanID() string
	SetPlanID(planID string) SubscriptionQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) SubscriptionQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) SubscriptionQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) SubscriptionQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) SubscriptionQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(withSoftDeleted bool) SubscriptionQueryInterface
}

// SubscriptionQuery is a shortcut alias for NewSubscriptionQuery
func SubscriptionQuery() SubscriptionQueryInterface {
	return NewSubscriptionQuery()
}

// NewSubscriptionQuery creates a new subscription query
func NewSubscriptionQuery() SubscriptionQueryInterface {
	return &subscriptionQueryImplementation{
		properties: make(map[string]interface{}),
	}
}

var _ SubscriptionQueryInterface = (*subscriptionQueryImplementation)(nil)

type subscriptionQueryImplementation struct {
	properties map[string]interface{}
}

func (q *subscriptionQueryImplementation) Validate() error {
	if q.HasID() && q.ID() == "" {
		return errors.New("subscription query. id cannot be empty")
	}
	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("subscription query. id_in cannot be empty array")
	}
	if q.HasStatus() && q.Status() == "" {
		return errors.New("subscription query. status cannot be empty")
	}
	if q.HasStatusIn() && len(q.StatusIn()) < 1 {
		return errors.New("subscription query. status_in cannot be empty array")
	}
	if q.HasSubscriberID() && q.SubscriberID() == "" {
		return errors.New("subscription query. subscriber_id cannot be empty")
	}
	if q.HasPlanID() && q.PlanID() == "" {
		return errors.New("subscription query. plan_id cannot be empty")
	}
	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("subscription query. limit cannot be negative")
	}
	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("subscription query. offset cannot be negative")
	}
	return nil
}

func (q *subscriptionQueryImplementation) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

func (q *subscriptionQueryImplementation) IsCountOnly() bool {
	return q.hasProperty("count_only") && q.properties["count_only"].(bool)
}

func (q *subscriptionQueryImplementation) SetCountOnly(countOnly bool) SubscriptionQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *subscriptionQueryImplementation) HasID() bool {
	return q.hasProperty("id")
}

func (q *subscriptionQueryImplementation) ID() string {
	return q.properties["id"].(string)
}

func (q *subscriptionQueryImplementation) SetID(id string) SubscriptionQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *subscriptionQueryImplementation) HasIDIn() bool {
	return q.hasProperty("id_in")
}

func (q *subscriptionQueryImplementation) IDIn() []string {
	return q.properties["id_in"].([]string)
}

func (q *subscriptionQueryImplementation) SetIDIn(idIn []string) SubscriptionQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

func (q *subscriptionQueryImplementation) HasStatus() bool {
	return q.hasProperty("status")
}

func (q *subscriptionQueryImplementation) Status() string {
	return q.properties["status"].(string)
}

func (q *subscriptionQueryImplementation) SetStatus(status string) SubscriptionQueryInterface {
	q.properties["status"] = status
	return q
}

func (q *subscriptionQueryImplementation) HasStatusIn() bool {
	return q.hasProperty("status_in")
}

func (q *subscriptionQueryImplementation) StatusIn() []string {
	return q.properties["status_in"].([]string)
}

func (q *subscriptionQueryImplementation) SetStatusIn(statusIn []string) SubscriptionQueryInterface {
	q.properties["status_in"] = statusIn
	return q
}

func (q *subscriptionQueryImplementation) HasSubscriberID() bool {
	return q.hasProperty("subscriber_id")
}

func (q *subscriptionQueryImplementation) SubscriberID() string {
	return q.properties["subscriber_id"].(string)
}

func (q *subscriptionQueryImplementation) SetSubscriberID(subscriberID string) SubscriptionQueryInterface {
	q.properties["subscriber_id"] = subscriberID
	return q
}

func (q *subscriptionQueryImplementation) HasPlanID() bool {
	return q.hasProperty("plan_id")
}

func (q *subscriptionQueryImplementation) PlanID() string {
	return q.properties["plan_id"].(string)
}

func (q *subscriptionQueryImplementation) SetPlanID(planID string) SubscriptionQueryInterface {
	q.properties["plan_id"] = planID
	return q
}

func (q *subscriptionQueryImplementation) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *subscriptionQueryImplementation) Offset() int {
	return q.properties["offset"].(int)
}

func (q *subscriptionQueryImplementation) SetOffset(offset int) SubscriptionQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *subscriptionQueryImplementation) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *subscriptionQueryImplementation) Limit() int {
	return q.properties["limit"].(int)
}

func (q *subscriptionQueryImplementation) SetLimit(limit int) SubscriptionQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *subscriptionQueryImplementation) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *subscriptionQueryImplementation) OrderBy() string {
	return q.properties["order_by"].(string)
}

func (q *subscriptionQueryImplementation) SetOrderBy(orderBy string) SubscriptionQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *subscriptionQueryImplementation) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *subscriptionQueryImplementation) SortOrder() string {
	return q.properties["sort_order"].(string)
}

func (q *subscriptionQueryImplementation) SetSortOrder(sortOrder string) SubscriptionQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *subscriptionQueryImplementation) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_deleted_included")
}

func (q *subscriptionQueryImplementation) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties["soft_deleted_included"].(bool)
}

func (q *subscriptionQueryImplementation) SetSoftDeletedIncluded(softDeletedIncluded bool) SubscriptionQueryInterface {
	q.properties["soft_deleted_included"] = softDeletedIncluded
	return q
}

func (q *subscriptionQueryImplementation) hasProperty(key string) bool {
	_, ok := q.properties[key]
	return ok
}
