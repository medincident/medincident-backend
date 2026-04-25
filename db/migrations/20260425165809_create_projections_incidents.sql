-- migrate:up
CREATE TABLE projections.incidents (
    id                              UUID PRIMARY KEY,
    organization_id                 UUID NOT NULL,
    clinic_id                       UUID NOT NULL,
    department_id                   UUID NOT NULL,
    category_id                     UUID NOT NULL,
    type_id                         UUID NOT NULL,
    status                          domain.incident_status NOT NULL DEFAULT 'pending',
    priority                        domain.incident_priority NOT NULL DEFAULT 'normal',
    description                     TEXT,
    patient_original_description    TEXT,
    occurred_at                     TIMESTAMPTZ NOT NULL,
    registrar_employee_id           UUID NOT NULL,
    registrar_display_name          TEXT NOT NULL,
    registrar_position              TEXT,
    registrar_organization_id       UUID NOT NULL,
    registrar_clinic_id             UUID NOT NULL,
    registrar_department_id         UUID NOT NULL,
    source_patient_zitadel_user_id  TEXT,
    source_buffer_id                UUID,
    reopened_from_incident_id       UUID,
    created_at                      TIMESTAMPTZ NOT NULL,
    updated_at                      TIMESTAMPTZ NOT NULL
);

CREATE INDEX proj_incidents_org_status_idx
    ON projections.incidents (organization_id, status);
CREATE INDEX proj_incidents_clinic_id_idx
    ON projections.incidents (clinic_id);
CREATE INDEX proj_incidents_dept_id_idx
    ON projections.incidents (department_id);
CREATE INDEX proj_incidents_registrar_idx
    ON projections.incidents (registrar_employee_id);
CREATE INDEX proj_incidents_source_patient_idx
    ON projections.incidents (source_patient_zitadel_user_id)
    WHERE source_patient_zitadel_user_id IS NOT NULL;

-- migrate:down
DROP TABLE IF EXISTS projections.incidents;
