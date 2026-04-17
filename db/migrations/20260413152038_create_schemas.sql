-- migrate:up

-- The command service uses a single `domain` schema for aggregate
-- tables (organizations, clinics, departments, etc.). Denormalised
-- read models land in a sibling `projections` schema created by a
-- later migration.
CREATE SCHEMA IF NOT EXISTS domain;

-- migrate:down
DROP SCHEMA IF EXISTS domain CASCADE;
