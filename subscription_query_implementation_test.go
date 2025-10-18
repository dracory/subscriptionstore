package subscriptionstore

import (
    "context"
    "strings"
    "testing"

    "github.com/dracory/sb"
)

type mockStore struct{}

func (m mockStore) AutoMigrate(context.Context) error { return nil }
func (m mockStore) EnableDebug(bool)                   {}
func (m mockStore) DatabaseDriverName() string         { return "sqlite3" }
func (m mockStore) PlanCount(context.Context, PlanQueryInterface) (int64, error) {
    return 0, nil
}
func (m mockStore) PlanCreate(context.Context, PlanInterface) error { return nil }
func (m mockStore) PlanDelete(context.Context, PlanInterface) error { return nil }
func (m mockStore) PlanDeleteByID(context.Context, string) error    { return nil }
func (m mockStore) PlanExists(context.Context, string) (bool, error) {
    return false, nil
}
func (m mockStore) PlanFindByID(context.Context, string) (PlanInterface, error) {
    return nil, nil
}
func (m mockStore) PlanList(context.Context, PlanQueryInterface) ([]PlanInterface, error) {
    return nil, nil
}
func (m mockStore) PlanSoftDelete(context.Context, PlanInterface) error { return nil }
func (m mockStore) PlanSoftDeleteByID(context.Context, string) error    { return nil }
func (m mockStore) PlanTableName() string                              { return "plans" }
func (m mockStore) PlanUpdate(context.Context, PlanInterface) error    { return nil }
func (m mockStore) SubscriptionCount(context.Context, SubscriptionQueryInterface) (int64, error) {
    return 0, nil
}
func (m mockStore) SubscriptionCreate(context.Context, SubscriptionInterface) error { return nil }
func (m mockStore) SubscriptionDelete(context.Context, SubscriptionInterface) error { return nil }
func (m mockStore) SubscriptionDeleteByID(context.Context, string) error            { return nil }
func (m mockStore) SubscriptionExists(context.Context, string) (bool, error) {
    return false, nil
}
func (m mockStore) SubscriptionFindByID(context.Context, string) (SubscriptionInterface, error) {
    return nil, nil
}
func (m mockStore) SubscriptionList(context.Context, SubscriptionQueryInterface) ([]SubscriptionInterface, error) {
    return nil, nil
}
func (m mockStore) SubscriptionSoftDelete(context.Context, SubscriptionInterface) error { return nil }
func (m mockStore) SubscriptionSoftDeleteByID(context.Context, string) error            { return nil }
func (m mockStore) SubscriptionTableName() string                                       { return "subscriptions" }
func (m mockStore) SubscriptionUpdate(context.Context, SubscriptionInterface) error     { return nil }

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
        SetOrderDirection(sb.ASC).
        SetCountOnly(true).
        SetWithDeleted(true)

    if !query.HasID() || query.ID() != "sub_1" {
        t.Fatalf("expected HasID true with value sub_1, got id=%s has=%v", query.ID(), query.HasID())
    }
    if !query.HasIDIn() || len(query.IDIn()) != 2 {
        t.Fatalf("expected HasIDIn true with 2 items, got len=%d", len(query.IDIn()))
    }
    if !query.HasStatus() || query.Status() != SUBSCRIPTION_STATUS_ACTIVE {
        t.Fatalf("expected HasStatus true with value %s", SUBSCRIPTION_STATUS_ACTIVE)
    }
    if !query.HasStatusIn() || len(query.StatusIn()) != 2 {
        t.Fatalf("expected HasStatusIn true with 2 items, got len=%d", len(query.StatusIn()))
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
    if !query.HasOrderDirection() || !strings.EqualFold(query.OrderDirection(), sb.ASC) {
        t.Fatalf("expected HasOrderDirection true with value %s", sb.ASC)
    }
    if !query.HasCountOnly() || !query.CountOnly() {
        t.Fatalf("expected HasCountOnly true with CountOnly true")
    }
    if !query.HasWithDeleted() || !query.WithDeleted() {
        t.Fatalf("expected HasWithDeleted true with WithDeleted true")
    }
}

func TestSubscriptionQueryToQueryWithFilters(t *testing.T) {
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
        SetOrderDirection(sb.ASC)

    dataset := query.ToQuery(mockStore{})
    sql, args, err := dataset.Prepared(true).ToSQL()
    if err != nil {
        t.Fatalf("unexpected error generating SQL: %v", err)
    }

    expectedSQL := "SELECT * FROM \"subscriptions\" WHERE ((\"id\" = ?) AND (\"id\" IN (?, ?)) AND (\"status\" = ?) AND (\"status\" IN (?, ?)) AND (\"subscriber_id\" = ?) AND (\"plan_id\" = ?) AND (\"soft_deleted_at\" > ?)) ORDER BY \"created_at\" ASC LIMIT ? OFFSET ?"
    if sql != expectedSQL {
        t.Fatalf("unexpected SQL generated:\nexpected: %s\nactual:   %s", expectedSQL, sql)
    }

    if len(args) != 11 {
        t.Fatalf("expected 11 args, got %d", len(args))
    }
    for i, arg := range args {
        if arg == "" {
            t.Fatalf("argument %d should not be empty", i)
        }
    }

    if args[len(args)-2] == "" || args[len(args)-1] == "" {
        t.Fatalf("limit and offset arguments should not be empty: %v %v", args[len(args)-2], args[len(args)-1])
    }
}

func TestSubscriptionQueryToQueryCountOnlySkipsPagination(t *testing.T) {
    query := NewSubscriptionQuery().
        SetStatus(SUBSCRIPTION_STATUS_ACTIVE).
        SetOffset(10).
        SetLimit(20).
        SetCountOnly(true)

    dataset := query.ToQuery(mockStore{})
    sql, _, err := dataset.Prepared(true).ToSQL()
    if err != nil {
        t.Fatalf("unexpected error generating SQL: %v", err)
    }

    if strings.Contains(sql, "LIMIT") {
        t.Fatalf("LIMIT should not appear in SQL when count only: %s", sql)
    }
    if strings.Contains(sql, "OFFSET") {
        t.Fatalf("OFFSET should not appear in SQL when count only: %s", sql)
    }
}

func TestSubscriptionQueryToQueryWithDeletedSkipsSoftDeleteFilter(t *testing.T) {
    query := NewSubscriptionQuery().
        SetWithDeleted(true)

    dataset := query.ToQuery(mockStore{})
    sql, args, err := dataset.Prepared(true).ToSQL()
    if err != nil {
        t.Fatalf("unexpected error generating SQL: %v", err)
    }

    if strings.Contains(sql, "soft_deleted_at") {
        t.Fatalf("soft_deleted_at filter should be omitted when with deleted: %s", sql)
    }
    if len(args) != 0 {
        t.Fatalf("expected no args when only with deleted set, got %d", len(args))
    }
}

func TestSubscriptionQueryToQueryDefaultOrderDirection(t *testing.T) {
    query := NewSubscriptionQuery().
        SetOrderBy("created_at")

    dataset := query.ToQuery(mockStore{})
    sql, _, err := dataset.Prepared(true).ToSQL()
    if err != nil {
        t.Fatalf("unexpected error generating SQL: %v", err)
    }

    if !strings.Contains(sql, "ORDER BY \"created_at\" DESC") {
        t.Fatalf("expected default order direction DESC, got SQL: %s", sql)
    }
}
