package subscriptionstore

import (
	"math"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestNewPlanDefaults(t *testing.T) {
	plan := NewPlan()

	if plan.GetID() == "" {
		t.Fatal("ID should not be empty")
	}
	if plan.GetStatus() != PLAN_STATUS_INACTIVE {
		t.Fatalf("expected status %s, got %s", PLAN_STATUS_INACTIVE, plan.GetStatus())
	}
	if plan.GetStripePriceID() != "" {
		t.Fatalf("expected empty stripe price id, got %s", plan.GetStripePriceID())
	}
	if plan.GetDescription() != "" {
		t.Fatalf("expected empty description, got %s", plan.GetDescription())
	}
	if plan.GetFeatures() != "" {
		t.Fatalf("expected empty features, got %s", plan.GetFeatures())
	}
	if plan.GetMemo() != "" {
		t.Fatalf("expected empty memo, got %s", plan.GetMemo())
	}
	if plan.GetPrice() != "" {
		t.Fatalf("expected empty price, got %s", plan.GetPrice())
	}
	if plan.GetPriceFloat() != 0 {
		t.Fatalf("expected price float 0, got %f", plan.GetPriceFloat())
	}
	if plan.GetSoftDeletedAt() != carbon.MaxValue().ToDateTimeString() {
		t.Fatalf("expected soft deleted at %s, got %s", carbon.MaxValue().ToDateTimeString(), plan.GetSoftDeletedAt())
	}
	if plan.GetCreatedAt() == "" {
		t.Fatal("created at should not be empty")
	}
	if plan.GetCreatedAtCarbon() == nil {
		t.Fatal("created at carbon should not be nil")
	}
	if plan.GetUpdatedAt() == "" {
		t.Fatal("updated at should not be empty")
	}
	if plan.GetUpdatedAtCarbon() == nil {
		t.Fatal("updated at carbon should not be nil")
	}

	metas, err := plan.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}
	if metas == nil {
		t.Fatal("metas should not be nil")
	}
	if len(metas) != 0 {
		t.Fatalf("expected empty metas, got %v", metas)
	}
}

func TestPlanSettersAndGetters(t *testing.T) {
	plan := NewPlan()

	plan = plan.SetTitle("Example Title").
		SetDescription("Example Description").
		SetFeatures("Feature A").
		SetPrice("9.99").
		SetStatus(PLAN_STATUS_ACTIVE).
		SetType(PLAN_TYPE_GOLD).
		SetInterval(PLAN_INTERVAL_YEARLY).
		SetCurrency(CURRENCY_EUR).
		SetMemo("Remember this").
		SetStripePriceID("price_123").
		SetSoftDeletedAt("2025-01-01 00:00:00").
		SetCreatedAt("2025-01-02 00:00:00").
		SetUpdatedAt("2025-01-03 00:00:00")

	if plan.GetTitle() != "Example Title" {
		t.Fatalf("expected title Example Title, got %s", plan.GetTitle())
	}
	if plan.GetDescription() != "Example Description" {
		t.Fatalf("expected description Example Description, got %s", plan.GetDescription())
	}
	if plan.GetFeatures() != "Feature A" {
		t.Fatalf("expected features Feature A, got %s", plan.GetFeatures())
	}
	if plan.GetPrice() != "9.99" {
		t.Fatalf("expected price 9.99, got %s", plan.GetPrice())
	}
	if math.Abs(plan.GetPriceFloat()-9.99) > 1e-9 {
		t.Fatalf("expected price float 9.99, got %f", plan.GetPriceFloat())
	}
	if plan.GetStatus() != PLAN_STATUS_ACTIVE {
		t.Fatalf("expected status %s, got %s", PLAN_STATUS_ACTIVE, plan.GetStatus())
	}
	if plan.GetType() != PLAN_TYPE_GOLD {
		t.Fatalf("expected type %s, got %s", PLAN_TYPE_GOLD, plan.GetType())
	}
	if plan.GetInterval() != PLAN_INTERVAL_YEARLY {
		t.Fatalf("expected interval %s, got %s", PLAN_INTERVAL_YEARLY, plan.GetInterval())
	}
	if plan.GetCurrency() != CURRENCY_EUR {
		t.Fatalf("expected currency %s, got %s", CURRENCY_EUR, plan.GetCurrency())
	}
	if plan.GetMemo() != "Remember this" {
		t.Fatalf("expected memo Remember this, got %s", plan.GetMemo())
	}
	if plan.GetStripePriceID() != "price_123" {
		t.Fatalf("expected stripe price id price_123, got %s", plan.GetStripePriceID())
	}
	if plan.GetSoftDeletedAt() != "2025-01-01 00:00:00" {
		t.Fatalf("expected soft deleted at 2025-01-01 00:00:00, got %s", plan.GetSoftDeletedAt())
	}
	if plan.GetCreatedAt() != "2025-01-02 00:00:00" {
		t.Fatalf("expected created at 2025-01-02 00:00:00, got %s", plan.GetCreatedAt())
	}
	if plan.GetUpdatedAt() != "2025-01-03 00:00:00" {
		t.Fatalf("expected updated at 2025-01-03 00:00:00, got %s", plan.GetUpdatedAt())
	}
}

func TestPlanMetasLifecycle(t *testing.T) {
	plan := NewPlan()

	var err error
	plan, err = plan.SetMetas(map[string]string{"foo": "bar"})
	if err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	metas, err := plan.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}
	if len(metas) != 1 || metas["foo"] != "bar" {
		t.Fatalf("unexpected metas value: %v", metas)
	}

	has, err := plan.HasMeta("foo")
	if err != nil {
		t.Fatalf("unexpected error checking meta: %v", err)
	}
	if !has {
		t.Fatal("expected meta foo to exist")
	}

	value, err := plan.Meta("foo")
	if err != nil {
		t.Fatalf("unexpected error retrieving meta: %v", err)
	}
	if value != "bar" {
		t.Fatalf("expected meta foo value bar, got %s", value)
	}

	plan, err = plan.SetMeta("baz", "qux")
	if err != nil {
		t.Fatalf("unexpected error setting meta: %v", err)
	}

	metas, err = plan.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas after set: %v", err)
	}
	if len(metas) != 2 || metas["baz"] != "qux" {
		t.Fatalf("unexpected metas after set: %v", metas)
	}

	plan, err = plan.DeleteMeta("foo")
	if err != nil {
		t.Fatalf("unexpected error deleting meta: %v", err)
	}

	has, err = plan.HasMeta("foo")
	if err != nil {
		t.Fatalf("unexpected error checking meta after delete: %v", err)
	}
	if has {
		t.Fatal("expected meta foo to be removed")
	}

	value, err = plan.Meta("foo")
	if err != nil {
		t.Fatalf("unexpected error retrieving meta after delete: %v", err)
	}
	if value != "" {
		t.Fatalf("expected empty value for removed meta, got %s", value)
	}
}

func TestPlanMetasEmptyHandling(t *testing.T) {
	plan := NewPlan()

	impl, ok := plan.(*planImplementation)
	if !ok {
		t.Fatal("unexpected type for plan implementation")
	}

	impl.MetasField = ""

	metas, err := plan.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}
	if metas != nil {
		t.Fatalf("expected nil metas, got %v", metas)
	}

	has, err := plan.HasMeta("missing")
	if err != nil {
		t.Fatalf("unexpected error checking missing meta: %v", err)
	}
	if has {
		t.Fatal("expected missing meta to return false")
	}

	value, err := plan.Meta("missing")
	if err != nil {
		t.Fatalf("unexpected error retrieving missing meta: %v", err)
	}
	if value != "" {
		t.Fatalf("expected empty string for missing meta, got %s", value)
	}

	plan, err = plan.DeleteMeta("missing")
	if err != nil {
		t.Fatalf("unexpected error deleting missing meta: %v", err)
	}
	_ = plan
}

func TestNewPlanFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:              "plan_1",
		COLUMN_TITLE:           "Stored Title",
		COLUMN_STATUS:          PLAN_STATUS_ACTIVE,
		COLUMN_PRICE:           "12.34",
		COLUMN_METAS:           "{\"key\":\"value\"}",
		COLUMN_CREATED_AT:      "2024-01-01 00:00:00",
		COLUMN_UPDATED_AT:      "2024-01-02 00:00:00",
		COLUMN_SOFT_DELETED_AT: carbon.MaxValue().ToDateTimeString(),
	}

	plan := NewPlanFromExistingData(data)

	if plan.GetID() != "plan_1" {
		t.Fatalf("expected id plan_1, got %s", plan.GetID())
	}
	if plan.GetTitle() != "Stored Title" {
		t.Fatalf("expected title Stored Title, got %s", plan.GetTitle())
	}
	if plan.GetStatus() != PLAN_STATUS_ACTIVE {
		t.Fatalf("expected status %s, got %s", PLAN_STATUS_ACTIVE, plan.GetStatus())
	}
	if plan.GetPrice() != "12.34" {
		t.Fatalf("expected price 12.34, got %s", plan.GetPrice())
	}

	metas, err := plan.GetMetas()
	if err != nil {
		t.Fatalf("unexpected error retrieving metas: %v", err)
	}
	if metas == nil || metas["key"] != "value" {
		t.Fatalf("expected metas with key=value, got %v", metas)
	}
}
