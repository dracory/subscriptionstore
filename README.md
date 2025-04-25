# Subscription Store

A powerful, extensible Go module for managing subscription plans and subscriptions, using a clean interface-driven architecture. All data access is performed through interfaces, supporting real database connections (including in-memory SQLite for testing) and allowing for easy extension or customization.

## Features
- **Plan & Subscription Management:** Create, update, query, and delete plans and subscriptions.
- **Interface-Driven:** All entities and stores are accessed via interfaces, enabling mocking, swapping implementations, and extension.
- **Custom Metadata (Metas):** Attach arbitrary key-value data to plans and subscriptions.
- **Flexible Queries:** Use query interfaces to filter, paginate, and sort results.
- **Real DB Testing:** Tests use in-memory SQLite for realistic, high-confidence tests.
- **Multi-Database Support:** Powered by [Goqu](https://github.com/doug-martin/goqu), you can use SQLite, MySQL, PostgreSQL, and more—just change the driver and connection string.

---

## Architecture Overview
- **Entities:** `Plan` and `Subscription` are private types, always accessed via their public interfaces (`PlanInterface`, `SubscriptionInterface`).
- **Stores:** Data access objects implement the `StoreInterface` and are responsible for all DB operations, using Goqu for SQL dialect independence.
- **Queries:** Query interfaces (`PlanQueryInterface`, `SubscriptionQueryInterface`) allow for composable, type-safe filtering.
- **Extensibility:** Swap or extend any entity, store, or query logic by implementing the relevant interface.

---

## Quick Start Example

### 1. Initialize the Store (SQLite Example, easily swappable)
```go
import (
    "context"
    "github.com/dracory/subscriptionstore"
    _ "modernc.org/sqlite" // Or use _ "github.com/go-sql-driver/mysql" for MySQL, etc.
    "database/sql"
)

// For SQLite (in-memory):
db, _ := sql.Open("sqlite", ":memory:")
// For MySQL: db, _ := sql.Open("mysql", "user:pass@tcp(localhost:3306)/dbname")
// For PostgreSQL: db, _ := sql.Open("postgres", "host=localhost user=... password=... dbname=... sslmode=disable")

store := subscriptionstore.NewStore(db, "sqlite", "plans", "subscriptions")
store.AutoMigrate(context.Background())
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
    SetPlanID(plan.ID()).
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