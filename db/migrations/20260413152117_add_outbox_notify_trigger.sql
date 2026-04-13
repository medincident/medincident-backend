-- migrate:up

-- notify_new_event fires a NOTIFY on the 'outbox_event' channel whenever
-- a row is inserted into outbox.events. The payload is the row's BIGINT
-- identity in text form; publisher listeners use it as a monotonic
-- cursor for their SELECT ... WHERE id > last_seen scan. NOTIFY is
-- delivered at transaction commit, so listeners only see events after
-- the aggregate write that produced them has already committed —
-- preserving transactional outbox atomicity.
CREATE OR REPLACE FUNCTION outbox.notify_new_event() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('outbox_event', NEW.id::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER outbox_events_notify
AFTER INSERT ON outbox.events
FOR EACH ROW
EXECUTE FUNCTION outbox.notify_new_event();

-- migrate:down
DROP TRIGGER IF EXISTS outbox_events_notify ON outbox.events;
DROP FUNCTION IF EXISTS outbox.notify_new_event();
