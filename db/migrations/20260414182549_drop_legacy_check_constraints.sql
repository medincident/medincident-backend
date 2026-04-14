-- migrate:up

-- AGENTS.md §9: the command service keeps only NOT NULL, PK and FK in
-- the schema; all other invariants live in the application layer.
-- Drop the legacy CHECK constraints that duplicated application-side
-- validators. Note: employee_vacations_no_overlap (EXCLUDE USING GIST)
-- stays — it enforces a multi-row invariant the application cannot
-- express atomically.
ALTER TABLE outbox.events
    DROP CONSTRAINT IF EXISTS outbox_events_subject_non_empty,
    DROP CONSTRAINT IF EXISTS outbox_events_payload_non_empty;

ALTER TABLE domain.employee_vacations
    DROP CONSTRAINT IF EXISTS employee_vacations_end_after_start;

-- migrate:down

ALTER TABLE outbox.events
    ADD CONSTRAINT outbox_events_subject_non_empty CHECK (length(subject) > 0),
    ADD CONSTRAINT outbox_events_payload_non_empty CHECK (octet_length(payload) > 0);

ALTER TABLE domain.employee_vacations
    ADD CONSTRAINT employee_vacations_end_after_start
        CHECK (ends_at IS NULL OR ends_at > starts_at);
