package subscriptionstore

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func initDB(filepath string) *sql.DB {
	if filepath != ":memory:" {
		if err := os.Remove(filepath); err != nil && !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	}
	dsn := filepath + "?parseTime=true&loc=UTC&_loc=UTC"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}
	return db
}

func initStore() (StoreInterface, error) {
	db := initDB(":memory:")
	store, err := NewStore(NewStoreOptions{
		DB:                    db,
		PlanTableName:         "plan_table",
		SubscriptionTableName: "subscription_table",
		AutomigrateEnabled:    true,
	})
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errors.New("unexpected nil store")
	}
	return store, nil
}

// == PLAN TESTS ===============================================================

func TestStorePlanCreate(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	plan := NewPlan().
		SetTitle("Pro Plan").
		SetDescription("Best for teams").
		SetFeatures("Unlimited users, Priority support").
		SetPrice("19.99").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	ctx := context.Background()
	err = store.PlanCreate(ctx, plan)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStorePlanFindByID(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	plan := NewPlan().
		SetTitle("Starter Plan").
		SetDescription("For individuals").
		SetFeatures("Single user").
		SetPrice("4.99").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	ctx := context.Background()
	err = store.PlanCreate(ctx, plan)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	planFound, errFind := store.PlanFindByID(ctx, plan.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if planFound == nil {
		t.Fatal("Plan MUST NOT be nil")
	}
	if planFound.GetTitle() != plan.GetTitle() {
		t.Errorf("expected Title %s, got %s", plan.GetTitle(), planFound.GetTitle())
	}
	if planFound.GetPriceFloat() != plan.GetPriceFloat() {
		t.Errorf("expected Price %s, got %s", plan.GetPrice(), planFound.GetPrice())
	}
}

func TestStorePlanDelete(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	plan := NewPlan().
		SetTitle("Delete Me").
		SetDescription("For deletion test").
		SetFeatures("Feature X").
		SetPrice("1.99").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	ctx := context.Background()
	err = store.PlanCreate(ctx, plan)
	if err != nil {
		t.Fatal("unexpected error creating plan:", err)
	}

	planFound, errFind := store.PlanFindByID(ctx, plan.GetID())
	if errFind != nil {
		t.Fatal("unexpected error finding plan:", errFind)
	}
	if planFound == nil {
		t.Fatal("Plan should exist")
	}

	err = store.PlanDelete(ctx, plan)
	if err != nil {
		t.Fatal("unexpected error deleting plan:", err)
	}

	planDeleted, errFindDeleted := store.PlanFindByID(ctx, plan.GetID())
	if errFindDeleted != nil {
		t.Fatal("unexpected error finding deleted plan:", errFindDeleted)
	}
	if planDeleted != nil {
		t.Fatal("Plan should be deleted")
	}
}

func TestStorePlanDeleteByID(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	plan := NewPlan().
		SetTitle("Delete By ID").
		SetDescription("Delete by ID test").
		SetFeatures("Feature Y").
		SetPrice("2.99").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	ctx := context.Background()
	err = store.PlanCreate(ctx, plan)
	if err != nil {
		t.Fatal("unexpected error creating plan:", err)
	}

	exists, errExists := store.PlanExists(ctx, plan.GetID())
	if errExists != nil {
		t.Fatal("unexpected error checking if plan exists:", errExists)
	}
	if !exists {
		t.Fatal("Plan should exist")
	}

	err = store.PlanDeleteByID(ctx, plan.GetID())
	if err != nil {
		t.Fatal("unexpected error deleting plan by ID:", err)
	}

	existsAfterDelete, errExistsAfterDelete := store.PlanExists(ctx, plan.GetID())
	if errExistsAfterDelete != nil {
		t.Fatal("unexpected error checking if plan exists after delete:", errExistsAfterDelete)
	}
	if existsAfterDelete {
		t.Fatal("Plan should be deleted")
	}
}

func TestStorePlanExists(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	plan1 := NewPlan().
		SetTitle("Exists Plan 1").
		SetDescription("Plan 1").
		SetFeatures("A").
		SetPrice("1.00").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	plan2 := NewPlan().
		SetTitle("Exists Plan 2").
		SetDescription("Plan 2").
		SetFeatures("B").
		SetPrice("2.00").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	err = store.PlanCreate(ctx, plan1)
	if err != nil {
		t.Fatal("unexpected error creating plan1:", err)
	}
	err = store.PlanCreate(ctx, plan2)
	if err != nil {
		t.Fatal("unexpected error creating plan2:", err)
	}

	exists, errExists := store.PlanExists(ctx, plan1.GetID())
	if errExists != nil {
		t.Fatal("unexpected error checking if plan exists by ID:", errExists)
	}
	if !exists {
		t.Fatal("Plan should exist by ID")
	}

	exists, errExists = store.PlanExists(ctx, "non-existent-id")
	if errExists != nil {
		t.Fatal("unexpected error checking if non-existent plan exists:", errExists)
	}
	if exists {
		t.Fatal("Plan should not exist with non-existent ID")
	}
}

func TestStorePlanCount(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	plan1 := NewPlan().
		SetTitle("Count Plan 1").
		SetDescription("Plan 1").
		SetFeatures("A").
		SetPrice("1.00").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	plan2 := NewPlan().
		SetTitle("Count Plan 2").
		SetDescription("Plan 2").
		SetFeatures("B").
		SetPrice("2.00").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	plan3 := NewPlan().
		SetTitle("Count Plan 3").
		SetDescription("Plan 3").
		SetFeatures("C").
		SetPrice("3.00").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	err = store.PlanCreate(ctx, plan1)
	if err != nil {
		t.Fatal("unexpected error creating plan1:", err)
	}
	err = store.PlanCreate(ctx, plan2)
	if err != nil {
		t.Fatal("unexpected error creating plan2:", err)
	}
	err = store.PlanCreate(ctx, plan3)
	if err != nil {
		t.Fatal("unexpected error creating plan3:", err)
	}

	count, errCount := store.PlanCount(ctx, PlanQuery())
	if errCount != nil {
		t.Fatal("unexpected error counting all plans:", errCount)
	}
	if count != 3 {
		t.Fatal("Count should be 3, got:", count)
	}

	count, errCount = store.PlanCount(ctx, PlanQuery().SetID(plan1.GetID()))
	if errCount != nil {
		t.Fatal("unexpected error counting plans by ID:", errCount)
	}
	if count != 1 {
		t.Fatal("Count should be 1 for ID filter, got:", count)
	}

	count, errCount = store.PlanCount(ctx, PlanQuery().SetIDIn([]string{plan1.GetID(), plan3.GetID()}))
	if errCount != nil {
		t.Fatal("unexpected error counting plans by multiple IDs:", errCount)
	}
	if count != 2 {
		t.Fatal("Count should be 2 for multiple IDs filter, got:", count)
	}
}

func TestStorePlanList(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	plan1 := NewPlan().
		SetTitle("List Plan 1").
		SetDescription("Plan 1").
		SetFeatures("A").
		SetPrice("1.00").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	plan2 := NewPlan().
		SetTitle("List Plan 2").
		SetDescription("Plan 2").
		SetFeatures("B").
		SetPrice("2.00").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_YEARLY).
		SetCurrency("usd")

	plan3 := NewPlan().
		SetTitle("List Plan 3").
		SetDescription("Plan 3").
		SetFeatures("C").
		SetPrice("3.00").
		SetStatus("inactive").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	err = store.PlanCreate(ctx, plan1)
	if err != nil {
		t.Fatal("unexpected error creating plan1:", err)
	}
	err = store.PlanCreate(ctx, plan2)
	if err != nil {
		t.Fatal("unexpected error creating plan2:", err)
	}
	err = store.PlanCreate(ctx, plan3)
	if err != nil {
		t.Fatal("unexpected error creating plan3:", err)
	}

	plans, errList := store.PlanList(ctx, PlanQuery().SetLimit(10))
	if errList != nil {
		t.Fatal("unexpected error listing plans:", errList)
	}
	if len(plans) != 3 {
		t.Fatal("Expected 3 plans, got:", len(plans))
	}

	plans, errList = store.PlanList(ctx, PlanQuery().SetStatus("active").SetLimit(10))
	if errList != nil {
		t.Fatal("unexpected error listing active plans:", errList)
	}
	if len(plans) != 2 {
		t.Fatal("Expected 2 active plans, got:", len(plans))
	}
}

func TestStorePlanSoftDelete(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	plan := NewPlan().
		SetTitle("Soft Delete Plan").
		SetDescription("Soft delete test").
		SetFeatures("Feature Z").
		SetPrice("5.99").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	ctx := context.Background()
	err = store.PlanCreate(ctx, plan)
	if err != nil {
		t.Fatal("unexpected error creating plan:", err)
	}

	err = store.PlanSoftDeleteByID(ctx, plan.GetID())
	if err != nil {
		t.Fatal("unexpected error soft deleting plan:", err)
	}

	planFound, errFind := store.PlanFindByID(ctx, plan.GetID())
	if errFind != nil {
		t.Fatal("unexpected error finding plan after soft delete:", errFind)
	}
	if planFound != nil {
		t.Fatal("Plan should be soft deleted and not found")
	}

	plansWithDeleted, errList := store.PlanList(ctx, PlanQuery().
		SetID(plan.GetID()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))
	if errList != nil {
		t.Fatal("unexpected error listing soft deleted plans:", errList)
	}
	if len(plansWithDeleted) != 1 {
		t.Fatal("Expected 1 soft deleted plan, got:", len(plansWithDeleted))
	}
	if !plansWithDeleted[0].IsSoftDeleted() {
		t.Fatal("Plan should be marked as soft deleted")
	}
}

func TestStorePlanUpdate(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	plan := NewPlan().
		SetTitle("Update Plan").
		SetDescription("Update test").
		SetFeatures("Feature U").
		SetPrice("6.99").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	ctx := context.Background()
	err = store.PlanCreate(ctx, plan)
	if err != nil {
		t.Fatal("unexpected error creating plan:", err)
	}

	plan.SetTitle("Updated Title")
	err = store.PlanUpdate(ctx, plan)
	if err != nil {
		t.Fatal("unexpected error updating plan:", err)
	}

	planFound, errFind := store.PlanFindByID(ctx, plan.GetID())
	if errFind != nil {
		t.Fatal("unexpected error finding updated plan:", errFind)
	}
	if planFound == nil || planFound.GetTitle() != "Updated Title" {
		t.Fatal("PlanUpdate did not update title correctly")
	}
}

// == SUBSCRIPTION TESTS ========================================================

func TestStoreSubscriptionCreate(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	sub := NewSubscription().
		SetStatus("active").
		SetSubscriberID("user123").
		SetPlanID("plan123").
		SetMemo("Test subscription")

	ctx := context.Background()
	err = store.SubscriptionCreate(ctx, sub)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreSubscriptionFindByID(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	sub := NewSubscription().
		SetStatus(SUBSCRIPTION_STATUS_ACTIVE).
		SetSubscriberID("user321").
		SetPlanID("plan321").
		SetMemo("FindByID test")

	ctx := context.Background()
	err = store.SubscriptionCreate(ctx, sub)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	subFound, errFind := store.SubscriptionFindByID(ctx, sub.GetID())
	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}
	if subFound == nil {
		t.Fatal("Subscription MUST NOT be nil")
	}
}

func TestStoreSubscriptionDelete(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	sub := NewSubscription().
		SetStatus(SUBSCRIPTION_STATUS_ACTIVE).
		SetSubscriberID("userDel").
		SetPlanID("planDel").
		SetMemo("Delete test")

	ctx := context.Background()
	err = store.SubscriptionCreate(ctx, sub)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SubscriptionDelete(ctx, sub)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	subFound, _ := store.SubscriptionFindByID(ctx, sub.GetID())
	if subFound != nil {
		t.Fatal("Subscription should be deleted")
	}
}

func TestStoreSubscriptionDeleteByID(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	sub := NewSubscription().
		SetStatus(SUBSCRIPTION_STATUS_ACTIVE).
		SetSubscriberID("userDelID").
		SetPlanID("planDelID").
		SetMemo("DeleteByID test")

	ctx := context.Background()
	err = store.SubscriptionCreate(ctx, sub)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SubscriptionDeleteByID(ctx, sub.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	subFound, _ := store.SubscriptionFindByID(ctx, sub.GetID())
	if subFound != nil {
		t.Fatal("Subscription should be deleted by ID")
	}
}

func TestStoreSubscriptionExists(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	sub := NewSubscription().
		SetStatus(SUBSCRIPTION_STATUS_ACTIVE).
		SetSubscriberID("userExists").
		SetPlanID("planExists").
		SetMemo("Exists test")

	ctx := context.Background()
	err = store.SubscriptionCreate(ctx, sub)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	exists, err := store.SubscriptionExists(ctx, sub.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if !exists {
		t.Fatal("Subscription should exist")
	}
}

func TestStoreSubscriptionCount(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	countBefore, err := store.SubscriptionCount(ctx, SubscriptionQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	sub := NewSubscription().
		SetStatus(SUBSCRIPTION_STATUS_ACTIVE).
		SetSubscriberID("userCount").
		SetPlanID("planCount").
		SetMemo("Count test")

	err = store.SubscriptionCreate(ctx, sub)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	countAfter, err := store.SubscriptionCount(ctx, SubscriptionQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if countAfter != countBefore+1 {
		t.Fatalf("SubscriptionCount should increase by 1, before=%d after=%d", countBefore, countAfter)
	}
}

func TestStoreSubscriptionList(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	_ = store.SubscriptionCreate(ctx, NewSubscription().SetStatus("list").SetSubscriberID("userList").SetPlanID("planList").SetMemo("List test 1"))
	_ = store.SubscriptionCreate(ctx, NewSubscription().SetStatus("list").SetSubscriberID("userList2").SetPlanID("planList2").SetMemo("List test 2"))

	list, err := store.SubscriptionList(ctx, SubscriptionQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) < 2 {
		t.Fatal("SubscriptionList should return at least 2 items")
	}
}

func TestStoreSubscriptionSoftDelete(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	sub := NewSubscription().
		SetStatus(SUBSCRIPTION_STATUS_ACTIVE).
		SetSubscriberID("userSoftDel").
		SetPlanID("planSoftDel").
		SetMemo("Soft delete test")

	ctx := context.Background()
	err = store.SubscriptionCreate(ctx, sub)
	if err != nil {
		t.Fatal("unexpected error creating subscription:", err)
	}

	err = store.SubscriptionSoftDeleteByID(ctx, sub.GetID())
	if err != nil {
		t.Fatal("unexpected error soft deleting subscription:", err)
	}

	subFound, errFind := store.SubscriptionFindByID(ctx, sub.GetID())
	if errFind != nil {
		t.Fatal("unexpected error finding subscription after soft delete:", errFind)
	}
	if subFound != nil {
		t.Fatal("Subscription should be soft deleted and not found")
	}

	subsWithDeleted, errList := store.SubscriptionList(ctx, SubscriptionQuery().
		SetID(sub.GetID()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))
	if errList != nil {
		t.Fatal("unexpected error listing soft deleted subscriptions:", errList)
	}
	if len(subsWithDeleted) != 1 {
		t.Fatal("Expected 1 soft deleted subscription, got:", len(subsWithDeleted))
	}
	if !subsWithDeleted[0].IsSoftDeleted() {
		t.Fatal("Subscription should be marked as soft deleted")
	}
}

func TestStoreSubscriptionUpdate(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	sub := NewSubscription().
		SetStatus(SUBSCRIPTION_STATUS_ACTIVE).
		SetSubscriberID("userUpdate").
		SetPlanID("planUpdate").
		SetMemo("Update test")

	ctx := context.Background()
	err = store.SubscriptionCreate(ctx, sub)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	sub.SetMemo("Updated memo")
	err = store.SubscriptionUpdate(ctx, sub)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	subFound, err := store.SubscriptionFindByID(ctx, sub.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if subFound == nil || subFound.GetMemo() != "Updated memo" {
		t.Fatal("SubscriptionUpdate did not update memo correctly")
	}
}

func TestStoreSubscriptionPeriodDefaults(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	sub := NewSubscription().
		SetStatus(SUBSCRIPTION_STATUS_ACTIVE).
		SetSubscriberID("userPeriod").
		SetPlanID("planPeriod").
		SetMemo("Period defaults test")

	ctx := context.Background()
	err = store.SubscriptionCreate(ctx, sub)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	subFound, err := store.SubscriptionFindByID(ctx, sub.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if subFound == nil {
		t.Fatal("Subscription MUST NOT be nil")
	}

	if subFound.GetPeriodStart() == "" {
		t.Fatal("PeriodStart should be set automatically")
	}
	if subFound.GetPeriodEnd() == "" {
		t.Fatal("PeriodEnd should be set automatically")
	}
}

func TestStoreSubscriptionCancelAtPeriodEnd(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	sub := NewSubscription().
		SetStatus(SUBSCRIPTION_STATUS_ACTIVE).
		SetSubscriberID("userCancel").
		SetPlanID("planCancel").
		SetMemo("Cancel test").
		SetCancelAtPeriodEnd(true)

	ctx := context.Background()
	err = store.SubscriptionCreate(ctx, sub)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	subFound, err := store.SubscriptionFindByID(ctx, sub.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if subFound == nil {
		t.Fatal("Subscription MUST NOT be nil")
	}
	if !subFound.GetCancelAtPeriodEnd() {
		t.Fatal("CancelAtPeriodEnd should be true")
	}
}
