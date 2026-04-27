-- migrate:up

-- Composite types for addresses.
CREATE TYPE domain.geo_point AS (
    longitude DOUBLE PRECISION,
    latitude  DOUBLE PRECISION
);

CREATE TYPE domain.address AS (
    text  TEXT,
    point domain.geo_point
);

-- Incident enums.
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

-- Announcement priority enum.
CREATE TYPE announcement_priority AS ENUM ('normal', 'high');

-- Service request status enum.
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
DROP TYPE IF EXISTS announcement_priority;
DROP TYPE IF EXISTS domain.buffer_status;
DROP TYPE IF EXISTS domain.incident_priority;
DROP TYPE IF EXISTS domain.incident_status;
DROP TYPE IF EXISTS domain.address;
DROP TYPE IF EXISTS domain.geo_point;
