-- migrate:up
CREATE OR REPLACE FUNCTION outbox.notify_inserted()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
  PERFORM pg_notify('outbox_events', '');
  RETURN NULL;
END;
$$;

CREATE TRIGGER outbox_notify
  AFTER INSERT ON outbox.events
  FOR EACH STATEMENT
  EXECUTE FUNCTION outbox.notify_inserted();

-- migrate:down
DROP TRIGGER IF EXISTS outbox_notify ON outbox.events;
DROP FUNCTION IF EXISTS outbox.notify_inserted();
