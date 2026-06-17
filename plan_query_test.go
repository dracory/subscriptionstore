package subscriptionstore

import (
	"strings"
	"testing"
)

func TestPlanQuerySettersAndGetters(t *testing.T) {
	query := NewPlanQuery().
		SetID("plan_1").
		SetIDIn([]string{"plan_1", "plan_2"}).
		SetStatus(PLAN_STATUS_ACTIVE).
		SetStatusIn([]string{PLAN_STATUS_ACTIVE, PLAN_STATUS_INACTIVE}).
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetIntervalIn([]string{PLAN_INTERVAL_MONTHLY, PLAN_INTERVAL_YEARLY}).
		SetType(PLAN_TYPE_GOLD).
		SetOffset(5).
		SetLimit(10).
		SetOrderBy("created_at").
		SetSortOrder("asc").
		SetCountOnly(true).
		SetSoftDeletedIncluded(true)

	if !query.HasID() || query.ID() != "plan_1" {
		t.Fatalf("expected HasID true with value plan_1")
	}
	if !query.HasIDIn() || len(query.IDIn()) != 2 {
		t.Fatalf("expected HasIDIn true with 2 items")
	}
	if !query.HasStatus() || query.Status() != PLAN_STATUS_ACTIVE {
		t.Fatalf("expected HasStatus true with value %s", PLAN_STATUS_ACTIVE)
	}
	if !query.HasStatusIn() || len(query.StatusIn()) != 2 {
		t.Fatalf("expected HasStatusIn true with 2 items")
	}
	if !query.HasInterval() || query.Interval() != PLAN_INTERVAL_MONTHLY {
		t.Fatalf("expected HasInterval true with value %s", PLAN_INTERVAL_MONTHLY)
	}
	if !query.HasIntervalIn() || len(query.IntervalIn()) != 2 {
		t.Fatalf("expected HasIntervalIn true with 2 items")
	}
	if !query.HasType() || query.Type() != PLAN_TYPE_GOLD {
		t.Fatalf("expected HasType true with value %s", PLAN_TYPE_GOLD)
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

func TestPlanQueryValidateErrors(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func(PlanQueryInterface)
		contains string
	}{
		{
			name: "id empty",
			setup: func(q PlanQueryInterface) {
				q.SetID("")
			},
			contains: "id cannot be empty",
		},
		{
			name: "id_in empty",
			setup: func(q PlanQueryInterface) {
				q.SetIDIn([]string{})
			},
			contains: "id_in cannot be empty array",
		},
		{
			name: "status empty",
			setup: func(q PlanQueryInterface) {
				q.SetStatus("")
			},
			contains: "status cannot be empty",
		},
		{
			name: "status_in empty",
			setup: func(q PlanQueryInterface) {
				q.SetStatusIn([]string{})
			},
			contains: "status_in cannot be empty array",
		},
		{
			name: "interval empty",
			setup: func(q PlanQueryInterface) {
				q.SetInterval("")
			},
			contains: "interval cannot be empty",
		},
		{
			name: "interval_in empty",
			setup: func(q PlanQueryInterface) {
				q.SetIntervalIn([]string{})
			},
			contains: "interval_in cannot be empty array",
		},
		{
			name: "type empty",
			setup: func(q PlanQueryInterface) {
				q.SetType("")
			},
			contains: "type cannot be empty",
		},
		{
			name: "limit negative",
			setup: func(q PlanQueryInterface) {
				q.SetLimit(-1)
			},
			contains: "limit cannot be negative",
		},
		{
			name: "offset negative",
			setup: func(q PlanQueryInterface) {
				q.SetOffset(-1)
			},
			contains: "offset cannot be negative",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := NewPlanQuery()
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

func TestPlanQueryValidateSuccess(t *testing.T) {
	query := NewPlanQuery().
		SetID("plan_1").
		SetLimit(10)

	if err := query.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}
