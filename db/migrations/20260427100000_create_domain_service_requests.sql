-- migrate:up

CREATE TABLE domain.request_types (
    id              UUID PRIMARY KEY,
    organization_id UUID NOT NULL
        REFERENCES domain.organizations(id) ON DELETE RESTRICT,
    name            TEXT NOT NULL,
    description     TEXT,
    is_active       BOOL NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX request_types_org_name_active_idx
    ON domain.request_types (organization_id, name)
    WHERE is_active = true;

CREATE INDEX request_types_organization_id_idx
    ON domain.request_types (organization_id);

CREATE TABLE domain.service_requests (
    id              UUID PRIMARY KEY,
    organization_id UUID NOT NULL
        REFERENCES domain.organizations(id) ON DELETE RESTRICT,
    clinic_id       UUID NOT NULL
        REFERENCES domain.clinics(id) ON DELETE RESTRICT,
    department_id   UUID NOT NULL
        REFERENCES domain.departments(id) ON DELETE RESTRICT,
    type_id         UUID NOT NULL
        REFERENCES domain.request_types(id) ON DELETE RESTRICT,
    incident_id     UUID
        REFERENCES domain.incidents(id) ON DELETE RESTRICT,
    description     TEXT NOT NULL,
    status          domain.request_status NOT NULL DEFAULT 'created',
    author_id       TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX service_requests_organization_id_idx
    ON domain.service_requests (organization_id);
CREATE INDEX service_requests_department_id_idx
    ON domain.service_requests (department_id);
CREATE INDEX service_requests_type_id_idx
    ON domain.service_requests (type_id);
CREATE INDEX service_requests_status_idx
    ON domain.service_requests (status);
CREATE INDEX service_requests_incident_id_idx
    ON domain.service_requests (incident_id)
    WHERE incident_id IS NOT NULL;

CREATE TABLE domain.service_request_executors (
    id              UUID PRIMARY KEY,
    request_id      UUID NOT NULL
        REFERENCES domain.service_requests(id) ON DELETE CASCADE,
    employee_id     UUID NOT NULL
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    assigned_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    assigned_by_id  TEXT NOT NULL
);

CREATE UNIQUE INDEX service_request_executors_request_employee_idx
    ON domain.service_request_executors (request_id, employee_id);

CREATE INDEX service_request_executors_request_id_idx
    ON domain.service_request_executors (request_id);

-- migrate:down
DROP TABLE IF EXISTS domain.service_request_executors;
DROP TABLE IF EXISTS domain.service_requests;
DROP TABLE IF EXISTS domain.request_types;
