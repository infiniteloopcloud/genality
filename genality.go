package genality

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
var tableName = "genality"
var dialect = goqu.Dialect(driverName)

type Opts struct {
	ConnectionString string
}

func New(opts Opts) (Descriptor, error) {
	return newGenality(opts)
}

func newGenality(opts Opts) (genality, error) {
	conn, err := sql.Open(driverName, opts.ConnectionString)
	if err != nil {
		return genality{}, err
	}
	return genality{
		db: conn,
	}, nil
}

type genality struct {
	db *sql.DB
}

func (m genality) Add(ctx context.Context, record string) error {
	return m.add(ctx, record, time.Now().UTC())
}

func (m genality) add(ctx context.Context, record string, t time.Time) error {
	query, params, err := dialect.Insert(tableName).Prepared(true).
		Cols("record", "time").Vals([]interface{}{record, t}).ToSQL()
	if err != nil {
		return err
	}
	_, err = m.db.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}
	return nil
}

func (m genality) GetCountFrom(ctx context.Context, record string, start time.Time) (int, error) {
	query, params, err := dialect.Select(goqu.COUNT("*")).
		From(tableName).Where(
		exp.NewBooleanExpression(exp.EqOp, goqu.L("record"), record),
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

func (m genality) GetCountBuckets(ctx context.Context, record string, start time.Time, bucketSize time.Duration) ([]BucketResponse, error) {
	query, params, err := dialect.Select(goqu.L(getBucketSize(bucketSize)).As("bucket"), goqu.COUNT("*")).
		From(tableName).Where(
		exp.NewBooleanExpression(exp.EqOp, goqu.L("record"), record),
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
