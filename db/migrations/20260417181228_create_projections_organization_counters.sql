-- migrate:up
CREATE TABLE projections.organization_counters (
    organization_id   uuid PRIMARY KEY,
    employees_total   int NOT NULL DEFAULT 0 CHECK (employees_total >= 0),
    clinics_total     int NOT NULL DEFAULT 0 CHECK (clinics_total >= 0),
    departments_total int NOT NULL DEFAULT 0 CHECK (departments_total >= 0),
    updated_at        timestamptz NOT NULL DEFAULT now()
);

-- migrate:down
DROP TABLE projections.organization_counters;
