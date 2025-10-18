package subscriptionstore

import (
    "math"
    "testing"

    "github.com/dromara/carbon/v2"
)

func TestNewPlanDefaults(t *testing.T) {
    plan := NewPlan()

    if plan.ID() == "" {
        t.Fatal("ID should not be empty")
    }
    if plan.Status() != PLAN_STATUS_INACTIVE {
        t.Fatalf("expected status %s, got %s", PLAN_STATUS_INACTIVE, plan.Status())
    }
    if plan.StripePriceID() != "" {
        t.Fatalf("expected empty stripe price id, got %s", plan.StripePriceID())
    }
    if plan.Description() != "" {
        t.Fatalf("expected empty description, got %s", plan.Description())
    }
    if plan.Features() != "" {
        t.Fatalf("expected empty features, got %s", plan.Features())
    }
    if plan.Memo() != "" {
        t.Fatalf("expected empty memo, got %s", plan.Memo())
    }
    if plan.Price() != "" {
        t.Fatalf("expected empty price, got %s", plan.Price())
    }
    if plan.PriceFloat() != 0 {
        t.Fatalf("expected price float 0, got %f", plan.PriceFloat())
    }
    if plan.SoftDeletedAt() != carbon.MaxValue().ToDateTimeString() {
        t.Fatalf("expected soft deleted at %s, got %s", carbon.MaxValue().ToDateTimeString(), plan.SoftDeletedAt())
    }
    if plan.CreatedAt() == "" {
        t.Fatal("created at should not be empty")
    }
    if plan.CreatedAtCarbon() == nil {
        t.Fatal("created at carbon should not be nil")
    }
    if plan.UpdatedAt() == "" {
        t.Fatal("updated at should not be empty")
    }
    if plan.UpdatedAtCarbon() == nil {
        t.Fatal("updated at carbon should not be nil")
    }

    metas, err := plan.Metas()
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

    if plan.Title() != "Example Title" {
        t.Fatalf("expected title Example Title, got %s", plan.Title())
    }
    if plan.Description() != "Example Description" {
        t.Fatalf("expected description Example Description, got %s", plan.Description())
    }
    if plan.Features() != "Feature A" {
        t.Fatalf("expected features Feature A, got %s", plan.Features())
    }
    if plan.Price() != "9.99" {
        t.Fatalf("expected price 9.99, got %s", plan.Price())
    }
    if math.Abs(plan.PriceFloat()-9.99) > 1e-9 {
        t.Fatalf("expected price float 9.99, got %f", plan.PriceFloat())
    }
    if plan.Status() != PLAN_STATUS_ACTIVE {
        t.Fatalf("expected status %s, got %s", PLAN_STATUS_ACTIVE, plan.Status())
    }
    if plan.Type() != PLAN_TYPE_GOLD {
        t.Fatalf("expected type %s, got %s", PLAN_TYPE_GOLD, plan.Type())
    }
    if plan.Interval() != PLAN_INTERVAL_YEARLY {
        t.Fatalf("expected interval %s, got %s", PLAN_INTERVAL_YEARLY, plan.Interval())
    }
    if plan.Currency() != CURRENCY_EUR {
        t.Fatalf("expected currency %s, got %s", CURRENCY_EUR, plan.Currency())
    }
    if plan.Memo() != "Remember this" {
        t.Fatalf("expected memo Remember this, got %s", plan.Memo())
    }
    if plan.StripePriceID() != "price_123" {
        t.Fatalf("expected stripe price id price_123, got %s", plan.StripePriceID())
    }
    if plan.SoftDeletedAt() != "2025-01-01 00:00:00" {
        t.Fatalf("expected soft deleted at 2025-01-01 00:00:00, got %s", plan.SoftDeletedAt())
    }
    if plan.CreatedAt() != "2025-01-02 00:00:00" {
        t.Fatalf("expected created at 2025-01-02 00:00:00, got %s", plan.CreatedAt())
    }
    if plan.UpdatedAt() != "2025-01-03 00:00:00" {
        t.Fatalf("expected updated at 2025-01-03 00:00:00, got %s", plan.UpdatedAt())
    }
}

func TestPlanMetasLifecycle(t *testing.T) {
    plan := NewPlan()

    var err error
    plan, err = plan.SetMetas(map[string]string{"foo": "bar"})
    if err != nil {
        t.Fatalf("unexpected error setting metas: %v", err)
    }

    metas, err := plan.Metas()
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

    metas, err = plan.Metas()
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

    impl.Set(COLUMN_METAS, "")

    metas, err := plan.Metas()
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

    if plan.ID() != "plan_1" {
        t.Fatalf("expected id plan_1, got %s", plan.ID())
    }
    if plan.Title() != "Stored Title" {
        t.Fatalf("expected title Stored Title, got %s", plan.Title())
    }
    if plan.Status() != PLAN_STATUS_ACTIVE {
        t.Fatalf("expected status %s, got %s", PLAN_STATUS_ACTIVE, plan.Status())
    }
    if plan.Price() != "12.34" {
        t.Fatalf("expected price 12.34, got %s", plan.Price())
    }
    if math.Abs(plan.PriceFloat()-12.34) > 1e-9 {
        t.Fatalf("expected price float 12.34, got %f", plan.PriceFloat())
    }
    if plan.CreatedAt() != "2024-01-01 00:00:00" {
        t.Fatalf("expected created at 2024-01-01 00:00:00, got %s", plan.CreatedAt())
    }
    if plan.UpdatedAt() != "2024-01-02 00:00:00" {
        t.Fatalf("expected updated at 2024-01-02 00:00:00, got %s", plan.UpdatedAt())
    }
    if plan.SoftDeletedAt() != carbon.MaxValue().ToDateTimeString() {
        t.Fatalf("expected soft deleted at %s, got %s", carbon.MaxValue().ToDateTimeString(), plan.SoftDeletedAt())
    }

    metas, err := plan.Metas()
    if err != nil {
        t.Fatalf("unexpected error retrieving metas: %v", err)
    }
    if len(metas) != 1 || metas["key"] != "value" {
        t.Fatalf("unexpected metas from existing data: %v", metas)
    }
}
