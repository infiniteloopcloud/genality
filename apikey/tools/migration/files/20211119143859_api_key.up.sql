CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

CREATE TABLE IF NOT EXISTS api_key_metrics
(
    api_key VARCHAR(255) NOT NULL,
    time timestamptz
    );

SELECT create_hypertable('api_key_metrics','time',chunk_time_interval => INTERVAL '1 day');
