-- migrate:up
CREATE TABLE projections.service_requests (
    id              UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    clinic_id       UUID NOT NULL,
    department_id   UUID NOT NULL,
    type_id         UUID NOT NULL,
    incident_id     UUID,
    description     TEXT NOT NULL,
    status          domain.request_status NOT NULL DEFAULT 'created',
    author_id       TEXT NOT NULL,
    author_display_name TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL,
    updated_at      TIMESTAMPTZ NOT NULL
);

CREATE INDEX proj_service_requests_org_status_idx
    ON projections.service_requests (organization_id, status);
CREATE INDEX proj_service_requests_dept_idx
    ON projections.service_requests (department_id);
CREATE INDEX proj_service_requests_incident_idx
    ON projections.service_requests (incident_id)
    WHERE incident_id IS NOT NULL;

-- migrate:down
DROP TABLE IF EXISTS projections.service_requests;
