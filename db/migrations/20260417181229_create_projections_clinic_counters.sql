-- migrate:up
CREATE TABLE projections.clinic_counters (
    clinic_id         uuid PRIMARY KEY,
    organization_id   uuid NOT NULL,
    employees_total   int NOT NULL DEFAULT 0 CHECK (employees_total >= 0),
    departments_total int NOT NULL DEFAULT 0 CHECK (departments_total >= 0),
    updated_at        timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX clinic_counters_org_idx
    ON projections.clinic_counters (organization_id);

-- migrate:down
DROP TABLE projections.clinic_counters;
