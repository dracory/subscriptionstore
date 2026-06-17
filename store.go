package subscriptionstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/dracory/neat"
	contractsorm "github.com/dracory/neat/contracts/database/orm"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

// StoreInterface defines the interface for the subscription store.
type StoreInterface interface {
	MigrateDown(ctx context.Context, tx ...*sql.Tx) error
	MigrateUp(ctx context.Context, tx ...*sql.Tx) error
	EnableDebug(debug bool)

	PlanCount(ctx context.Context, query PlanQueryInterface) (int64, error)
	PlanCreate(ctx context.Context, plan PlanInterface) error
	PlanDelete(ctx context.Context, plan PlanInterface) error
	PlanDeleteByID(ctx context.Context, id string) error
	PlanExists(ctx context.Context, planID string) (bool, error)
	PlanFindByID(ctx context.Context, id string) (PlanInterface, error)
	PlanList(ctx context.Context, query PlanQueryInterface) ([]PlanInterface, error)
	PlanSoftDelete(ctx context.Context, plan PlanInterface) error
	PlanSoftDeleteByID(ctx context.Context, id string) error
	PlanTableName() string
	PlanUpdate(ctx context.Context, plan PlanInterface) error

	SubscriptionCount(ctx context.Context, query SubscriptionQueryInterface) (int64, error)
	SubscriptionCreate(ctx context.Context, subscription SubscriptionInterface) error
	SubscriptionDelete(ctx context.Context, subscription SubscriptionInterface) error
	SubscriptionDeleteByID(ctx context.Context, id string) error
	SubscriptionExists(ctx context.Context, subscriptionID string) (bool, error)
	SubscriptionFindByID(ctx context.Context, id string) (SubscriptionInterface, error)
	SubscriptionList(ctx context.Context, query SubscriptionQueryInterface) ([]SubscriptionInterface, error)
	SubscriptionSoftDelete(ctx context.Context, subscription SubscriptionInterface) error
	SubscriptionSoftDeleteByID(ctx context.Context, id string) error
	SubscriptionTableName() string
	SubscriptionUpdate(ctx context.Context, subscription SubscriptionInterface) error
}

var _ StoreInterface = (*storeImplementation)(nil)

// == TYPE =====================================================================

type storeImplementation struct {
	planTableName         string
	subscriptionTableName string
	db                    *neat.Database
	automigrateEnabled    bool
	debugEnabled          bool
	sqlLogger             *slog.Logger
}

// PUBLIC METHODS ==============================================================

// MigrateUp creates the plans and subscriptions tables if they do not exist
func (st *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if st.db.Schema().HasTable(st.planTableName) {
		if st.debugEnabled {
			st.sqlLogger.Info("MigrateUp: plan table already exists", "table", st.planTableName)
		}
	} else {
		err := st.db.Schema().Create(st.planTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 40)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_TYPE, 50)
			table.String(COLUMN_STATUS, 40)
			table.String(COLUMN_TITLE, 100)
			table.Text(COLUMN_DESCRIPTION)
			table.String(COLUMN_INTERVAL, 40)
			table.String(COLUMN_CURRENCY, 40)
			table.String(COLUMN_PRICE, 40)
			table.String(COLUMN_STRIPE_PRICE_ID, 100)
			table.Text(COLUMN_FEATURES)
			table.Text(COLUMN_MEMO)
			table.Text(COLUMN_METAS)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT)
		})
		if err != nil {
			if st.debugEnabled {
				st.sqlLogger.Error("MigrateUp: plan table failed", "error", err)
			}
			return err
		}
	}

	if st.db.Schema().HasTable(st.subscriptionTableName) {
		if st.debugEnabled {
			st.sqlLogger.Info("MigrateUp: subscription table already exists", "table", st.subscriptionTableName)
		}
	} else {
		err := st.db.Schema().Create(st.subscriptionTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 40)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_STATUS, 40)
			table.String(COLUMN_SUBSCRIBER_ID, 50)
			table.String(COLUMN_PLAN_ID, 50)
			table.DateTime(COLUMN_PERIOD_START)
			table.DateTime(COLUMN_PERIOD_END)
			table.String(COLUMN_CANCEL_AT_PERIOD_END, 3)
			table.String(COLUMN_PAYMENT_METHOD_ID, 40)
			table.Text(COLUMN_MEMO)
			table.Text(COLUMN_METAS)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT)
		})
		if err != nil {
			if st.debugEnabled {
				st.sqlLogger.Error("MigrateUp: subscription table failed", "error", err)
			}
			return err
		}
	}

	return nil
}

// MigrateDown drops the plans and subscriptions tables
func (st *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	if st.db.Schema().HasTable(st.planTableName) {
		if err := st.db.Schema().Drop(st.planTableName); err != nil {
			if st.debugEnabled {
				st.sqlLogger.Error("MigrateDown: plan table failed", "error", err)
			}
			return err
		}
	}
	if st.db.Schema().HasTable(st.subscriptionTableName) {
		if err := st.db.Schema().Drop(st.subscriptionTableName); err != nil {
			if st.debugEnabled {
				st.sqlLogger.Error("MigrateDown: subscription table failed", "error", err)
			}
			return err
		}
	}
	return nil
}

// EnableDebug enables the debug option
func (st *storeImplementation) EnableDebug(debug bool) {
	st.debugEnabled = debug
	if debug {
		st.db.EnableDebug()
		st.sqlLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		st.db.DisableDebug()
		st.sqlLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
}

// PlanTableName returns the plan table name
func (st *storeImplementation) PlanTableName() string {
	return st.planTableName
}

// SubscriptionTableName returns the subscription table name
func (st *storeImplementation) SubscriptionTableName() string {
	return st.subscriptionTableName
}

// == PLAN METHODS =============================================================

// PlanCount returns the number of plans based on the given query options
func (st *storeImplementation) PlanCount(ctx context.Context, query PlanQueryInterface) (int64, error) {
	if query == nil {
		return 0, errors.New("plan query: cannot be nil")
	}
	if err := query.Validate(); err != nil {
		return 0, err
	}

	q := st.buildPlanQuery(query)

	var count int64
	err := q.Table(st.planTableName).Count(&count)
	return count, err
}

// PlanCreate creates a new plan
func (st *storeImplementation) PlanCreate(ctx context.Context, plan PlanInterface) error {
	if plan == nil {
		return errors.New("subscriptionstore > plan create. plan cannot be nil")
	}

	if plan.GetCreatedAt() == "" {
		plan.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	}
	if plan.GetUpdatedAt() == "" {
		plan.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	}
	if plan.GetSoftDeletedAt() == "" {
		plan.SetSoftDeletedAt(MAX_DATETIME)
	}

	metasMap, err := plan.GetMetas()
	if err != nil {
		return err
	}
	var metasStr string
	if metasMap != nil {
		b, err := json.Marshal(metasMap)
		if err != nil {
			return err
		}
		metasStr = string(b)
	}

	row := map[string]any{
		COLUMN_ID:              plan.GetID(),
		COLUMN_TYPE:            plan.GetType(),
		COLUMN_STATUS:          plan.GetStatus(),
		COLUMN_TITLE:           plan.GetTitle(),
		COLUMN_DESCRIPTION:     plan.GetDescription(),
		COLUMN_INTERVAL:        plan.GetInterval(),
		COLUMN_CURRENCY:        plan.GetCurrency(),
		COLUMN_PRICE:           plan.GetPrice(),
		COLUMN_STRIPE_PRICE_ID: plan.GetStripePriceID(),
		COLUMN_FEATURES:        plan.GetFeatures(),
		COLUMN_MEMO:            plan.GetMemo(),
		COLUMN_METAS:           metasStr,
		COLUMN_CREATED_AT:      plan.GetCreatedAtCarbon().StdTime(),
		COLUMN_UPDATED_AT:      plan.GetUpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: plan.GetSoftDeletedAtCarbon().StdTime(),
	}

	return st.db.Query().Table(st.planTableName).Create(row)
}

// PlanDelete deletes a plan
func (st *storeImplementation) PlanDelete(ctx context.Context, plan PlanInterface) error {
	if plan == nil {
		return errors.New("plan is nil")
	}
	return st.PlanDeleteByID(ctx, plan.GetID())
}

// PlanDeleteByID deletes a plan by id
func (st *storeImplementation) PlanDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("plan id is empty")
	}
	_, err := st.db.Query().Table(st.planTableName).Where(COLUMN_ID+" = ?", id).Delete()
	return err
}

// PlanExists returns true if a plan exists
func (st *storeImplementation) PlanExists(ctx context.Context, planID string) (bool, error) {
	if planID == "" {
		return false, errors.New("plan id is empty")
	}
	count, err := st.PlanCount(ctx, PlanQuery().SetID(planID))
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// PlanFindByID finds a plan by id
func (st *storeImplementation) PlanFindByID(ctx context.Context, id string) (PlanInterface, error) {
	if id == "" {
		return nil, errors.New("plan id is empty")
	}
	list, err := st.PlanList(ctx, PlanQuery().SetID(id).SetLimit(1))
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// PlanList retrieves a list of plans
func (st *storeImplementation) PlanList(ctx context.Context, query PlanQueryInterface) ([]PlanInterface, error) {
	if query == nil {
		return []PlanInterface{}, errors.New("at plan list > plan query is nil")
	}
	if err := query.Validate(); err != nil {
		return []PlanInterface{}, err
	}

	q := st.buildPlanQuery(query)

	type planRow struct {
		ID            string    `db:"id"`
		Type          string    `db:"type"`
		Status        string    `db:"status"`
		Title         string    `db:"title"`
		Description   string    `db:"description"`
		Interval      string    `db:"interval"`
		Currency      string    `db:"currency"`
		Price         string    `db:"price"`
		StripePriceID string    `db:"stripe_price_id"`
		Features      string    `db:"features"`
		Memo          string    `db:"memo"`
		Metas         string    `db:"metas"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
		SoftDeletedAt time.Time `db:"soft_deleted_at"`
	}

	var rows []planRow
	if err := q.Table(st.planTableName).Get(&rows); err != nil {
		return []PlanInterface{}, err
	}

	list := make([]PlanInterface, 0, len(rows))
	for _, r := range rows {
		p := &planImplementation{}
		p.SetID(r.ID)
		p.SetType(r.Type)
		p.SetStatus(r.Status)
		p.SetTitle(r.Title)
		p.SetDescription(r.Description)
		p.SetInterval(r.Interval)
		p.SetCurrency(r.Currency)
		p.SetPrice(r.Price)
		p.SetStripePriceID(r.StripePriceID)
		p.SetFeatures(r.Features)
		p.SetMemo(r.Memo)
		p.MetasField = r.Metas
		p.CreatedAtField.CreatedAt = r.CreatedAt
		p.UpdatedAtField.UpdatedAt = r.UpdatedAt
		p.SoftDeletesMaxDate.SoftDeletedAt = r.SoftDeletedAt
		list = append(list, p)
	}

	return list, nil
}

// PlanSoftDelete soft deletes a plan
func (st *storeImplementation) PlanSoftDelete(ctx context.Context, plan PlanInterface) error {
	if plan == nil {
		return errors.New("plan is nil")
	}
	plan.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	return st.PlanUpdate(ctx, plan)
}

// PlanSoftDeleteByID soft deletes a plan by id
func (st *storeImplementation) PlanSoftDeleteByID(ctx context.Context, id string) error {
	plan, err := st.PlanFindByID(ctx, id)
	if err != nil {
		return err
	}
	return st.PlanSoftDelete(ctx, plan)
}

// PlanUpdate updates a plan
func (st *storeImplementation) PlanUpdate(ctx context.Context, plan PlanInterface) error {
	if plan == nil {
		return errors.New("subscriptionstore > plan update. plan cannot be nil")
	}

	plan.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	metasMap, err := plan.GetMetas()
	if err != nil {
		return err
	}
	var metasStr string
	if metasMap != nil {
		b, err := json.Marshal(metasMap)
		if err != nil {
			return err
		}
		metasStr = string(b)
	}

	row := map[string]any{
		COLUMN_TYPE:            plan.GetType(),
		COLUMN_STATUS:          plan.GetStatus(),
		COLUMN_TITLE:           plan.GetTitle(),
		COLUMN_DESCRIPTION:     plan.GetDescription(),
		COLUMN_INTERVAL:        plan.GetInterval(),
		COLUMN_CURRENCY:        plan.GetCurrency(),
		COLUMN_PRICE:           plan.GetPrice(),
		COLUMN_STRIPE_PRICE_ID: plan.GetStripePriceID(),
		COLUMN_FEATURES:        plan.GetFeatures(),
		COLUMN_MEMO:            plan.GetMemo(),
		COLUMN_METAS:           metasStr,
		COLUMN_UPDATED_AT:      plan.GetUpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: plan.GetSoftDeletedAtCarbon().StdTime(),
	}

	_, err = st.db.Query().Table(st.planTableName).Where(COLUMN_ID+" = ?", plan.GetID()).Update(row)
	return err
}

// == SUBSCRIPTION METHODS ======================================================

// SubscriptionCount returns the number of subscriptions based on the given query options
func (st *storeImplementation) SubscriptionCount(ctx context.Context, query SubscriptionQueryInterface) (int64, error) {
	if query == nil {
		return 0, errors.New("subscription query: cannot be nil")
	}
	if err := query.Validate(); err != nil {
		return 0, err
	}

	q := st.buildSubscriptionQuery(query)

	var count int64
	err := q.Table(st.subscriptionTableName).Count(&count)
	return count, err
}

// SubscriptionCreate creates a new subscription
func (st *storeImplementation) SubscriptionCreate(ctx context.Context, subscription SubscriptionInterface) error {
	if subscription == nil {
		return errors.New("subscriptionstore > subscription create. subscription cannot be nil")
	}

	if subscription.GetPeriodStart() == "" {
		subscription.SetPeriodStart(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if subscription.GetPeriodEnd() == "" {
		subscription.SetPeriodEnd(MAX_DATETIME)
	}
	if subscription.GetCreatedAt() == "" {
		subscription.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	}
	if subscription.GetUpdatedAt() == "" {
		subscription.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	}
	if subscription.GetSoftDeletedAt() == "" {
		subscription.SetSoftDeletedAt(MAX_DATETIME)
	}

	metasMap, err := subscription.GetMetas()
	if err != nil {
		return err
	}
	var metasStr string
	if metasMap != nil {
		b, err := json.Marshal(metasMap)
		if err != nil {
			return err
		}
		metasStr = string(b)
	}

	row := map[string]any{
		COLUMN_ID:                   subscription.GetID(),
		COLUMN_STATUS:               subscription.GetStatus(),
		COLUMN_SUBSCRIBER_ID:        subscription.GetSubscriberID(),
		COLUMN_PLAN_ID:              subscription.GetPlanID(),
		COLUMN_PERIOD_START:         subscription.GetPeriodStartCarbon().StdTime(),
		COLUMN_PERIOD_END:           subscription.GetPeriodEndCarbon().StdTime(),
		COLUMN_CANCEL_AT_PERIOD_END: lo.Ternary(subscription.GetCancelAtPeriodEnd(), YES, NO),
		COLUMN_PAYMENT_METHOD_ID:    subscription.GetPaymentMethodID(),
		COLUMN_MEMO:                 subscription.GetMemo(),
		COLUMN_METAS:                metasStr,
		COLUMN_CREATED_AT:           subscription.GetCreatedAtCarbon().StdTime(),
		COLUMN_UPDATED_AT:           subscription.GetUpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT:      subscription.GetSoftDeletedAtCarbon().StdTime(),
	}

	return st.db.Query().Table(st.subscriptionTableName).Create(row)
}

// SubscriptionDelete deletes a subscription
func (st *storeImplementation) SubscriptionDelete(ctx context.Context, subscription SubscriptionInterface) error {
	if subscription == nil {
		return errors.New("subscription is nil")
	}
	return st.SubscriptionDeleteByID(ctx, subscription.GetID())
}

// SubscriptionDeleteByID deletes a subscription by id
func (st *storeImplementation) SubscriptionDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("subscription id is empty")
	}
	_, err := st.db.Query().Table(st.subscriptionTableName).Where(COLUMN_ID+" = ?", id).Delete()
	return err
}

// SubscriptionExists returns true if a subscription exists
func (st *storeImplementation) SubscriptionExists(ctx context.Context, subscriptionID string) (bool, error) {
	if subscriptionID == "" {
		return false, errors.New("subscription id is empty")
	}
	count, err := st.SubscriptionCount(ctx, SubscriptionQuery().SetID(subscriptionID))
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SubscriptionFindByID finds a subscription by id
func (st *storeImplementation) SubscriptionFindByID(ctx context.Context, id string) (SubscriptionInterface, error) {
	if id == "" {
		return nil, errors.New("subscription id is empty")
	}
	list, err := st.SubscriptionList(ctx, SubscriptionQuery().SetID(id).SetLimit(1))
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// SubscriptionList retrieves a list of subscriptions
func (st *storeImplementation) SubscriptionList(ctx context.Context, query SubscriptionQueryInterface) ([]SubscriptionInterface, error) {
	if query == nil {
		return []SubscriptionInterface{}, errors.New("at subscription list > subscription query is nil")
	}
	if err := query.Validate(); err != nil {
		return []SubscriptionInterface{}, err
	}

	q := st.buildSubscriptionQuery(query)

	type subscriptionRow struct {
		ID                string    `db:"id"`
		Status            string    `db:"status"`
		SubscriberID      string    `db:"subscriber_id"`
		PlanID            string    `db:"plan_id"`
		PeriodStart       time.Time `db:"period_start"`
		PeriodEnd         time.Time `db:"period_end"`
		CancelAtPeriodEnd string    `db:"cancel_at_period_end"`
		PaymentMethodID   string    `db:"payment_method_id"`
		Memo              string    `db:"memo"`
		Metas             string    `db:"metas"`
		CreatedAt         time.Time `db:"created_at"`
		UpdatedAt         time.Time `db:"updated_at"`
		SoftDeletedAt     time.Time `db:"soft_deleted_at"`
	}

	var rows []subscriptionRow
	if err := q.Table(st.subscriptionTableName).Get(&rows); err != nil {
		return []SubscriptionInterface{}, err
	}

	list := make([]SubscriptionInterface, 0, len(rows))
	for _, r := range rows {
		s := &subscriptionImplementation{}
		s.SetID(r.ID)
		s.SetStatus(r.Status)
		s.SetSubscriberID(r.SubscriberID)
		s.SetPlanID(r.PlanID)
		s.SetPeriodStart(carbon.CreateFromStdTime(r.PeriodStart).ToDateTimeString())
		s.SetPeriodEnd(carbon.CreateFromStdTime(r.PeriodEnd).ToDateTimeString())
		s.SetCancelAtPeriodEnd(r.CancelAtPeriodEnd == YES)
		s.SetPaymentMethodID(r.PaymentMethodID)
		s.SetMemo(r.Memo)
		s.MetasField = r.Metas
		s.CreatedAtField.CreatedAt = r.CreatedAt
		s.UpdatedAtField.UpdatedAt = r.UpdatedAt
		s.SoftDeletesMaxDate.SoftDeletedAt = r.SoftDeletedAt
		list = append(list, s)
	}

	return list, nil
}

// SubscriptionSoftDelete soft deletes a subscription
func (st *storeImplementation) SubscriptionSoftDelete(ctx context.Context, subscription SubscriptionInterface) error {
	if subscription == nil {
		return errors.New("subscription is nil")
	}
	subscription.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	return st.SubscriptionUpdate(ctx, subscription)
}

// SubscriptionSoftDeleteByID soft deletes a subscription by id
func (st *storeImplementation) SubscriptionSoftDeleteByID(ctx context.Context, id string) error {
	subscription, err := st.SubscriptionFindByID(ctx, id)
	if err != nil {
		return err
	}
	return st.SubscriptionSoftDelete(ctx, subscription)
}

// SubscriptionUpdate updates a subscription
func (st *storeImplementation) SubscriptionUpdate(ctx context.Context, subscription SubscriptionInterface) error {
	if subscription == nil {
		return errors.New("subscriptionstore > subscription update. subscription cannot be nil")
	}

	subscription.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	metasMap, err := subscription.GetMetas()
	if err != nil {
		return err
	}
	var metasStr string
	if metasMap != nil {
		b, err := json.Marshal(metasMap)
		if err != nil {
			return err
		}
		metasStr = string(b)
	}

	row := map[string]any{
		COLUMN_STATUS:               subscription.GetStatus(),
		COLUMN_SUBSCRIBER_ID:        subscription.GetSubscriberID(),
		COLUMN_PLAN_ID:              subscription.GetPlanID(),
		COLUMN_PERIOD_START:         subscription.GetPeriodStartCarbon().StdTime(),
		COLUMN_PERIOD_END:           subscription.GetPeriodEndCarbon().StdTime(),
		COLUMN_CANCEL_AT_PERIOD_END: lo.Ternary(subscription.GetCancelAtPeriodEnd(), YES, NO),
		COLUMN_PAYMENT_METHOD_ID:    subscription.GetPaymentMethodID(),
		COLUMN_MEMO:                 subscription.GetMemo(),
		COLUMN_METAS:                metasStr,
		COLUMN_UPDATED_AT:           subscription.GetUpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT:      subscription.GetSoftDeletedAtCarbon().StdTime(),
	}

	_, err = st.db.Query().Table(st.subscriptionTableName).Where(COLUMN_ID+" = ?", subscription.GetID()).Update(row)
	return err
}

// == QUERY BUILDERS ===========================================================

// buildPlanQuery builds a neat query from the plan query interface.
func (st *storeImplementation) buildPlanQuery(query PlanQueryInterface) contractsorm.Query {
	q := st.db.Query()

	if query == nil {
		return q
	}

	if query.HasID() && query.ID() != "" {
		q = q.Where(COLUMN_ID+" = ?", query.ID())
	}
	if query.HasIDIn() && len(query.IDIn()) > 0 {
		args := make([]any, len(query.IDIn()))
		for i, id := range query.IDIn() {
			args[i] = id
		}
		q = q.WhereIn(COLUMN_ID, args)
	}
	if query.HasStatus() && query.Status() != "" {
		q = q.Where(COLUMN_STATUS+" = ?", query.Status())
	}
	if query.HasStatusIn() && len(query.StatusIn()) > 0 {
		args := make([]any, len(query.StatusIn()))
		for i, v := range query.StatusIn() {
			args[i] = v
		}
		q = q.WhereIn(COLUMN_STATUS, args)
	}
	if query.HasInterval() && query.Interval() != "" {
		q = q.Where(COLUMN_INTERVAL+" = ?", query.Interval())
	}
	if query.HasIntervalIn() && len(query.IntervalIn()) > 0 {
		args := make([]any, len(query.IntervalIn()))
		for i, v := range query.IntervalIn() {
			args[i] = v
		}
		q = q.WhereIn(COLUMN_INTERVAL, args)
	}
	if query.HasType() && query.Type() != "" {
		q = q.Where(COLUMN_TYPE+" = ?", query.Type())
	}
	if query.HasLimit() && query.Limit() > 0 {
		q = q.Limit(query.Limit())
	}
	if query.HasOffset() && query.Offset() > 0 {
		q = q.Offset(query.Offset())
	}
	if query.HasOrderBy() && query.OrderBy() != "" {
		sortOrder := "desc"
		if query.HasSortOrder() && query.SortOrder() != "" {
			sortOrder = query.SortOrder()
		}
		q = q.OrderBy(query.OrderBy(), sortOrder)
	}

	if query.HasSoftDeletedIncluded() && query.SoftDeletedIncluded() {
		q = q.WithSoftDeleted()
	} else {
		q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", carbon.Parse(MAX_DATETIME, carbon.UTC).StdTime())
	}

	return q
}

// buildSubscriptionQuery builds a neat query from the subscription query interface.
func (st *storeImplementation) buildSubscriptionQuery(query SubscriptionQueryInterface) contractsorm.Query {
	q := st.db.Query()

	if query == nil {
		return q
	}

	if query.HasID() && query.ID() != "" {
		q = q.Where(COLUMN_ID+" = ?", query.ID())
	}
	if query.HasIDIn() && len(query.IDIn()) > 0 {
		args := make([]any, len(query.IDIn()))
		for i, id := range query.IDIn() {
			args[i] = id
		}
		q = q.WhereIn(COLUMN_ID, args)
	}
	if query.HasStatus() && query.Status() != "" {
		q = q.Where(COLUMN_STATUS+" = ?", query.Status())
	}
	if query.HasStatusIn() && len(query.StatusIn()) > 0 {
		args := make([]any, len(query.StatusIn()))
		for i, v := range query.StatusIn() {
			args[i] = v
		}
		q = q.WhereIn(COLUMN_STATUS, args)
	}
	if query.HasSubscriberID() && query.SubscriberID() != "" {
		q = q.Where(COLUMN_SUBSCRIBER_ID+" = ?", query.SubscriberID())
	}
	if query.HasPlanID() && query.PlanID() != "" {
		q = q.Where(COLUMN_PLAN_ID+" = ?", query.PlanID())
	}
	if query.HasLimit() && query.Limit() > 0 {
		q = q.Limit(query.Limit())
	}
	if query.HasOffset() && query.Offset() > 0 {
		q = q.Offset(query.Offset())
	}
	if query.HasOrderBy() && query.OrderBy() != "" {
		sortOrder := "desc"
		if query.HasSortOrder() && query.SortOrder() != "" {
			sortOrder = query.SortOrder()
		}
		q = q.OrderBy(query.OrderBy(), sortOrder)
	}

	if query.HasSoftDeletedIncluded() && query.SoftDeletedIncluded() {
		q = q.WithSoftDeleted()
	} else {
		q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", carbon.Parse(MAX_DATETIME, carbon.UTC).StdTime())
	}

	return q
}
