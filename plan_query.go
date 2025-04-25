package subscriptionstore

import (
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
)

// PlanQuery is a shortcut alias for NewPlanQuery
func PlanQuery() PlanQueryInterface {
	return NewPlanQuery()
}

func NewPlanQuery() PlanQueryInterface {
	return &planQueryImplementation{}
}

type planQueryImplementation struct {
	id    string
	hasID bool

	idIn    []string
	hasIDIn bool

	status    string
	hasStatus bool

	statusIn    []string
	hasStatusIn bool

	type_   string
	hasType bool

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

func (o *planQueryImplementation) ID() string {
	return o.id
}

func (o *planQueryImplementation) HasID() bool {
	return o.hasID
}

func (o *planQueryImplementation) SetID(id string) PlanQueryInterface {
	o.id = id
	o.hasID = true
	return o
}

func (o *planQueryImplementation) IDIn() []string {
	return o.idIn
}

func (o *planQueryImplementation) HasIDIn() bool {
	return o.hasIDIn
}

func (o *planQueryImplementation) SetIDIn(idIn []string) PlanQueryInterface {
	o.idIn = idIn
	o.hasIDIn = true
	return o
}

func (o *planQueryImplementation) Status() string {
	return o.status
}

func (o *planQueryImplementation) HasStatus() bool {
	return o.hasStatus
}

func (o *planQueryImplementation) SetStatus(status string) PlanQueryInterface {
	o.status = status
	o.hasStatus = true
	return o
}

func (o *planQueryImplementation) StatusIn() []string {
	return o.statusIn
}

func (o *planQueryImplementation) HasStatusIn() bool {
	return o.hasStatusIn
}

func (o *planQueryImplementation) SetStatusIn(statusIn []string) PlanQueryInterface {
	o.statusIn = statusIn
	o.hasStatusIn = true
	return o
}

func (o *planQueryImplementation) Type() string {
	return o.type_
}

func (o *planQueryImplementation) HasType() bool {
	return o.hasType
}

func (o *planQueryImplementation) SetType(type_ string) PlanQueryInterface {
	o.type_ = type_
	o.hasType = true
	return o
}

func (o *planQueryImplementation) Offset() int {
	return o.offset
}

func (o *planQueryImplementation) HasOffset() bool {
	return o.hasOffset
}

func (o *planQueryImplementation) SetOffset(offset int) PlanQueryInterface {
	o.offset = offset
	o.hasOffset = true
	return o
}

func (o *planQueryImplementation) Limit() int {
	return o.limit
}

func (o *planQueryImplementation) HasLimit() bool {
	return o.hasLimit
}

func (o *planQueryImplementation) SetLimit(limit int) PlanQueryInterface {
	o.limit = limit
	o.hasLimit = true
	return o
}

func (o *planQueryImplementation) OrderBy() string {
	return o.orderBy
}

func (o *planQueryImplementation) HasOrderBy() bool {
	return o.hasOrderBy
}

func (o *planQueryImplementation) SetOrderBy(orderBy string) PlanQueryInterface {
	o.orderBy = orderBy
	o.hasOrderBy = true
	return o
}

func (o *planQueryImplementation) OrderDirection() string {
	return o.orderDirection
}

func (o *planQueryImplementation) HasOrderDirection() bool {
	return o.hasOrderDirection
}

func (o *planQueryImplementation) SetOrderDirection(orderDirection string) PlanQueryInterface {
	o.orderDirection = orderDirection
	o.hasOrderDirection = true
	return o
}

func (o *planQueryImplementation) CountOnly() bool {
	return o.countOnly
}

func (o *planQueryImplementation) HasCountOnly() bool {
	return o.hasCountOnly
}

func (o *planQueryImplementation) SetCountOnly(countOnly bool) PlanQueryInterface {
	o.countOnly = countOnly
	o.hasCountOnly = true
	return o
}

func (o *planQueryImplementation) WithDeleted() bool {
	return o.withDeleted
}

func (o *planQueryImplementation) HasWithDeleted() bool {
	return o.hasWithDeleted
}

func (o *planQueryImplementation) SetWithDeleted(withDeleted bool) PlanQueryInterface {
	o.withDeleted = withDeleted
	o.hasWithDeleted = true
	return o
}

func (o *planQueryImplementation) ToQuery(store StoreInterface) *goqu.SelectDataset {

	q := goqu.Dialect(store.DatabaseDriverName()).From(store.PlanTableName())

	if o.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(o.ID()))
	}

	if o.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(o.IDIn()))
	}

	if o.HasStatus() {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(o.Status()))
	}

	if o.HasStatusIn() {
		q = q.Where(goqu.C(COLUMN_STATUS).In(o.StatusIn()))
	}

	if o.HasType() {
		q = q.Where(goqu.C(COLUMN_TYPE).Eq(o.Type()))
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
