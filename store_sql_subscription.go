package subscriptionstore

import (
	"errors"

	"github.com/gouniverse/sb"
)

// sqlSubscriptionTableCreate returns a SQL string for creating the country table
func (st *storeImplementation) sqlSubscriptionTableCreate() (string, error) {
	if st.db == nil {
		return "", errors.New("subscription store: db is nil")
	}

	if st.subscriptionTableName == "" {
		return "", errors.New("subscription store: subscription table name is empty")
	}

	sql := sb.NewBuilder(sb.DatabaseDriverName(st.db)).
		Table(st.subscriptionTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     40,
		}).
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_SUBSCRIBER_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 50,
		}).
		Column(sb.Column{
			Name:   COLUMN_PLAN_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 50,
		}).
		Column(sb.Column{
			Name: COLUMN_PERIOD_START,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_PERIOD_END,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name:   COLUMN_CANCEL_AT_PERIOD_END,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 3,
		}).
		Column(sb.Column{
			Name:   COLUMN_PAYMENT_METHOD_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name: COLUMN_MEMO,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_METAS,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_SOFT_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql, nil
}
