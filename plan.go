package subscriptionstore

import (
	"encoding/json"
	"log"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// PlanInterface defines the methods for a Plan entity
type PlanInterface interface {
	IsSoftDeleted() bool

	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) PlanInterface

	GetCurrency() string
	SetCurrency(currency string) PlanInterface

	GetDescription() string
	SetDescription(description string) PlanInterface

	GetFeatures() string
	SetFeatures(features string) PlanInterface

	GetID() string
	SetID(id string) PlanInterface

	GetInterval() string
	SetInterval(interval string) PlanInterface

	GetMemo() string
	SetMemo(memo string) PlanInterface

	GetMetas() (map[string]string, error)
	SetMetas(data map[string]string) (PlanInterface, error)
	HasMeta(key string) (bool, error)
	Meta(key string) (string, error)
	SetMeta(key string, value string) (PlanInterface, error)
	DeleteMeta(key string) (PlanInterface, error)

	GetPrice() string
	GetPriceFloat() float64
	SetPrice(price string) PlanInterface

	GetSoftDeletedAt() string
	GetSoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) PlanInterface

	GetStatus() string
	SetStatus(status string) PlanInterface

	GetStripePriceID() string
	SetStripePriceID(stripePriceID string) PlanInterface

	GetTitle() string
	SetTitle(title string) PlanInterface

	GetType() string
	SetType(type_ string) PlanInterface

	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) PlanInterface
}

var _ PlanInterface = (*planImplementation)(nil)

// == TYPE =====================================================================

type planImplementation struct {
	orm.ShortID

	TypeField          string `db:"type"`
	StatusField        string `db:"status"`
	TitleField         string `db:"title"`
	DescriptionField   string `db:"description"`
	IntervalField      string `db:"interval"`
	CurrencyField      string `db:"currency"`
	PriceField         string `db:"price"`
	StripePriceIDField string `db:"stripe_price_id"`
	FeaturesField      string `db:"features"`
	MemoField          string `db:"memo"`
	MetasField         string `db:"metas"`

	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate
}

// == CONSTRUCTORS =============================================================

func NewPlan() PlanInterface {
	o := &planImplementation{}
	o.SetID(neatuid.GenerateShortID())
	o.SetStatus(PLAN_STATUS_INACTIVE)
	o.SetStripePriceID("")
	o.SetDescription("")
	o.SetFeatures("")
	o.SetMemo("")
	o.SetSoftDeletedAt(MAX_DATETIME)
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	if _, err := o.SetMetas(map[string]string{}); err != nil {
		log.Println(err.Error())
	}
	return o
}

func NewPlanFromExistingData(data map[string]string) PlanInterface {
	o := &planImplementation{}
	o.SetID(data[COLUMN_ID])
	o.SetType(data[COLUMN_TYPE])
	o.SetStatus(data[COLUMN_STATUS])
	o.SetTitle(data[COLUMN_TITLE])
	o.SetDescription(data[COLUMN_DESCRIPTION])
	o.SetInterval(data[COLUMN_INTERVAL])
	o.SetCurrency(data[COLUMN_CURRENCY])
	o.SetPrice(data[COLUMN_PRICE])
	o.SetStripePriceID(data[COLUMN_STRIPE_PRICE_ID])
	o.SetFeatures(data[COLUMN_FEATURES])
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

func (o *planImplementation) IsSoftDeleted() bool {
	return o.SoftDeletesMaxDate.IsSoftDeleted()
}

// == SETTERS AND GETTERS ======================================================

func (o *planImplementation) GetID() string {
	return o.ShortID.ID
}

func (o *planImplementation) SetID(id string) PlanInterface {
	o.ShortID.ID = id
	return o
}

func (o *planImplementation) GetType() string {
	return o.TypeField
}

func (o *planImplementation) SetType(type_ string) PlanInterface {
	o.TypeField = type_
	return o
}

func (o *planImplementation) GetStatus() string {
	return o.StatusField
}

func (o *planImplementation) SetStatus(status string) PlanInterface {
	o.StatusField = status
	return o
}

func (o *planImplementation) GetTitle() string {
	return o.TitleField
}

func (o *planImplementation) SetTitle(title string) PlanInterface {
	o.TitleField = title
	return o
}

func (o *planImplementation) GetDescription() string {
	return o.DescriptionField
}

func (o *planImplementation) SetDescription(description string) PlanInterface {
	o.DescriptionField = description
	return o
}

func (o *planImplementation) GetInterval() string {
	return o.IntervalField
}

func (o *planImplementation) SetInterval(interval string) PlanInterface {
	o.IntervalField = interval
	return o
}

func (o *planImplementation) GetCurrency() string {
	return o.CurrencyField
}

func (o *planImplementation) SetCurrency(currency string) PlanInterface {
	o.CurrencyField = currency
	return o
}

func (o *planImplementation) GetPrice() string {
	return o.PriceField
}

func (o *planImplementation) GetPriceFloat() float64 {
	return cast.ToFloat64(o.PriceField)
}

func (o *planImplementation) SetPrice(price string) PlanInterface {
	o.PriceField = price
	return o
}

func (o *planImplementation) GetStripePriceID() string {
	return o.StripePriceIDField
}

func (o *planImplementation) SetStripePriceID(stripePriceID string) PlanInterface {
	o.StripePriceIDField = stripePriceID
	return o
}

func (o *planImplementation) GetFeatures() string {
	return o.FeaturesField
}

func (o *planImplementation) SetFeatures(features string) PlanInterface {
	o.FeaturesField = features
	return o
}

func (o *planImplementation) GetMemo() string {
	return o.MemoField
}

func (o *planImplementation) SetMemo(memo string) PlanInterface {
	o.MemoField = memo
	return o
}

func (o *planImplementation) GetMetas() (map[string]string, error) {
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

func (o *planImplementation) SetMetas(data map[string]string) (PlanInterface, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	o.MetasField = string(jsonBytes)
	return o, nil
}

func (o *planImplementation) HasMeta(key string) (bool, error) {
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

func (o *planImplementation) Meta(key string) (string, error) {
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

func (o *planImplementation) SetMeta(key string, value string) (PlanInterface, error) {
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

func (o *planImplementation) DeleteMeta(key string) (PlanInterface, error) {
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

func (o *planImplementation) GetCreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

func (o *planImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

func (o *planImplementation) SetCreatedAt(createdAt string) PlanInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

func (o *planImplementation) GetUpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

func (o *planImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

func (o *planImplementation) SetUpdatedAt(updatedAt string) PlanInterface {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

func (o *planImplementation) GetSoftDeletedAt() string {
	if o.SoftDeletesMaxDate.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt).ToDateTimeString()
}

func (o *planImplementation) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt)
}

func (o *planImplementation) SetSoftDeletedAt(deletedAt string) PlanInterface {
	if deletedAt == "" {
		return o
	}
	o.SoftDeletesMaxDate.SoftDeletedAt = carbon.Parse(deletedAt, carbon.UTC).StdTime()
	return o
}
