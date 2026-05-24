-- migrate:up

CREATE TYPE domain.buffer_priority AS ENUM ('normal', 'high');

-- domain table: make description NOT NULL, add summary + priority
UPDATE domain.patient_incident_buffer SET description = '' WHERE description IS NULL;
ALTER TABLE domain.patient_incident_buffer ALTER COLUMN description SET NOT NULL;

ALTER TABLE domain.patient_incident_buffer
    ADD COLUMN summary  TEXT                  NOT NULL DEFAULT '',
    ADD COLUMN priority domain.buffer_priority NOT NULL DEFAULT 'normal';

UPDATE domain.patient_incident_buffer SET summary = description;

ALTER TABLE domain.patient_incident_buffer ALTER COLUMN summary  DROP DEFAULT;
ALTER TABLE domain.patient_incident_buffer ALTER COLUMN priority DROP DEFAULT;

-- projection table: same changes
UPDATE projections.patient_incident_buffer SET description = '' WHERE description IS NULL;
ALTER TABLE projections.patient_incident_buffer ALTER COLUMN description SET NOT NULL;

ALTER TABLE projections.patient_incident_buffer
    ADD COLUMN summary  TEXT                  NOT NULL DEFAULT '',
    ADD COLUMN priority domain.buffer_priority NOT NULL DEFAULT 'normal';

UPDATE projections.patient_incident_buffer SET summary = description;

ALTER TABLE projections.patient_incident_buffer ALTER COLUMN summary  DROP DEFAULT;
ALTER TABLE projections.patient_incident_buffer ALTER COLUMN priority DROP DEFAULT;

-- migrate:down

ALTER TABLE projections.patient_incident_buffer DROP COLUMN IF EXISTS summary, DROP COLUMN IF EXISTS priority;
ALTER TABLE projections.patient_incident_buffer ALTER COLUMN description DROP NOT NULL;

ALTER TABLE domain.patient_incident_buffer DROP COLUMN IF EXISTS summary, DROP COLUMN IF EXISTS priority;
ALTER TABLE domain.patient_incident_buffer ALTER COLUMN description DROP NOT NULL;

DROP TYPE IF EXISTS domain.buffer_priority;
