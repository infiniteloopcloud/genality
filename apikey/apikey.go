package apikey

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/doug-martin/goqu/v9/exp"
	_ "github.com/lib/pq"
)

var driverName = "postgres"
var tableName = "api_key_metrics"
var dialect = goqu.Dialect("postgres")

func NewMetrics(connString string) (Metrics, error) {
	conn, err := sql.Open(driverName, connString)
	if err != nil {
		return Metrics{}, err
	}
	return Metrics{
		db: conn,
	}, nil
}

type Metrics struct {
	db *sql.DB
}

func (m Metrics) Add(ctx context.Context, apiKey string) error {
	return m.add(ctx, apiKey, time.Now().UTC())
}

func (m Metrics) add(ctx context.Context, apiKey string, t time.Time) error {
	query, params, err := dialect.Insert(tableName).Prepared(true).
		Cols("api_key", "time").Vals([]interface{}{apiKey, t}).ToSQL()
	if err != nil {
		return err
	}
	_, err = m.db.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}
	return nil
}

func (m Metrics) GetCountFrom(ctx context.Context, apiKey string, start time.Time) (int, error) {
	query, params, err := dialect.Select(goqu.COUNT("*")).
		From(tableName).Where(
		exp.NewBooleanExpression(exp.EqOp, goqu.L("api_key"), apiKey),
		exp.NewBooleanExpression(exp.GteOp, goqu.L("time"), start),
	).
		Prepared(true).ToSQL()
	if err != nil {
		return 0, err
	}

	row := m.db.QueryRowContext(ctx, query, params...)

	var v int
	if err := row.Scan(&v); err != nil {
		return 0, err
	}
	return v, nil
}

type BucketResponse struct {
	Bucket time.Time `json:"bucket"`
	Count  uint      `json:"count"`
}

func (m Metrics) GetCountBuckets(ctx context.Context, apiKey string, start time.Time, bucketSize time.Duration) ([]BucketResponse, error) {
	query, params, err := dialect.Select(goqu.L(getBucketSize(bucketSize)).As("bucket"), goqu.COUNT("*")).
		From(tableName).Where(
		exp.NewBooleanExpression(exp.EqOp, goqu.L("api_key"), apiKey),
		exp.NewBooleanExpression(exp.GteOp, goqu.L("time"), start),
	).GroupBy("bucket").
		Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := m.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	var res []BucketResponse
	for rows.Next() {
		var r BucketResponse
		if err := rows.Scan(&r.Bucket, &r.Count); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func getBucketSize(d time.Duration) string {
	return fmt.Sprintf(`time_bucket('%d hour',time)`, uint(d.Hours()))
}
