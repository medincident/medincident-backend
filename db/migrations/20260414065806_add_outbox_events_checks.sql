-- migrate:up
ALTER TABLE outbox.events
    ADD CONSTRAINT outbox_events_subject_non_empty CHECK (length(subject) > 0),
    ADD CONSTRAINT outbox_events_payload_non_empty CHECK (octet_length(payload) > 0);

-- migrate:down
ALTER TABLE outbox.events
    DROP CONSTRAINT IF EXISTS outbox_events_payload_non_empty,
    DROP CONSTRAINT IF EXISTS outbox_events_subject_non_empty;
