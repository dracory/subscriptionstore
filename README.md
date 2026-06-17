# Subscription Store <a href="https://gitpod.io/#https://github.com/dracory/subscriptionstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/dracory/subscriptionstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/subscriptionstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/subscriptionstore)](https://goreportcard.com/report/github.com/dracory/subscriptionstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/subscriptionstore)](https://pkg.go.dev/github.com/dracory/subscriptionstore)


## Introduction

A powerful, extensible Go module for managing subscription plans and subscriptions, using a clean interface-driven architecture. All data access is performed through interfaces, supporting real database connections (including in-memory SQLite for testing) and allowing for easy extension or customization.

## Features
- **Plan & Subscription Management:** Create, update, query, and delete plans and subscriptions.
- **Interface-Driven:** All entities and stores are accessed via interfaces, enabling mocking, swapping implementations, and extension.
- **Custom Metadata (Metas):** Attach arbitrary key-value data to plans and subscriptions.
- **Flexible Queries:** Use query interfaces to filter, paginate, and sort results.
- **Real DB Testing:** Tests use in-memory SQLite for realistic, high-confidence tests.
- **ORM Powered:** Built on [dracory/neat](https://github.com/dracory/neat) for clean database abstraction.

---

## Architecture Overview
- **Entities:** `Plan` and `Subscription` are private types, always accessed via their public interfaces (`PlanInterface`, `SubscriptionInterface`).
- **Stores:** Data access objects implement the `StoreInterface` and are responsible for all DB operations via neat ORM.
- **Queries:** Query interfaces (`PlanQueryInterface`, `SubscriptionQueryInterface`) allow for composable, type-safe filtering.
- **Extensibility:** Swap or extend any entity, store, or query logic by implementing the relevant interface.

---

## Quick Start Example

### 1. Initialize the Store (SQLite Example)
```go
import (
    "context"
    "database/sql"
    "github.com/dracory/subscriptionstore"
    _ "modernc.org/sqlite"
)

db, _ := sql.Open("sqlite", ":memory:")

store, err := subscriptionstore.NewStore(subscriptionstore.NewStoreOptions{
    DB:                    db,
    PlanTableName:         "plans",
    SubscriptionTableName: "subscriptions",
    AutomigrateEnabled:    true,
})
if err != nil {
    panic(err)
}
```

### 2. Create a New Plan
```go
plan := subscriptionstore.NewPlan().
    SetTitle("Pro Plan").
    SetDescription("Best for teams").
    SetFeatures("Unlimited users, Priority support").
    SetPrice("19.99").
    SetStatus("active").
    SetType("subscription").
    SetInterval(subscriptionstore.PLAN_INTERVAL_MONTHLY).
    SetCurrency("usd")

// Add custom metadata
plan.SetMeta("max_projects", "100")
plan.SetMeta("support_level", "priority")

err := store.PlanCreate(context.Background(), plan)
```

### 3. Create a Subscription and Attach to a Plan
```go
subscription := subscriptionstore.NewSubscription().
    SetSubscriberID("user_123").
    SetPlanID(plan.GetID()).
    SetStatus("active").
    SetPaymentMethodID("pm_abc123")

// Add custom metadata to subscription
subscription.SetMeta("trial", "true")

err = store.SubscriptionCreate(context.Background(), subscription)
```

### 4. Querying Plans and Subscriptions
```go
// List all active plans
query := subscriptionstore.PlanQuery().SetStatus("active")
plans, err := store.PlanList(context.Background(), query)

// Find subscriptions for a user
subQuery := subscriptionstore.SubscriptionQuery().SetSubscriberID("user_123")
subs, err := store.SubscriptionList(context.Background(), subQuery)
```

### 5. Using Metas for Custom Data
```go
// Set a meta value
plan.SetMeta("custom_key", "custom_value")

// Get a meta value
value, _ := plan.Meta("custom_key")

// Check if a meta exists
exists, _ := plan.HasMeta("custom_key")

// Remove a meta value
plan.DeleteMeta("custom_key")
```

---

## Extending the System

Everything in `subscriptionstore` is accessed via interfaces. To extend or customize:
- **Custom Entities:** Implement `PlanInterface` or `SubscriptionInterface` if you want to add new fields or behaviors.
- **Custom Stores:** Implement `StoreInterface` for new data sources or custom logic (e.g., sharding, caching).
- **Custom Queries:** Implement the query interfaces for advanced filtering or external integrations.

This design enables easy swapping of implementations, mocking for tests, or plugging in new backends.

---

## Testing

Tests use a real, in-memory SQLite database. No mocks are used—tests exercise the actual store logic for maximum reliability.

---

## License

This software is subject to both open-source and commercial licensing.

- See [LICENSE](LICENSE) for standard terms.
- See [LICENSE_COMMERCIAL.txt](LICENSE_COMMERCIAL.txt) for commercial licensing terms.

For commercial licensing inquiries, please contact: [https://lesichkov.co.uk/contact](https://lesichkov.co.uk/contact)
