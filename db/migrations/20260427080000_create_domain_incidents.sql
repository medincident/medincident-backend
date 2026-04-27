-- migrate:up

CREATE TABLE domain.incidents (
    id                              UUID PRIMARY KEY,
    organization_id                 UUID NOT NULL
        REFERENCES domain.organizations(id) ON DELETE RESTRICT,
    clinic_id                       UUID NOT NULL
        REFERENCES domain.clinics(id) ON DELETE RESTRICT,
    department_id                   UUID NOT NULL
        REFERENCES domain.departments(id) ON DELETE RESTRICT,
    category_id                     UUID NOT NULL
        REFERENCES domain.incident_categories(id) ON DELETE RESTRICT,
    type_id                         UUID NOT NULL
        REFERENCES domain.incident_types(id) ON DELETE RESTRICT,
    status                          domain.incident_status NOT NULL DEFAULT 'pending',
    priority                        domain.incident_priority NOT NULL DEFAULT 'normal',
    description                     TEXT,
    patient_original_description    TEXT,
    occurred_at                     TIMESTAMPTZ NOT NULL,
    registrar_employee_id           UUID NOT NULL
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    source_patient_zitadel_user_id  TEXT,
    source_buffer_id                UUID,
    reopened_from_incident_id       UUID
        REFERENCES domain.incidents(id) ON DELETE RESTRICT,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX incidents_organization_id_idx ON domain.incidents (organization_id);
CREATE INDEX incidents_clinic_id_idx ON domain.incidents (clinic_id);
CREATE INDEX incidents_department_id_idx ON domain.incidents (department_id);
CREATE INDEX incidents_registrar_employee_id_idx ON domain.incidents (registrar_employee_id);
CREATE INDEX incidents_status_idx ON domain.incidents (status);
CREATE INDEX incidents_source_patient_idx
    ON domain.incidents (source_patient_zitadel_user_id)
    WHERE source_patient_zitadel_user_id IS NOT NULL;
CREATE INDEX incidents_source_buffer_id_idx
    ON domain.incidents (source_buffer_id)
    WHERE source_buffer_id IS NOT NULL;

CREATE TABLE domain.patient_incident_buffer (
    id                      UUID PRIMARY KEY,
    organization_id         UUID NOT NULL
        REFERENCES domain.organizations(id) ON DELETE RESTRICT,
    patient_zitadel_user_id TEXT NOT NULL,
    category_id             UUID
        REFERENCES domain.incident_categories(id) ON DELETE RESTRICT,
    type_id                 UUID
        REFERENCES domain.incident_types(id) ON DELETE RESTRICT,
    description             TEXT,
    occurred_at             TIMESTAMPTZ,
    status                  domain.buffer_status NOT NULL DEFAULT 'pending',
    published_incident_id   UUID
        REFERENCES domain.incidents(id) ON DELETE RESTRICT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX buffer_patient_zitadel_idx
    ON domain.patient_incident_buffer (patient_zitadel_user_id);
CREATE INDEX buffer_organization_id_idx
    ON domain.patient_incident_buffer (organization_id);
CREATE INDEX buffer_status_idx
    ON domain.patient_incident_buffer (status);

-- migrate:down
DROP TABLE IF EXISTS domain.patient_incident_buffer;
DROP TABLE IF EXISTS domain.incidents;
