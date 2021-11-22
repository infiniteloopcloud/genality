package migration

import (
	"context"
	"database/sql"
)

var (
	extension = `CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;`
	table     = `CREATE TABLE IF NOT EXISTS genality
	(
	   record VARCHAR(255) NOT NULL,
	   time timestamptz
	);`
	hypertable = `SELECT create_hypertable('genality','time',chunk_time_interval => INTERVAL '1 day', if_not_exists => true);`
)

func Up(ctx context.Context, db *sql.DB) error {
	if err := exec(ctx, db, extension); err != nil {
		return err
	}
	if err := exec(ctx, db, table); err != nil {
		return err
	}
	return exec(ctx, db, hypertable)
}

func Down(ctx context.Context, db *sql.DB) error {
	return exec(ctx, db, "DROP TABLE genality;")
}

func exec(ctx context.Context, db *sql.DB, query string) error {
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
