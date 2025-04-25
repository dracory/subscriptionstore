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

// PlanCount returns the number of plans based on the given query options
func (store *storeImplementation) PlanCount(ctx context.Context, query PlanQueryInterface) (int64, error) {
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

func (store *storeImplementation) PlanCreate(ctx context.Context, plan PlanInterface) error {
	plan.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	plan.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := plan.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.planTableName).
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

	plan.MarkAsNotDirty()

	return nil
}

func (store *storeImplementation) PlanDelete(ctx context.Context, plan PlanInterface) error {
	if plan == nil {
		return errors.New("plan is nil")
	}

	return store.PlanDeleteByID(ctx, plan.ID())
}

func (store *storeImplementation) PlanDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("plan id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.planTableName).
		Prepared(true).
		Where(goqu.C("id").Eq(id)).
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

// PlanExists returns true if a plan exists based on the given query options
func (store *storeImplementation) PlanExists(ctx context.Context, planID string) (bool, error) {
	if planID == "" {
		return false, errors.New("plan id is empty")
	}

	count, err := store.PlanCount(ctx, PlanQuery().SetID(planID))

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (store *storeImplementation) PlanFindByID(ctx context.Context, id string) (PlanInterface, error) {
	if id == "" {
		return nil, errors.New("plan id is empty")
	}

	list, err := store.PlanList(ctx, NewPlanQuery().
		SetID(id).
		SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *storeImplementation) PlanList(ctx context.Context, query PlanQueryInterface) ([]PlanInterface, error) {
	sqlStr, _, errSql := query.ToQuery(store).Select().ToSQL()

	if errSql != nil {
		return []PlanInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	modelMaps, err := db.SelectToMapString(sqlStr)
	if err != nil {
		return []PlanInterface{}, err
	}

	list := []PlanInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewPlanFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *storeImplementation) PlanSoftDelete(ctx context.Context, plan PlanInterface) error {
	if plan == nil {
		return errors.New("plan is nil")
	}

	plan.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.PlanUpdate(ctx, plan)
}

func (store *storeImplementation) PlanSoftDeleteByID(ctx context.Context, id string) error {
	plan, err := store.PlanFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.PlanSoftDelete(ctx, plan)
}

func (store *storeImplementation) PlanUpdate(ctx context.Context, plan PlanInterface) error {
	if plan == nil {
		return errors.New("order is nil")
	}

	plan.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := plan.DataChanged()

	delete(dataChanged, "id") // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.planTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C("id").Eq(plan.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

	plan.MarkAsNotDirty()

	return err
}
