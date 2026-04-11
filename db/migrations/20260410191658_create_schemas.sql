-- migrate:up

-- The database is split into two schemas with deliberately different
-- responsibilities and access patterns:
--
--   domain — aggregate tables (organizations, future clinics, patients,
--            etc.). Write-heavy, latency-sensitive, driven by the
--            command service's business logic. These are the source of
--            truth for the platform.
--
--   outbox — the transactional outbox staging area (outbox.events) and
--            its supporting objects (NOTIFY trigger, indexes). Pure
--            infrastructure: the command service writes rows here as
--            part of every aggregate write, and a separate publisher
--            process drains them to NATS JetStream.
--
-- Splitting the schemas lets operators grant different privileges to
-- the command-service role vs. the publisher role (publisher only
-- needs SELECT/UPDATE on outbox.events; it never touches domain.*).
-- It also makes backup/replication scoping obvious: only domain.* is
-- authoritative state worth snapshotting long-term.
CREATE SCHEMA IF NOT EXISTS domain;
CREATE SCHEMA IF NOT EXISTS outbox;

-- migrate:down
DROP SCHEMA IF EXISTS outbox CASCADE;
DROP SCHEMA IF EXISTS domain CASCADE;
