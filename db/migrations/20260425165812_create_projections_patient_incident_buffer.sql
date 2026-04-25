-- migrate:up
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

-- migrate:down
DROP TABLE IF EXISTS projections.patient_incident_buffer;
