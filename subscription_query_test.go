package subscriptionstore

import (
	"strings"
	"testing"
)

func TestSubscriptionQuerySettersAndGetters(t *testing.T) {
	query := NewSubscriptionQuery().
		SetID("sub_1").
		SetIDIn([]string{"sub_1", "sub_2"}).
		SetStatus(SUBSCRIPTION_STATUS_ACTIVE).
		SetStatusIn([]string{SUBSCRIPTION_STATUS_ACTIVE, SUBSCRIPTION_STATUS_CANCELLED}).
		SetSubscriberID("subscriber_1").
		SetPlanID("plan_1").
		SetOffset(5).
		SetLimit(10).
		SetOrderBy("created_at").
		SetSortOrder("asc").
		SetCountOnly(true).
		SetSoftDeletedIncluded(true)

	if !query.HasID() || query.ID() != "sub_1" {
		t.Fatalf("expected HasID true with value sub_1")
	}
	if !query.HasIDIn() || len(query.IDIn()) != 2 {
		t.Fatalf("expected HasIDIn true with 2 items")
	}
	if !query.HasStatus() || query.Status() != SUBSCRIPTION_STATUS_ACTIVE {
		t.Fatalf("expected HasStatus true with value %s", SUBSCRIPTION_STATUS_ACTIVE)
	}
	if !query.HasStatusIn() || len(query.StatusIn()) != 2 {
		t.Fatalf("expected HasStatusIn true with 2 items")
	}
	if !query.HasSubscriberID() || query.SubscriberID() != "subscriber_1" {
		t.Fatalf("expected HasSubscriberID true with value subscriber_1")
	}
	if !query.HasPlanID() || query.PlanID() != "plan_1" {
		t.Fatalf("expected HasPlanID true with value plan_1")
	}
	if !query.HasOffset() || query.Offset() != 5 {
		t.Fatalf("expected HasOffset true with value 5")
	}
	if !query.HasLimit() || query.Limit() != 10 {
		t.Fatalf("expected HasLimit true with value 10")
	}
	if !query.HasOrderBy() || query.OrderBy() != "created_at" {
		t.Fatalf("expected HasOrderBy true with value created_at")
	}
	if !query.HasSortOrder() || query.SortOrder() != "asc" {
		t.Fatalf("expected HasSortOrder true with value asc")
	}
	if !query.HasCountOnly() || !query.IsCountOnly() {
		t.Fatalf("expected HasCountOnly true with IsCountOnly true")
	}
	if !query.HasSoftDeletedIncluded() || !query.SoftDeletedIncluded() {
		t.Fatalf("expected HasSoftDeletedIncluded true with SoftDeletedIncluded true")
	}
}

func TestSubscriptionQueryValidateErrors(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func(SubscriptionQueryInterface)
		contains string
	}{
		{
			name: "id empty",
			setup: func(q SubscriptionQueryInterface) {
				q.SetID("")
			},
			contains: "id cannot be empty",
		},
		{
			name: "id_in empty",
			setup: func(q SubscriptionQueryInterface) {
				q.SetIDIn([]string{})
			},
			contains: "id_in cannot be empty array",
		},
		{
			name: "status empty",
			setup: func(q SubscriptionQueryInterface) {
				q.SetStatus("")
			},
			contains: "status cannot be empty",
		},
		{
			name: "status_in empty",
			setup: func(q SubscriptionQueryInterface) {
				q.SetStatusIn([]string{})
			},
			contains: "status_in cannot be empty array",
		},
		{
			name: "subscriber_id empty",
			setup: func(q SubscriptionQueryInterface) {
				q.SetSubscriberID("")
			},
			contains: "subscriber_id cannot be empty",
		},
		{
			name: "plan_id empty",
			setup: func(q SubscriptionQueryInterface) {
				q.SetPlanID("")
			},
			contains: "plan_id cannot be empty",
		},
		{
			name: "limit negative",
			setup: func(q SubscriptionQueryInterface) {
				q.SetLimit(-1)
			},
			contains: "limit cannot be negative",
		},
		{
			name: "offset negative",
			setup: func(q SubscriptionQueryInterface) {
				q.SetOffset(-1)
			},
			contains: "offset cannot be negative",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := NewSubscriptionQuery()
			tc.setup(query)
			err := query.Validate()
			if err == nil {
				t.Fatalf("expected error containing %q", tc.contains)
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Fatalf("unexpected error %q, expected to contain %q", err.Error(), tc.contains)
			}
		})
	}
}

func TestSubscriptionQueryValidateSuccess(t *testing.T) {
	query := NewSubscriptionQuery().
		SetID("sub_1").
		SetLimit(10)

	if err := query.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}
