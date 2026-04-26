-- migrate:up
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

-- migrate:down
DROP TABLE IF EXISTS domain.service_requests;
