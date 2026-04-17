-- migrate:up

-- The projections schema holds the read-side denormalized tables that
-- the query-server reads from. The command-server writes into both
-- domain.* and projections.* atomically via sync projector functions
-- inside a single transaction (per the monorepo CQRS merge spec).
--
-- This migration only creates the schema; the table migrations follow
-- as separate files. Plan 1 of the merge ships an unused projections
-- schema so Plan 2 can wire the projector functions on top of an
-- already-existing schema.
CREATE SCHEMA IF NOT EXISTS projections;

-- migrate:down
DROP SCHEMA IF EXISTS projections CASCADE;
