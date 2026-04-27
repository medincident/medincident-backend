-- migrate:up

CREATE TABLE projections.request_types (
    id              UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    name            TEXT NOT NULL,
    description     TEXT,
    is_active       BOOL NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL,
    updated_at      TIMESTAMPTZ NOT NULL
);

CREATE INDEX proj_request_types_org_idx
    ON projections.request_types (organization_id);

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

CREATE TABLE projections.service_request_status_history (
    id          UUID PRIMARY KEY,
    request_id  UUID NOT NULL,
    old_status  domain.request_status,
    new_status  domain.request_status NOT NULL,
    actor_id    TEXT NOT NULL,
    actor_name  TEXT NOT NULL,
    changed_at  TIMESTAMPTZ NOT NULL
);

CREATE INDEX proj_sr_status_hist_request_idx
    ON projections.service_request_status_history (request_id, changed_at);

CREATE TABLE projections.service_request_executor_history (
    id            UUID PRIMARY KEY,
    request_id    UUID NOT NULL,
    action        TEXT NOT NULL,
    employee_id   UUID NOT NULL,
    employee_name TEXT NOT NULL,
    actor_id      TEXT NOT NULL,
    actor_name    TEXT NOT NULL,
    changed_at    TIMESTAMPTZ NOT NULL
);

CREATE INDEX proj_sr_exec_hist_request_idx
    ON projections.service_request_executor_history (request_id, changed_at);

-- migrate:down
DROP TABLE IF EXISTS projections.service_request_executor_history;
DROP TABLE IF EXISTS projections.service_request_status_history;
DROP TABLE IF EXISTS projections.service_requests;
DROP TABLE IF EXISTS projections.request_types;
