-- migrate:up
CREATE TYPE domain.request_status AS ENUM (
    'created',
    'in_work',
    'on_hold',
    'pending_review',
    'completed',
    'cancelled'
);

-- migrate:down
DROP TYPE IF EXISTS domain.request_status;
