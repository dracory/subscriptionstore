package subscriptionstore

import (
	"encoding/json"
	"log"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
)

// SubscriptionInterface defines the methods for a Subscription entity
type SubscriptionInterface interface {
	IsSoftDeleted() bool

	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) SubscriptionInterface

	GetID() string
	SetID(id string) SubscriptionInterface

	GetStatus() string
	SetStatus(status string) SubscriptionInterface

	GetSubscriberID() string
	SetSubscriberID(subscriberID string) SubscriptionInterface

	GetPlanID() string
	SetPlanID(planID string) SubscriptionInterface

	GetPeriodStart() string
	GetPeriodStartCarbon() *carbon.Carbon
	SetPeriodStart(periodStart string) SubscriptionInterface

	GetPeriodEnd() string
	GetPeriodEndCarbon() *carbon.Carbon
	SetPeriodEnd(periodEnd string) SubscriptionInterface

	GetCancelAtPeriodEnd() bool
	SetCancelAtPeriodEnd(cancelAtPeriodEnd bool) SubscriptionInterface

	GetPaymentMethodID() string
	SetPaymentMethodID(paymentMethodID string) SubscriptionInterface

	GetMemo() string
	SetMemo(memo string) SubscriptionInterface

	GetMetas() (map[string]string, error)
	SetMetas(data map[string]string) (SubscriptionInterface, error)
	HasMeta(key string) (bool, error)
	Meta(key string) (string, error)
	SetMeta(key string, value string) (SubscriptionInterface, error)
	DeleteMeta(key string) (SubscriptionInterface, error)

	GetSoftDeletedAt() string
	GetSoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) SubscriptionInterface

	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) SubscriptionInterface
}

var _ SubscriptionInterface = (*subscriptionImplementation)(nil)

// == TYPE =====================================================================

type subscriptionImplementation struct {
	orm.ShortID

	StatusField            string `db:"status"`
	SubscriberIDField      string `db:"subscriber_id"`
	PlanIDField            string `db:"plan_id"`
	PeriodStartField       string `db:"period_start"`
	PeriodEndField         string `db:"period_end"`
	CancelAtPeriodEndField string `db:"cancel_at_period_end"`
	PaymentMethodIDField   string `db:"payment_method_id"`
	MemoField              string `db:"memo"`
	MetasField             string `db:"metas"`

	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate
}

// == CONSTRUCTORS =============================================================

func NewSubscription() SubscriptionInterface {
	o := &subscriptionImplementation{}
	o.SetID(neatuid.GenerateShortID())
	o.SetMemo("")
	o.SetSoftDeletedAt(MAX_DATETIME)
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetStatus(SUBSCRIPTION_STATUS_INACTIVE)
	o.SetPlanID("")
	o.SetSubscriberID("")
	o.SetPaymentMethodID("")
	o.SetPeriodStart(MAX_DATETIME)
	o.SetPeriodEnd(MAX_DATETIME)
	o.SetCancelAtPeriodEnd(false)
	if _, err := o.SetMetas(map[string]string{}); err != nil {
		log.Println(err.Error())
	}
	return o
}

func NewSubscriptionFromExistingData(data map[string]string) SubscriptionInterface {
	o := &subscriptionImplementation{}
	o.SetID(data[COLUMN_ID])
	o.SetStatus(data[COLUMN_STATUS])
	o.SetSubscriberID(data[COLUMN_SUBSCRIBER_ID])
	o.SetPlanID(data[COLUMN_PLAN_ID])
	o.SetPeriodStart(data[COLUMN_PERIOD_START])
	o.SetPeriodEnd(data[COLUMN_PERIOD_END])
	o.SetCancelAtPeriodEnd(data[COLUMN_CANCEL_AT_PERIOD_END] == YES)
	o.SetPaymentMethodID(data[COLUMN_PAYMENT_METHOD_ID])
	o.SetMemo(data[COLUMN_MEMO])
	o.MetasField = data[COLUMN_METAS]
	if v, ok := data[COLUMN_CREATED_AT]; ok {
		o.SetCreatedAt(v)
	}
	if v, ok := data[COLUMN_UPDATED_AT]; ok {
		o.SetUpdatedAt(v)
	}
	if v, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		o.SetSoftDeletedAt(v)
	}
	return o
}

// == METHODS ==================================================================

func (o *subscriptionImplementation) IsSoftDeleted() bool {
	return o.SoftDeletesMaxDate.IsSoftDeleted()
}

// == SETTERS AND GETTERS ======================================================

func (o *subscriptionImplementation) GetID() string {
	return o.ShortID.ID
}

func (o *subscriptionImplementation) SetID(id string) SubscriptionInterface {
	o.ShortID.ID = id
	return o
}

func (o *subscriptionImplementation) GetStatus() string {
	return o.StatusField
}

func (o *subscriptionImplementation) SetStatus(status string) SubscriptionInterface {
	o.StatusField = status
	return o
}

func (o *subscriptionImplementation) GetSubscriberID() string {
	return o.SubscriberIDField
}

func (o *subscriptionImplementation) SetSubscriberID(subscriberID string) SubscriptionInterface {
	o.SubscriberIDField = subscriberID
	return o
}

func (o *subscriptionImplementation) GetPlanID() string {
	return o.PlanIDField
}

func (o *subscriptionImplementation) SetPlanID(planID string) SubscriptionInterface {
	o.PlanIDField = planID
	return o
}

func (o *subscriptionImplementation) GetPeriodStart() string {
	if o.PeriodStartField == "" {
		return ""
	}
	return o.PeriodStartField
}

func (o *subscriptionImplementation) GetPeriodStartCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetPeriodStart(), carbon.UTC)
}

func (o *subscriptionImplementation) SetPeriodStart(periodStart string) SubscriptionInterface {
	o.PeriodStartField = periodStart
	return o
}

func (o *subscriptionImplementation) GetPeriodEnd() string {
	if o.PeriodEndField == "" {
		return ""
	}
	return o.PeriodEndField
}

func (o *subscriptionImplementation) GetPeriodEndCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetPeriodEnd(), carbon.UTC)
}

func (o *subscriptionImplementation) SetPeriodEnd(periodEnd string) SubscriptionInterface {
	o.PeriodEndField = periodEnd
	return o
}

func (o *subscriptionImplementation) GetCancelAtPeriodEnd() bool {
	return o.CancelAtPeriodEndField == YES
}

func (o *subscriptionImplementation) SetCancelAtPeriodEnd(cancelAtPeriodEnd bool) SubscriptionInterface {
	if cancelAtPeriodEnd {
		o.CancelAtPeriodEndField = YES
	} else {
		o.CancelAtPeriodEndField = NO
	}
	return o
}

func (o *subscriptionImplementation) GetPaymentMethodID() string {
	return o.PaymentMethodIDField
}

func (o *subscriptionImplementation) SetPaymentMethodID(paymentMethodID string) SubscriptionInterface {
	o.PaymentMethodIDField = paymentMethodID
	return o
}

func (o *subscriptionImplementation) GetMemo() string {
	return o.MemoField
}

func (o *subscriptionImplementation) SetMemo(memo string) SubscriptionInterface {
	o.MemoField = memo
	return o
}

func (o *subscriptionImplementation) GetMetas() (map[string]string, error) {
	if o.MetasField == "" {
		return nil, nil
	}
	var metas map[string]string
	err := json.Unmarshal([]byte(o.MetasField), &metas)
	if err != nil {
		return nil, err
	}
	return metas, nil
}

func (o *subscriptionImplementation) SetMetas(data map[string]string) (SubscriptionInterface, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	o.MetasField = string(jsonBytes)
	return o, nil
}

func (o *subscriptionImplementation) HasMeta(key string) (bool, error) {
	metas, err := o.GetMetas()
	if err != nil {
		return false, err
	}
	if metas == nil {
		return false, nil
	}
	_, exists := metas[key]
	return exists, nil
}

func (o *subscriptionImplementation) Meta(key string) (string, error) {
	metas, err := o.GetMetas()
	if err != nil {
		return "", err
	}
	if metas == nil {
		return "", nil
	}
	val, exists := metas[key]
	if !exists {
		return "", nil
	}
	return val, nil
}

func (o *subscriptionImplementation) SetMeta(key string, value string) (SubscriptionInterface, error) {
	metas, err := o.GetMetas()
	if err != nil {
		return nil, err
	}
	if metas == nil {
		metas = map[string]string{}
	}
	metas[key] = value
	return o.SetMetas(metas)
}

func (o *subscriptionImplementation) DeleteMeta(key string) (SubscriptionInterface, error) {
	metas, err := o.GetMetas()
	if err != nil {
		return nil, err
	}
	if metas == nil {
		return o, nil
	}
	delete(metas, key)
	return o.SetMetas(metas)
}

func (o *subscriptionImplementation) GetCreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

func (o *subscriptionImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

func (o *subscriptionImplementation) SetCreatedAt(createdAt string) SubscriptionInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

func (o *subscriptionImplementation) GetUpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

func (o *subscriptionImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

func (o *subscriptionImplementation) SetUpdatedAt(updatedAt string) SubscriptionInterface {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

func (o *subscriptionImplementation) GetSoftDeletedAt() string {
	if o.SoftDeletesMaxDate.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt).ToDateTimeString()
}

func (o *subscriptionImplementation) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt)
}

func (o *subscriptionImplementation) SetSoftDeletedAt(deletedAt string) SubscriptionInterface {
	if deletedAt == "" {
		return o
	}
	o.SoftDeletesMaxDate.SoftDeletedAt = carbon.Parse(deletedAt, carbon.UTC).StdTime()
	return o
}
