package subscriptionstore

import (
    "testing"

    "github.com/dromara/carbon/v2"
)

func TestNewSubscriptionDefaults(t *testing.T) {
    subscription := NewSubscription()

    if subscription.ID() == "" {
        t.Fatal("ID should not be empty")
    }
    if subscription.Status() != SUBSCRIPTION_STATUS_INACTIVE {
        t.Fatalf("expected status %s, got %s", SUBSCRIPTION_STATUS_INACTIVE, subscription.Status())
    }
    if subscription.PlanID() != "" {
        t.Fatalf("expected empty plan ID, got %s", subscription.PlanID())
    }
    if subscription.SubscriberID() != "" {
        t.Fatalf("expected empty subscriber ID, got %s", subscription.SubscriberID())
    }
    if subscription.PaymentMethodID() != "" {
        t.Fatalf("expected empty payment method ID, got %s", subscription.PaymentMethodID())
    }
    if subscription.PeriodStart() != carbon.MaxValue().ToDateTimeString() {
        t.Fatalf("expected period start %s, got %s", carbon.MaxValue().ToDateTimeString(), subscription.PeriodStart())
    }
    if subscription.PeriodEnd() != carbon.MaxValue().ToDateTimeString() {
        t.Fatalf("expected period end %s, got %s", carbon.MaxValue().ToDateTimeString(), subscription.PeriodEnd())
    }
    if subscription.CancelAtPeriodEnd() {
        t.Fatal("expected cancel at period end to be false")
    }
    if subscription.Memo() != "" {
        t.Fatalf("expected empty memo, got %s", subscription.Memo())
    }
    if subscription.SoftDeletedAt() != carbon.MaxValue().ToDateTimeString() {
        t.Fatalf("expected soft deleted at %s, got %s", carbon.MaxValue().ToDateTimeString(), subscription.SoftDeletedAt())
    }
    if subscription.CreatedAt() == "" {
        t.Fatal("created at should not be empty")
    }
    if subscription.CreatedAtCarbon() == nil {
        t.Fatal("created at carbon should not be nil")
    }
    if subscription.UpdatedAt() == "" {
        t.Fatal("updated at should not be empty")
    }
    if subscription.UpdatedAtCarbon() == nil {
        t.Fatal("updated at carbon should not be nil")
    }
    metas, err := subscription.Metas()
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

func TestSubscriptionSettersAndGetters(t *testing.T) {
    subscription := NewSubscription()

    subscription = subscription.SetPlanID("plan_123").
        SetSubscriberID("subscriber_456").
        SetPaymentMethodID("pm_789").
        SetStatus(SUBSCRIPTION_STATUS_ACTIVE).
        SetMemo("Important notes").
        SetSoftDeletedAt("2025-01-01 00:00:00").
        SetCreatedAt("2025-01-02 00:00:00").
        SetUpdatedAt("2025-01-03 00:00:00").
        SetPeriodStart("2025-02-01 00:00:00").
        SetPeriodEnd("2025-03-01 00:00:00")

    subscription = subscription.SetCancelAtPeriodEnd(true)

    if !subscription.CancelAtPeriodEnd() {
        t.Fatal("expected cancel at period end to be true")
    }
    if subscription.PlanID() != "plan_123" {
        t.Fatalf("expected plan id plan_123, got %s", subscription.PlanID())
    }
    if subscription.SubscriberID() != "subscriber_456" {
        t.Fatalf("expected subscriber id subscriber_456, got %s", subscription.SubscriberID())
    }
    if subscription.PaymentMethodID() != "pm_789" {
        t.Fatalf("expected payment method id pm_789, got %s", subscription.PaymentMethodID())
    }
    if subscription.Status() != SUBSCRIPTION_STATUS_ACTIVE {
        t.Fatalf("expected status %s, got %s", SUBSCRIPTION_STATUS_ACTIVE, subscription.Status())
    }
    if subscription.Memo() != "Important notes" {
        t.Fatalf("expected memo Important notes, got %s", subscription.Memo())
    }
    if subscription.SoftDeletedAt() != "2025-01-01 00:00:00" {
        t.Fatalf("expected soft deleted at 2025-01-01 00:00:00, got %s", subscription.SoftDeletedAt())
    }
    if subscription.CreatedAt() != "2025-01-02 00:00:00" {
        t.Fatalf("expected created at 2025-01-02 00:00:00, got %s", subscription.CreatedAt())
    }
    if subscription.UpdatedAt() != "2025-01-03 00:00:00" {
        t.Fatalf("expected updated at 2025-01-03 00:00:00, got %s", subscription.UpdatedAt())
    }
    if subscription.PeriodStart() != "2025-02-01 00:00:00" {
        t.Fatalf("expected period start 2025-02-01 00:00:00, got %s", subscription.PeriodStart())
    }
    if subscription.PeriodEnd() != "2025-03-01 00:00:00" {
        t.Fatalf("expected period end 2025-03-01 00:00:00, got %s", subscription.PeriodEnd())
    }
}

func TestSubscriptionMetasLifecycle(t *testing.T) {
    subscription := NewSubscription()

    var err error
    subscription, err = subscription.SetMetas(map[string]string{"foo": "bar"})
    if err != nil {
        t.Fatalf("unexpected error setting metas: %v", err)
    }

    metas, err := subscription.Metas()
    if err != nil {
        t.Fatalf("unexpected error retrieving metas: %v", err)
    }
    if len(metas) != 1 || metas["foo"] != "bar" {
        t.Fatalf("unexpected metas value: %v", metas)
    }

    has, err := subscription.HasMeta("foo")
    if err != nil {
        t.Fatalf("unexpected error checking meta: %v", err)
    }
    if !has {
        t.Fatal("expected meta foo to exist")
    }

    value, err := subscription.Meta("foo")
    if err != nil {
        t.Fatalf("unexpected error retrieving meta: %v", err)
    }
    if value != "bar" {
        t.Fatalf("expected meta foo value bar, got %s", value)
    }

    subscription, err = subscription.SetMeta("baz", "qux")
    if err != nil {
        t.Fatalf("unexpected error setting meta: %v", err)
    }

    metas, err = subscription.Metas()
    if err != nil {
        t.Fatalf("unexpected error retrieving metas after set: %v", err)
    }
    if len(metas) != 2 || metas["baz"] != "qux" {
        t.Fatalf("unexpected metas after set: %v", metas)
    }

    subscription, err = subscription.DeleteMeta("foo")
    if err != nil {
        t.Fatalf("unexpected error deleting meta: %v", err)
    }

    has, err = subscription.HasMeta("foo")
    if err != nil {
        t.Fatalf("unexpected error checking meta after delete: %v", err)
    }
    if has {
        t.Fatal("expected meta foo to be removed")
    }

    value, err = subscription.Meta("foo")
    if err != nil {
        t.Fatalf("unexpected error retrieving meta after delete: %v", err)
    }
    if value != "" {
        t.Fatalf("expected empty value for removed meta, got %s", value)
    }
}

func TestSubscriptionMetasEmptyHandling(t *testing.T) {
    subscription := NewSubscription()

    impl, ok := subscription.(*subscriptionImplementation)
    if !ok {
        t.Fatal("unexpected type for subscription implementation")
    }

    impl.Set(COLUMN_METAS, "")

    metas, err := subscription.Metas()
    if err != nil {
        t.Fatalf("unexpected error retrieving metas: %v", err)
    }
    if metas != nil {
        t.Fatalf("expected nil metas, got %v", metas)
    }

    has, err := subscription.HasMeta("missing")
    if err != nil {
        t.Fatalf("unexpected error checking missing meta: %v", err)
    }
    if has {
        t.Fatal("expected missing meta to return false")
    }

    value, err := subscription.Meta("missing")
    if err != nil {
        t.Fatalf("unexpected error retrieving missing meta: %v", err)
    }
    if value != "" {
        t.Fatalf("expected empty string for missing meta, got %s", value)
    }

    subscription, err = subscription.DeleteMeta("missing")
    if err != nil {
        t.Fatalf("unexpected error deleting missing meta: %v", err)
    }
    _ = subscription
}

func TestNewSubscriptionFromExistingData(t *testing.T) {
    data := map[string]string{
        COLUMN_ID:              "sub_1",
        COLUMN_PLAN_ID:         "plan_123",
        COLUMN_SUBSCRIBER_ID:   "sub_456",
        COLUMN_PAYMENT_METHOD_ID: "pm_789",
        COLUMN_STATUS:          SUBSCRIPTION_STATUS_CANCELLED,
        COLUMN_METAS:           "{\"key\":\"value\"}",
        COLUMN_CREATED_AT:      "2024-01-01 00:00:00",
        COLUMN_UPDATED_AT:      "2024-01-02 00:00:00",
        COLUMN_SOFT_DELETED_AT: carbon.MaxValue().ToDateTimeString(),
        COLUMN_PERIOD_START:    "2024-02-01 00:00:00",
        COLUMN_PERIOD_END:      "2024-03-01 00:00:00",
        COLUMN_CANCEL_AT_PERIOD_END: YES,
        COLUMN_MEMO:            "stored memo",
    }

    subscription := NewSubscriptionFromExistingData(data)

    if subscription.ID() != "sub_1" {
        t.Fatalf("expected id sub_1, got %s", subscription.ID())
    }
    if subscription.PlanID() != "plan_123" {
        t.Fatalf("expected plan id plan_123, got %s", subscription.PlanID())
    }
    if subscription.SubscriberID() != "sub_456" {
        t.Fatalf("expected subscriber id sub_456, got %s", subscription.SubscriberID())
    }
    if subscription.PaymentMethodID() != "pm_789" {
        t.Fatalf("expected payment method id pm_789, got %s", subscription.PaymentMethodID())
    }
    if subscription.Status() != SUBSCRIPTION_STATUS_CANCELLED {
        t.Fatalf("expected status %s, got %s", SUBSCRIPTION_STATUS_CANCELLED, subscription.Status())
    }
    if subscription.Memo() != "stored memo" {
        t.Fatalf("expected memo stored memo, got %s", subscription.Memo())
    }
    if subscription.SoftDeletedAt() != carbon.MaxValue().ToDateTimeString() {
        t.Fatalf("expected soft deleted at %s, got %s", carbon.MaxValue().ToDateTimeString(), subscription.SoftDeletedAt())
    }
    if subscription.CreatedAt() != "2024-01-01 00:00:00" {
        t.Fatalf("expected created at 2024-01-01 00:00:00, got %s", subscription.CreatedAt())
    }
    if subscription.UpdatedAt() != "2024-01-02 00:00:00" {
        t.Fatalf("expected updated at 2024-01-02 00:00:00, got %s", subscription.UpdatedAt())
    }
    if subscription.PeriodStart() != "2024-02-01 00:00:00" {
        t.Fatalf("expected period start 2024-02-01 00:00:00, got %s", subscription.PeriodStart())
    }
    if subscription.PeriodEnd() != "2024-03-01 00:00:00" {
        t.Fatalf("expected period end 2024-03-01 00:00:00, got %s", subscription.PeriodEnd())
    }
    if !subscription.CancelAtPeriodEnd() {
        t.Fatal("expected cancel at period end to be true")
    }

    metas, err := subscription.Metas()
    if err != nil {
        t.Fatalf("unexpected error retrieving metas: %v", err)
    }
    if len(metas) != 1 || metas["key"] != "value" {
        t.Fatalf("unexpected metas from existing data: %v", metas)
    }
}
