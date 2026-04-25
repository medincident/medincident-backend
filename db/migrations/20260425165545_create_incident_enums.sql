-- migrate:up
CREATE TYPE domain.incident_status AS ENUM (
    'pending',
    'in_progress',
    'done',
    'rejected',
    'cancelled'
);

CREATE TYPE domain.incident_priority AS ENUM (
    'low',
    'normal',
    'high',
    'critical'
);

CREATE TYPE domain.buffer_status AS ENUM (
    'pending',
    'published',
    'rejected',
    'cancelled'
);

-- migrate:down
DROP TYPE IF EXISTS domain.buffer_status;
DROP TYPE IF EXISTS domain.incident_priority;
DROP TYPE IF EXISTS domain.incident_status;
