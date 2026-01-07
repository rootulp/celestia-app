-- Initialize PostgreSQL database with required TimescaleDB extensions
-- This script is automatically executed when the PostgreSQL container starts for the first time

-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Enable TimescaleDB Toolkit extension (provides percentile_agg and other functions)
CREATE EXTENSION IF NOT EXISTS timescaledb_toolkit;
