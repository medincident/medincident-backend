-- migrate:up
CREATE TABLE projections.organization_counters (
    organization_id   uuid PRIMARY KEY,
    employees_total   int NOT NULL DEFAULT 0,
    clinics_total     int NOT NULL DEFAULT 0,
    departments_total int NOT NULL DEFAULT 0,
    updated_at        timestamptz NOT NULL DEFAULT now()
);

-- migrate:down
DROP TABLE IF EXISTS projections.organization_counters;
