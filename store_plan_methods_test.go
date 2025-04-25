package subscriptionstore

import (
	"context"
	"testing"

	_ "modernc.org/sqlite"
)

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

	planFound, errFind := store.PlanFindByID(ctx, plan.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if planFound == nil {
		t.Fatal("Plan MUST NOT be nil")
	}
	if planFound.Title() != plan.Title() {
		t.Errorf("expected Title %s, got %s", plan.Title(), planFound.Title())
	}
	if planFound.PriceFloat() != plan.PriceFloat() {
		t.Errorf("expected Price %s, got %s", plan.Price(), planFound.Price())
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

	// Create the plan first
	err = store.PlanCreate(ctx, plan)
	if err != nil {
		t.Fatal("unexpected error creating plan:", err)
	}

	// Verify it exists
	planFound, errFind := store.PlanFindByID(ctx, plan.ID())
	if errFind != nil {
		t.Fatal("unexpected error finding plan:", errFind)
	}
	if planFound == nil {
		t.Fatal("Plan should exist")
	}

	// Delete the plan
	err = store.PlanDelete(ctx, plan)
	if err != nil {
		t.Fatal("unexpected error deleting plan:", err)
	}

	// Verify it was deleted
	planDeleted, errFindDeleted := store.PlanFindByID(ctx, plan.ID())
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

	// Create the plan first
	err = store.PlanCreate(ctx, plan)
	if err != nil {
		t.Fatal("unexpected error creating plan:", err)
	}

	// Verify it exists
	exists, errExists := store.PlanExists(ctx, plan.ID())
	if errExists != nil {
		t.Fatal("unexpected error checking if plan exists:", errExists)
	}
	if !exists {
		t.Fatal("Plan should exist")
	}

	// Delete the plan by ID
	err = store.PlanDeleteByID(ctx, plan.ID())
	if err != nil {
		t.Fatal("unexpected error deleting plan by ID:", err)
	}

	// Verify it was deleted
	existsAfterDelete, errExistsAfterDelete := store.PlanExists(ctx, plan.ID())
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

	// Create plans
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

	// Test exists by ID
	exists, errExists := store.PlanExists(ctx, plan1.ID())
	if errExists != nil {
		t.Fatal("unexpected error checking if plan exists by ID:", errExists)
	}
	if !exists {
		t.Fatal("Plan should exist by ID")
	}

	// Test not exists with non-existent ID
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

	// Create plans
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

	// Test count all
	count, errCount := store.PlanCount(ctx, PlanQuery())
	if errCount != nil {
		t.Fatal("unexpected error counting all plans:", errCount)
	}
	if count != 3 {
		t.Fatal("Count should be 3, got:", count)
	}

	// Test count with ID filter (single result)
	count, errCount = store.PlanCount(ctx, PlanQuery().SetID(plan1.ID()))
	if errCount != nil {
		t.Fatal("unexpected error counting plans by ID:", errCount)
	}
	if count != 1 {
		t.Fatal("Count should be 1 for ID filter, got:", count)
	}

	// Test count with multiple IDs
	count, errCount = store.PlanCount(ctx, PlanQuery().SetIDIn([]string{plan1.ID(), plan3.ID()}))
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

	// Create plans
	plan1 := NewPlan().
		SetTitle("List Plan 1").
		SetDescription("Plan 1").
		SetFeatures("A").
		SetPrice("1.00").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	err = store.PlanCreate(ctx, plan1)
	if err != nil {
		t.Fatal("unexpected error creating plan1:", err)
	}

	plan2 := NewPlan().
		SetTitle("List Plan 2").
		SetDescription("Plan 2").
		SetFeatures("B").
		SetPrice("2.00").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	err = store.PlanCreate(ctx, plan2)
	if err != nil {
		t.Fatal("unexpected error creating plan2:", err)
	}

	plan3 := NewPlan().
		SetTitle("List Plan 3").
		SetDescription("Plan 3").
		SetFeatures("C").
		SetPrice("3.00").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	err = store.PlanCreate(ctx, plan3)
	if err != nil {
		t.Fatal("unexpected error creating plan3:", err)
	}

	// Test list all
	plans, errList := store.PlanList(ctx, PlanQuery())
	if errList != nil {
		t.Fatal("unexpected error listing all plans:", errList)
	}
	if len(plans) != 3 {
		t.Fatal("Should list 3 plans, got:", len(plans))
	}
}

func TestStorePlanUpdate(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create a plan
	plan := NewPlan().
		SetTitle("Update Plan").
		SetDescription("Plan to update").
		SetFeatures("Old Features").
		SetPrice("5.00").
		SetStatus("active").
		SetType("subscription").
		SetInterval(PLAN_INTERVAL_MONTHLY).
		SetCurrency("usd")

	err = store.PlanCreate(ctx, plan)
	if err != nil {
		t.Fatal("unexpected error creating plan:", err)
	}

	// Update the plan values
	plan.SetTitle("Updated Plan").
		SetDescription("Updated description").
		SetFeatures("New Features").
		SetPrice("9.99").
		SetStatus("inactive")

	// Save the updates
	err = store.PlanUpdate(ctx, plan)
	if err != nil {
		t.Fatal("unexpected error updating plan:", err)
	}

	// Retrieve the updated plan
	updatedPlan, errFind := store.PlanFindByID(ctx, plan.ID())
	if errFind != nil {
		t.Fatal("unexpected error finding updated plan:", errFind)
	}

	if updatedPlan == nil {
		t.Fatal("Updated plan should exist")
	}

	// Check if the values were updated correctly
	if updatedPlan.Title() != "Updated Plan" {
		t.Fatal("Plan title should be updated to 'Updated Plan', got:", updatedPlan.Title())
	}
	if updatedPlan.Description() != "Updated description" {
		t.Fatal("Plan description should be updated to 'Updated description', got:", updatedPlan.Description())
	}
	if updatedPlan.Features() != "New Features" {
		t.Fatal("Plan features should be updated to 'New Features', got:", updatedPlan.Features())
	}
	if updatedPlan.PriceFloat() != 9.99 {
		t.Fatal("Plan price should be updated to '9.99', got:", updatedPlan.Price())
	}
	if updatedPlan.Status() != "inactive" {
		t.Fatal("Plan status should be updated to 'inactive', got:", updatedPlan.Status())
	}

	// Check that the ID remains unchanged
	if updatedPlan.ID() != plan.ID() {
		t.Fatal("Plan ID should remain unchanged")
	}
}
