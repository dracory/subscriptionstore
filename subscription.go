package subscriptionstore

import (
	"encoding/json"

	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/uid"
)

// SubscriptionInterface defines the methods for a Subscription entity
// This interface can be implemented by any Subscription struct for flexibility and testability.
var _ SubscriptionInterface = (*subscriptionImplementation)(nil) // verify it extends the interface

// ============================================================================
// == Subscription
// ============================================================================

type subscriptionImplementation struct {
	dataobject.DataObject
}

// ============================================================================
// == Constructors
// ============================================================================

func NewSubscription() SubscriptionInterface {
	o := &subscriptionImplementation{}

	o.SetID(uid.HumanUid())
	o.SetMemo("")
	o.SetSoftDeletedAt(carbon.MaxValue().ToDateTimeString())
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	o.SetStatus(SUBSCRIPTION_STATUS_INACTIVE)
	o.SetPlanID("")
	o.SetSubscriberID("")
	o.SetPaymentMethodID("")
	o.SetPeriodStart(carbon.MaxValue().ToDateTimeString())
	o.SetPeriodEnd(carbon.MaxValue().ToDateTimeString())
	o.SetCancelAtPeriodEnd(false)
	o.SetMetas(map[string]string{})
	return o
}

func NewSubscriptionFromExistingData(data map[string]string) SubscriptionInterface {
	o := &subscriptionImplementation{}
	o.Hydrate(data)
	return o
}

// ============================================================================
// == Setters and Getters
// ============================================================================

func (o *subscriptionImplementation) CancelAtPeriodEnd() bool {
	return o.Get(COLUMN_CANCEL_AT_PERIOD_END) == YES
}

func (o *subscriptionImplementation) SetCancelAtPeriodEnd(cancelAtPeriodEnd bool) SubscriptionInterface {
	if cancelAtPeriodEnd {
		o.Set(COLUMN_CANCEL_AT_PERIOD_END, YES)
	} else {
		o.Set(COLUMN_CANCEL_AT_PERIOD_END, NO)
	}
	return o
}

func (o *subscriptionImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *subscriptionImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_CREATED_AT), carbon.UTC)
}

func (o *subscriptionImplementation) SetCreatedAt(createdAt string) SubscriptionInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *subscriptionImplementation) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *subscriptionImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_SOFT_DELETED_AT), carbon.UTC)
}

func (o *subscriptionImplementation) SetSoftDeletedAt(softDeletedAt string) SubscriptionInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *subscriptionImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *subscriptionImplementation) SetID(id string) SubscriptionInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *subscriptionImplementation) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *subscriptionImplementation) SetMemo(memo string) SubscriptionInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *subscriptionImplementation) Metas() (map[string]string, error) {
	metasString := o.Get(COLUMN_METAS)

	if metasString == "" {
		return nil, nil
	}

	var metas map[string]string
	err := json.Unmarshal([]byte(metasString), &metas)

	if err != nil {
		return nil, err
	}

	return metas, nil
}

func (o *subscriptionImplementation) SetMetas(data map[string]string) (SubscriptionInterface, error) {
	json, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	o.Set(COLUMN_METAS, string(json))
	return o, nil
}

func (o *subscriptionImplementation) PaymentMethodID() string {
	return o.Get(COLUMN_PAYMENT_METHOD_ID)
}

func (o *subscriptionImplementation) SetPaymentMethodID(paymentMethodID string) SubscriptionInterface {
	o.Set(COLUMN_PAYMENT_METHOD_ID, paymentMethodID)
	return o
}

func (o *subscriptionImplementation) SetPeriodStart(periodStart string) SubscriptionInterface {
	o.Set(COLUMN_PERIOD_START, periodStart)
	return o
}

func (o *subscriptionImplementation) PeriodEnd() string {
	return o.Get(COLUMN_PERIOD_END)
}

func (o *subscriptionImplementation) SetPeriodEnd(periodEnd string) SubscriptionInterface {
	o.Set(COLUMN_PERIOD_END, periodEnd)
	return o
}

func (o *subscriptionImplementation) PlanID() string {
	return o.Get(COLUMN_PLAN_ID)
}

func (o *subscriptionImplementation) SetPlanID(planID string) SubscriptionInterface {
	o.Set(COLUMN_PLAN_ID, planID)
	return o
}

func (o *subscriptionImplementation) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *subscriptionImplementation) SetStatus(status string) SubscriptionInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *subscriptionImplementation) SubscriberID() string {
	return o.Get(COLUMN_SUBSCRIBER_ID)
}

func (o *subscriptionImplementation) SetSubscriberID(subscriberID string) SubscriptionInterface {
	o.Set(COLUMN_SUBSCRIBER_ID, subscriberID)
	return o
}

func (o *subscriptionImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *subscriptionImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_UPDATED_AT), carbon.UTC)
}

func (o *subscriptionImplementation) SetUpdatedAt(updatedAt string) SubscriptionInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *subscriptionImplementation) PeriodStart() string {
	return o.Get(COLUMN_PERIOD_START)
}

func (o *subscriptionImplementation) HasMeta(key string) (bool, error) {
	metas, err := o.Metas()
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
	metas, err := o.Metas()
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
	metas, err := o.Metas()
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
	metas, err := o.Metas()
	if err != nil {
		return nil, err
	}
	if metas == nil {
		return o, nil
	}
	delete(metas, key)
	return o.SetMetas(metas)
}
