-- migrate:up

-- The database is split into two schemas with deliberately different
-- responsibilities and access patterns:
--
--   domain — aggregate tables (organizations, clinics, departments).
--            Write path of the command service. Source of truth for
--            organisational structure.
--
--   outbox — the transactional outbox staging area (outbox.events)
--            and its supporting objects (NOTIFY trigger, indexes).
--            Pure infrastructure: command-service appends rows here
--            atomically with aggregate writes; a separate publisher
--            process drains them to NATS JetStream.
--
-- Splitting the schemas lets operators grant different privileges to
-- command-service and publisher roles.
CREATE SCHEMA IF NOT EXISTS domain;
CREATE SCHEMA IF NOT EXISTS outbox;

-- migrate:down
DROP SCHEMA IF EXISTS outbox CASCADE;
DROP SCHEMA IF EXISTS domain CASCADE;
