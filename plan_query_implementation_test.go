package subscriptionstore

import (
    "strings"
    "testing"

    "github.com/dracory/sb"
)

func TestPlanQueryIntervalSetters(t *testing.T) {
    query := NewPlanQuery().
        SetInterval(PLAN_INTERVAL_MONTHLY).
        SetIntervalIn([]string{PLAN_INTERVAL_MONTHLY, PLAN_INTERVAL_YEARLY})

    if !query.HasInterval() {
        t.Fatal("expected HasInterval to be true")
    }
    if query.Interval() != PLAN_INTERVAL_MONTHLY {
        t.Fatalf("expected interval %s, got %s", PLAN_INTERVAL_MONTHLY, query.Interval())
    }
    if !query.HasIntervalIn() {
        t.Fatal("expected HasIntervalIn to be true")
    }
    if len(query.IntervalIn()) != 2 {
        t.Fatalf("expected IntervalIn length 2, got %d", len(query.IntervalIn()))
    }
}

func TestPlanQueryToQueryWithIntervalFilters(t *testing.T) {
    query := NewPlanQuery().
        SetID("plan_1").
        SetIDIn([]string{"plan_1", "plan_2"}).
        SetStatus(PLAN_STATUS_ACTIVE).
        SetStatusIn([]string{PLAN_STATUS_ACTIVE, PLAN_STATUS_INACTIVE}).
        SetInterval(PLAN_INTERVAL_MONTHLY).
        SetIntervalIn([]string{PLAN_INTERVAL_MONTHLY, PLAN_INTERVAL_YEARLY}).
        SetOrderBy("created_at").
        SetOrderDirection(sb.ASC).
        SetLimit(10).
        SetOffset(5)

    dataset := query.ToQuery(mockStore{})
    sql, args, err := dataset.Prepared(true).ToSQL()
    if err != nil {
        t.Fatalf("unexpected error generating SQL: %v", err)
    }

    expectedSQL := "SELECT * FROM \"plans\" WHERE ((\"id\" = ?) AND (\"id\" IN (?, ?)) AND (\"status\" = ?) AND (\"status\" IN (?, ?)) AND (\"interval\" = ?) AND (\"interval\" IN (?, ?)) AND (\"soft_deleted_at\" > ?)) ORDER BY \"created_at\" ASC LIMIT ? OFFSET ?"
    if sql != expectedSQL {
        t.Fatalf("unexpected SQL generated:\nexpected: %s\nactual:   %s", expectedSQL, sql)
    }

    if len(args) != 12 {
        t.Fatalf("expected 12 args, got %d", len(args))
    }
    for i, arg := range args {
        if arg == "" {
            t.Fatalf("argument %d should not be empty", i)
        }
    }

    if !strings.Contains(sql, "\"interval\" = ?") {
        t.Fatalf("expected SQL to contain interval equality filter: %s", sql)
    }
    if !strings.Contains(sql, "\"interval\" IN (?, ?)") {
        t.Fatalf("expected SQL to contain interval IN filter: %s", sql)
    }
}
