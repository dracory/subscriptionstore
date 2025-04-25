package subscriptionstore

import (
	"context"
	"testing"

	_ "modernc.org/sqlite"
)

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

	subFound, errFind := store.SubscriptionFindByID(ctx, sub.ID())
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

	subFound, _ := store.SubscriptionFindByID(ctx, sub.ID())
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

	err = store.SubscriptionDeleteByID(ctx, sub.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	subFound, _ := store.SubscriptionFindByID(ctx, sub.ID())
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

	exists, err := store.SubscriptionExists(ctx, sub.ID())
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

	subFound, err := store.SubscriptionFindByID(ctx, sub.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if subFound == nil || subFound.Memo() != "Updated memo" {
		t.Fatal("SubscriptionUpdate did not update memo correctly")
	}
}
