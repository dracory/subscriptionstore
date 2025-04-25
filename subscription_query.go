package subscriptionstore

import (
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
)

// SubscriptionQuery is a shortcut alias for NewSubscriptionQuery
func SubscriptionQuery() SubscriptionQueryInterface {
	return NewSubscriptionQuery()
}

func NewSubscriptionQuery() SubscriptionQueryInterface {
	return &subscriptionQueryImplementation{}
}

type subscriptionQueryImplementation struct {
	id    string
	hasID bool

	idIn    []string
	hasIDIn bool

	status    string
	hasStatus bool

	statusIn    []string
	hasStatusIn bool

	subscriberID    string
	hasSubscriberID bool

	planID    string
	hasPlanID bool

	offset    int
	hasOffset bool

	limit    int
	hasLimit bool

	orderBy    string
	hasOrderBy bool

	orderDirection    string
	hasOrderDirection bool

	countOnly    bool
	hasCountOnly bool

	withDeleted    bool
	hasWithDeleted bool
}

// Implement SubscriptionQueryInterface methods
func (o *subscriptionQueryImplementation) ID() string  { return o.id }
func (o *subscriptionQueryImplementation) HasID() bool { return o.hasID }
func (o *subscriptionQueryImplementation) SetID(id string) SubscriptionQueryInterface {
	o.id = id
	o.hasID = true
	return o
}

func (o *subscriptionQueryImplementation) IDIn() []string { return o.idIn }
func (o *subscriptionQueryImplementation) HasIDIn() bool  { return o.hasIDIn }
func (o *subscriptionQueryImplementation) SetIDIn(idIn []string) SubscriptionQueryInterface {
	o.idIn = idIn
	o.hasIDIn = true
	return o
}

func (o *subscriptionQueryImplementation) Status() string  { return o.status }
func (o *subscriptionQueryImplementation) HasStatus() bool { return o.hasStatus }
func (o *subscriptionQueryImplementation) SetStatus(status string) SubscriptionQueryInterface {
	o.status = status
	o.hasStatus = true
	return o
}

func (o *subscriptionQueryImplementation) StatusIn() []string { return o.statusIn }
func (o *subscriptionQueryImplementation) HasStatusIn() bool  { return o.hasStatusIn }
func (o *subscriptionQueryImplementation) SetStatusIn(statusIn []string) SubscriptionQueryInterface {
	o.statusIn = statusIn
	o.hasStatusIn = true
	return o
}

func (o *subscriptionQueryImplementation) SubscriberID() string  { return o.subscriberID }
func (o *subscriptionQueryImplementation) HasSubscriberID() bool { return o.hasSubscriberID }
func (o *subscriptionQueryImplementation) SetSubscriberID(subscriberID string) SubscriptionQueryInterface {
	o.subscriberID = subscriberID
	o.hasSubscriberID = true
	return o
}

func (o *subscriptionQueryImplementation) PlanID() string  { return o.planID }
func (o *subscriptionQueryImplementation) HasPlanID() bool { return o.hasPlanID }
func (o *subscriptionQueryImplementation) SetPlanID(planID string) SubscriptionQueryInterface {
	o.planID = planID
	o.hasPlanID = true
	return o
}

func (o *subscriptionQueryImplementation) Offset() int     { return o.offset }
func (o *subscriptionQueryImplementation) HasOffset() bool { return o.hasOffset }
func (o *subscriptionQueryImplementation) SetOffset(offset int) SubscriptionQueryInterface {
	o.offset = offset
	o.hasOffset = true
	return o
}

func (o *subscriptionQueryImplementation) Limit() int     { return o.limit }
func (o *subscriptionQueryImplementation) HasLimit() bool { return o.hasLimit }
func (o *subscriptionQueryImplementation) SetLimit(limit int) SubscriptionQueryInterface {
	o.limit = limit
	o.hasLimit = true
	return o
}

func (o *subscriptionQueryImplementation) OrderBy() string  { return o.orderBy }
func (o *subscriptionQueryImplementation) HasOrderBy() bool { return o.hasOrderBy }
func (o *subscriptionQueryImplementation) SetOrderBy(orderBy string) SubscriptionQueryInterface {
	o.orderBy = orderBy
	o.hasOrderBy = true
	return o
}

func (o *subscriptionQueryImplementation) OrderDirection() string  { return o.orderDirection }
func (o *subscriptionQueryImplementation) HasOrderDirection() bool { return o.hasOrderDirection }
func (o *subscriptionQueryImplementation) SetOrderDirection(orderDirection string) SubscriptionQueryInterface {
	o.orderDirection = orderDirection
	o.hasOrderDirection = true
	return o
}

func (o *subscriptionQueryImplementation) CountOnly() bool    { return o.countOnly }
func (o *subscriptionQueryImplementation) HasCountOnly() bool { return o.hasCountOnly }
func (o *subscriptionQueryImplementation) SetCountOnly(countOnly bool) SubscriptionQueryInterface {
	o.countOnly = countOnly
	o.hasCountOnly = true
	return o
}

func (o *subscriptionQueryImplementation) WithDeleted() bool    { return o.withDeleted }
func (o *subscriptionQueryImplementation) HasWithDeleted() bool { return o.hasWithDeleted }
func (o *subscriptionQueryImplementation) SetWithDeleted(withDeleted bool) SubscriptionQueryInterface {
	o.withDeleted = withDeleted
	o.hasWithDeleted = true
	return o
}

// ToQuery builds a goqu.SelectDataset based on the query options
func (o *subscriptionQueryImplementation) ToQuery(store StoreInterface) *goqu.SelectDataset {
	q := goqu.Dialect(store.DatabaseDriverName()).From(store.SubscriptionTableName())

	if o.hasID {
		q = q.Where(goqu.C(COLUMN_ID).Eq(o.ID()))
	}
	if o.hasIDIn && len(o.idIn) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(o.IDIn()))
	}
	if o.hasStatus {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(o.Status()))
	}
	if o.hasStatusIn && len(o.statusIn) > 0 {
		q = q.Where(goqu.C(COLUMN_STATUS).In(o.StatusIn()))
	}
	if o.hasSubscriberID {
		q = q.Where(goqu.C(COLUMN_SUBSCRIBER_ID).Eq(o.SubscriberID()))
	}
	if o.hasPlanID {
		q = q.Where(goqu.C(COLUMN_PLAN_ID).Eq(o.PlanID()))
	}

	if !o.HasCountOnly() && o.HasOffset() {
		if o.Offset() > 0 {
			q = q.Offset(uint(o.Offset()))
		}
	}

	if !o.HasCountOnly() && o.HasLimit() {
		if o.Limit() > 0 {
			q = q.Limit(uint(o.Limit()))
		}
	}

	sortOrder := lo.Ternary(o.HasOrderDirection(), o.OrderDirection(), "desc")

	if o.HasOrderBy() {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(o.OrderBy()).Asc())
		} else {
			q = q.Order(goqu.I(o.OrderBy()).Desc())
		}
	}

	if o.WithDeleted() {
		return q
	}

	q = q.Where(goqu.C(COLUMN_SOFT_DELETED_AT).Gt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)))

	return q
}
