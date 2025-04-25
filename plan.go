package subscriptionstore

import (
	"encoding/json"

	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/uid"
	"github.com/spf13/cast"
)

var _ PlanInterface = (*planImplementation)(nil)

type planImplementation struct {
	dataobject.DataObject
}

func NewPlan() PlanInterface {
	o := &planImplementation{}

	o.SetID(uid.HumanUid())
	o.SetStatus(PLAN_STATUS_INACTIVE)
	o.SetStripePriceID("")
	o.SetDescription("")
	o.SetFeatures("")
	o.SetMemo("")
	o.SetMetas(map[string]string{})
	o.SetSoftDeletedAt(carbon.MaxValue().ToDateTimeString())
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	return o
}

func NewPlanFromExistingData(data map[string]string) PlanInterface {
	o := &planImplementation{}
	o.Hydrate(data)
	return o
}

func (o *planImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *planImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_CREATED_AT), carbon.UTC)
}

func (o *planImplementation) SetCreatedAt(createdAt string) PlanInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *planImplementation) Currency() string {
	return o.Get(COLUMN_CURRENCY)
}

func (o *planImplementation) SetCurrency(currency string) PlanInterface {
	o.Set(COLUMN_CURRENCY, currency)
	return o
}

func (o *planImplementation) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *planImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_SOFT_DELETED_AT), carbon.UTC)
}

func (o *planImplementation) SetSoftDeletedAt(softDeletedAt string) PlanInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *planImplementation) Description() string {
	return o.Get(COLUMN_DESCRIPTION)
}

func (o *planImplementation) SetDescription(description string) PlanInterface {
	o.Set(COLUMN_DESCRIPTION, description)
	return o
}

func (o *planImplementation) Features() string {
	return o.Get(COLUMN_FEATURES)
}

func (o *planImplementation) SetFeatures(features string) PlanInterface {
	o.Set(COLUMN_FEATURES, features)
	return o
}

func (o *planImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *planImplementation) SetID(id string) PlanInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *planImplementation) Interval() string {
	return o.Get(COLUMN_INTERVAL)
}

func (o *planImplementation) SetInterval(interval string) PlanInterface {
	o.Set(COLUMN_INTERVAL, interval)
	return o
}

func (o *planImplementation) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *planImplementation) SetMemo(memo string) PlanInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *planImplementation) Metas() (map[string]string, error) {
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

func (o *planImplementation) SetMetas(data map[string]string) (PlanInterface, error) {
	json, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	o.Set(COLUMN_METAS, string(json))
	return o, nil
}

func (o *planImplementation) HasMeta(key string) (bool, error) {
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

func (o *planImplementation) Meta(key string) (string, error) {
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

func (o *planImplementation) SetMeta(key string, value string) (PlanInterface, error) {
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

func (o *planImplementation) DeleteMeta(key string) (PlanInterface, error) {
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

func (o *planImplementation) Price() string {
	return o.Get(COLUMN_PRICE)
}

func (o *planImplementation) PriceFloat() float64 {
	return cast.ToFloat64(o.Get(COLUMN_PRICE))
}

func (o *planImplementation) SetPrice(price string) PlanInterface {
	o.Set(COLUMN_PRICE, price)
	return o
}

func (o *planImplementation) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *planImplementation) SetStatus(status string) PlanInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *planImplementation) StripePriceID() string {
	return o.Get(COLUMN_STRIPE_PRICE_ID)
}

func (o *planImplementation) SetStripePriceID(stripePriceID string) PlanInterface {
	o.Set(COLUMN_STRIPE_PRICE_ID, stripePriceID)
	return o
}

func (o *planImplementation) Title() string {
	return o.Get(COLUMN_TITLE)
}

func (o *planImplementation) SetTitle(title string) PlanInterface {
	o.Set(COLUMN_TITLE, title)
	return o
}

func (o *planImplementation) Type() string {
	return o.Get(COLUMN_TYPE)
}

func (o *planImplementation) SetType(type_ string) PlanInterface {
	o.Set(COLUMN_TYPE, type_)
	return o
}

func (o *planImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *planImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_UPDATED_AT), carbon.UTC)
}

func (o *planImplementation) SetUpdatedAt(updatedAt string) PlanInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}
