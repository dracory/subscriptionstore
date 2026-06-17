package subscriptionstore

import (
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestNewSubscriptionDefaults(t *testing.T) {
	subscription := NewSubscription()

	if subscription.GetID() == "" {
		t.Fatal("ID should not be empty")
	}
	if subscription.GetStatus() != SUBSCRIPTION_STATUS_INACTIVE {
		t.Fatalf("expected status %s, got %s", SUBSCRIPTION_STATUS_INACTIVE, subscription.GetStatus())
	}
	if subscription.GetPlanID() != "" {
		t.Fatalf("expected empty plan ID, got %s", subscription.GetPlanID())
	}
	if subscription.GetSubscriberID() != "" {
		t.Fatalf("expected empty subscriber ID, got %s", subscription.GetSubscriberID())
	}
	if subscription.GetPaymentMethodID() != "" {
		t.Fatalf("expected empty payment method ID, got %s", subscription.GetPaymentMethodID())
	}
	if subscription.GetPeriodStart() != carbon.MaxValue().ToDateTimeString() {
		t.Fatalf("expected period start %s, got %s", carbon.MaxValue().ToDateTimeString(), subscription.GetPeriodStart())
	}
	if subscription.GetPeriodEnd() != carbon.MaxValue().ToDateTimeString() {
		t.Fatalf("expected period end %s, got %s", carbon.MaxValue().ToDateTimeString(), subscription.GetPeriodEnd())
	}
	if subscription.GetCancelAtPeriodEnd() {
		t.Fatal("expected cancel at period end to be false")
	}
	if subscription.GetMemo() != "" {
		t.Fatalf("expected empty memo, got %s", subscription.GetMemo())
	}
	if subscription.GetSoftDeletedAt() != carbon.MaxValue().ToDateTimeString() {
		t.Fatalf("expected soft deleted at %s, got %s", carbon.MaxValue().ToDateTimeString(), subscription.GetSoftDeletedAt())
	}
	if subscription.GetCreatedAt() == "" {
		t.Fatal("created at should not be empty")
	}
	if subscription.GetCreatedAtCarbon() == nil {
		t.Fatal("created at carbon should not be nil")
	}
	if subscription.GetUpdatedAt() == "" {
		t.Fatal("updated at should not be empty")
	}
	if subscription.GetUpdatedAtCarbon() == nil {
		t.Fatal("updated at carbon should not be nil")
	}
	metas, err := subscription.GetMetas()
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

	if !subscription.GetCancelAtPeriodEnd() {
		t.Fatal("expected cancel at period end to be true")
	}
	if subscription.GetPlanID() != "plan_123" {
		t.Fatalf("expected plan id plan_123, got %s", subscription.GetPlanID())
	}
	if subscription.GetSubscriberID() != "subscriber_456" {
		t.Fatalf("expected subscriber id subscriber_456, got %s", subscription.GetSubscriberID())
	}
	if subscription.GetPaymentMethodID() != "pm_789" {
		t.Fatalf("expected payment method id pm_789, got %s", subscription.GetPaymentMethodID())
	}
	if subscription.GetStatus() != SUBSCRIPTION_STATUS_ACTIVE {
		t.Fatalf("expected status %s, got %s", SUBSCRIPTION_STATUS_ACTIVE, subscription.GetStatus())
	}
	if subscription.GetMemo() != "Important notes" {
		t.Fatalf("expected memo Important notes, got %s", subscription.GetMemo())
	}
	if subscription.GetSoftDeletedAt() != "2025-01-01 00:00:00" {
		t.Fatalf("expected soft deleted at 2025-01-01 00:00:00, got %s", subscription.GetSoftDeletedAt())
	}
	if subscription.GetCreatedAt() != "2025-01-02 00:00:00" {
		t.Fatalf("expected created at 2025-01-02 00:00:00, got %s", subscription.GetCreatedAt())
	}
	if subscription.GetUpdatedAt() != "2025-01-03 00:00:00" {
		t.Fatalf("expected updated at 2025-01-03 00:00:00, got %s", subscription.GetUpdatedAt())
	}
	if subscription.GetPeriodStart() != "2025-02-01 00:00:00" {
		t.Fatalf("expected period start 2025-02-01 00:00:00, got %s", subscription.GetPeriodStart())
	}
	if subscription.GetPeriodEnd() != "2025-03-01 00:00:00" {
		t.Fatalf("expected period end 2025-03-01 00:00:00, got %s", subscription.GetPeriodEnd())
	}
}

func TestSubscriptionMetasLifecycle(t *testing.T) {
	subscription := NewSubscription()

	var err error
	subscription, err = subscription.SetMetas(map[string]string{"foo": "bar"})
	if err != nil {
		t.Fatalf("unexpected error setting metas: %v", err)
	}

	metas, err := subscription.GetMetas()
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

	metas, err = subscription.GetMetas()
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

	impl.MetasField = ""

	metas, err := subscription.GetMetas()
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
		COLUMN_ID:                   "sub_1",
		COLUMN_PLAN_ID:              "plan_123",
		COLUMN_SUBSCRIBER_ID:        "sub_456",
		COLUMN_PAYMENT_METHOD_ID:    "pm_789",
		COLUMN_STATUS:               SUBSCRIPTION_STATUS_CANCELLED,
		COLUMN_METAS:                "{\"key\":\"value\"}",
		COLUMN_CREATED_AT:           "2024-01-01 00:00:00",
		COLUMN_UPDATED_AT:           "2024-01-02 00:00:00",
		COLUMN_SOFT_DELETED_AT:      carbon.MaxValue().ToDateTimeString(),
		COLUMN_PERIOD_START:         "2024-02-01 00:00:00",
		COLUMN_PERIOD_END:           "2024-03-01 00:00:00",
		COLUMN_CANCEL_AT_PERIOD_END: YES,
	}

	subscription := NewSubscriptionFromExistingData(data)

	if subscription.GetID() != "sub_1" {
		t.Fatalf("expected id sub_1, got %s", subscription.GetID())
	}
	if subscription.GetPlanID() != "plan_123" {
		t.Fatalf("expected plan id plan_123, got %s", subscription.GetPlanID())
	}
	if subscription.GetStatus() != SUBSCRIPTION_STATUS_CANCELLED {
		t.Fatalf("expected status %s, got %s", SUBSCRIPTION_STATUS_CANCELLED, subscription.GetStatus())
	}
	if !subscription.GetCancelAtPeriodEnd() {
		t.Fatal("expected cancel at period end to be true")
	}
	if subscription.GetPeriodStart() != "2024-02-01 00:00:00" {
		t.Fatalf("expected period start 2024-02-01 00:00:00, got %s", subscription.GetPeriodStart())
	}
	if subscription.GetPeriodEnd() != "2024-03-01 00:00:00" {
		t.Fatalf("expected period end 2024-03-01 00:00:00, got %s", subscription.GetPeriodEnd())
	}
}
