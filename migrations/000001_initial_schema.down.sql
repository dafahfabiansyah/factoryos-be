-- Down migration for FactoryOS V1 initial schema

-- Drop triggers
DROP TRIGGER IF EXISTS update_alert_rules_updated_at ON alert_rules;
DROP TRIGGER IF EXISTS update_sensors_updated_at ON sensors;
DROP TRIGGER IF EXISTS update_machines_updated_at ON machines;
DROP TRIGGER IF EXISTS update_production_lines_updated_at ON production_lines;
DROP TRIGGER IF EXISTS update_factories_updated_at ON factories;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_organizations_updated_at ON organizations;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order (respecting foreign keys)
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS alerts;
DROP TABLE IF EXISTS alert_rules;
DROP TABLE IF EXISTS telemetry;
DROP TABLE IF EXISTS sensors;
DROP TABLE IF EXISTS machines;
DROP TABLE IF EXISTS production_lines;
DROP TABLE IF EXISTS factories;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS organizations;

-- Drop UUID extension (optional, only if not used elsewhere)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
