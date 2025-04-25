package subscriptionstore

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/base/database"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
)

func (store *storeImplementation) SubscriptionCreate(ctx context.Context, subscription SubscriptionInterface) error {
	// Set period_start to now if not set
	if subscription.PeriodStart() == "" {
		subscription.SetPeriodStart(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	// Set period_end to max if not set
	if subscription.PeriodEnd() == "" {
		subscription.SetPeriodEnd(carbon.MaxValue().ToDateTimeString(carbon.UTC))
	}

	subscription.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	subscription.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := subscription.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.subscriptionTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

	if err != nil {
		return err
	}

	subscription.MarkAsNotDirty()

	return nil
}

// SubscriptionCount returns the number of subscriptions based on the given query options
func (store *storeImplementation) SubscriptionCount(ctx context.Context, query SubscriptionQueryInterface) (int64, error) {
	query.SetCountOnly(true)

	q := query.ToQuery(store)

	sqlStr, sqlParams, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	store.logSql("count", sqlStr, sqlParams...)

	mapped, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, sqlParams...)

	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]

	i, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err
	}

	return i, nil
}

func (store *storeImplementation) SubscriptionDelete(ctx context.Context, subscription SubscriptionInterface) error {
	if subscription == nil {
		return errors.New("subscription is nil")
	}

	return store.SubscriptionDeleteByID(ctx, subscription.ID())
}

func (store *storeImplementation) SubscriptionDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("subscription id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.subscriptionTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

	return err
}

func (store *storeImplementation) SubscriptionFindByID(ctx context.Context, id string) (SubscriptionInterface, error) {
	if id == "" {
		return nil, errors.New("subscription id is empty")
	}

	list, err := store.SubscriptionList(ctx, SubscriptionQuery().
		SetID(id).SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// SubscriptionExists returns true if a subscription exists based on the given subscription ID
func (store *storeImplementation) SubscriptionExists(ctx context.Context, subscriptionID string) (bool, error) {
	count, err := store.SubscriptionCount(ctx, SubscriptionQuery().SetID(subscriptionID))

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (store *storeImplementation) SubscriptionList(ctx context.Context, query SubscriptionQueryInterface) ([]SubscriptionInterface, error) {
	q := query.ToQuery(store)

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []SubscriptionInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	modelMaps, err := db.SelectToMapString(sqlStr)
	if err != nil {
		return []SubscriptionInterface{}, err
	}

	list := []SubscriptionInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewSubscriptionFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *storeImplementation) SubscriptionSoftDelete(ctx context.Context, subscription SubscriptionInterface) error {
	if subscription == nil {
		return errors.New("subscription is nil")
	}

	subscription.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.SubscriptionUpdate(ctx, subscription)
}

func (store *storeImplementation) SubscriptionSoftDeleteByID(ctx context.Context, id string) error {
	subscription, err := store.SubscriptionFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.SubscriptionSoftDelete(ctx, subscription)
}

func (store *storeImplementation) SubscriptionUpdate(ctx context.Context, subscription SubscriptionInterface) error {
	if subscription == nil {
		return errors.New("order is nil")
	}

	subscription.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := subscription.DataChanged()

	delete(dataChanged, "id") // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.subscriptionTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C("id").Eq(subscription.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

	subscription.MarkAsNotDirty()

	return err
}
