-- migrate:up
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
