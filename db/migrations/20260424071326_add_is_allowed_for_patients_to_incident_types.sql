-- migrate:up
ALTER TABLE domain.incident_types
    ADD COLUMN is_allowed_for_patients BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE projections.incident_types
    ADD COLUMN is_allowed_for_patients BOOLEAN NOT NULL DEFAULT FALSE;

-- migrate:down
ALTER TABLE projections.incident_types
    DROP COLUMN is_allowed_for_patients;

ALTER TABLE domain.incident_types
    DROP COLUMN is_allowed_for_patients;
