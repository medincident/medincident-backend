-- migrate:up
CREATE INDEX IF NOT EXISTS projections_organizations_created_at_idx
    ON projections.organizations (created_at DESC);

CREATE INDEX IF NOT EXISTS projections_clinics_created_at_idx
    ON projections.clinics (created_at DESC);

CREATE INDEX IF NOT EXISTS projections_departments_created_at_idx
    ON projections.departments (created_at DESC);

CREATE INDEX IF NOT EXISTS projections_incidents_created_at_idx
    ON projections.incidents (created_at DESC);

CREATE INDEX IF NOT EXISTS projections_service_requests_created_at_idx
    ON projections.service_requests (created_at DESC);

-- migrate:down
DROP INDEX IF EXISTS projections_organizations_created_at_idx;
DROP INDEX IF EXISTS projections_clinics_created_at_idx;
DROP INDEX IF EXISTS projections_departments_created_at_idx;
DROP INDEX IF EXISTS projections_incidents_created_at_idx;
DROP INDEX IF EXISTS projections_service_requests_created_at_idx;
