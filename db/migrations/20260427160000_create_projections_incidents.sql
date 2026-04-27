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
CREATE INDEX proj_incidents_source_buffer_id_idx
    ON projections.incidents (source_buffer_id)
    WHERE source_buffer_id IS NOT NULL;

CREATE TABLE projections.patient_incident_buffer (
    id                      UUID PRIMARY KEY,
    organization_id         UUID NOT NULL,
    patient_zitadel_user_id TEXT NOT NULL,
    category_id             UUID,
    type_id                 UUID,
    description             TEXT,
    occurred_at             TIMESTAMPTZ,
    status                  domain.buffer_status NOT NULL DEFAULT 'pending',
    published_incident_id   UUID,
    created_at              TIMESTAMPTZ NOT NULL,
    updated_at              TIMESTAMPTZ NOT NULL
);

CREATE INDEX proj_buffer_patient_idx
    ON projections.patient_incident_buffer (patient_zitadel_user_id);
CREATE INDEX proj_buffer_org_status_idx
    ON projections.patient_incident_buffer (organization_id, status);

CREATE TABLE projections.incident_status_history (
    id                  UUID PRIMARY KEY,
    incident_id         UUID NOT NULL,
    old_status          domain.incident_status,
    new_status          domain.incident_status NOT NULL,
    actor_employee_id   UUID,
    actor_display_name  TEXT,
    changed_at          TIMESTAMPTZ NOT NULL
);

CREATE INDEX proj_status_hist_incident_idx
    ON projections.incident_status_history (incident_id, changed_at);

CREATE TABLE projections.incident_priority_history (
    id                  UUID PRIMARY KEY,
    incident_id         UUID NOT NULL,
    old_priority        domain.incident_priority NOT NULL,
    new_priority        domain.incident_priority NOT NULL,
    actor_employee_id   UUID,
    actor_display_name  TEXT,
    changed_at          TIMESTAMPTZ NOT NULL
);

CREATE INDEX proj_priority_hist_incident_idx
    ON projections.incident_priority_history (incident_id, changed_at);

-- migrate:down
DROP TABLE IF EXISTS projections.incident_priority_history;
DROP TABLE IF EXISTS projections.incident_status_history;
DROP TABLE IF EXISTS projections.patient_incident_buffer;
DROP TABLE IF EXISTS projections.incidents;
