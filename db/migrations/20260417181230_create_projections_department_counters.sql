-- migrate:up
CREATE TABLE projections.department_counters (
    department_id   uuid PRIMARY KEY,
    clinic_id       uuid NULL,
    organization_id uuid NOT NULL,
    employees_total int NOT NULL DEFAULT 0,
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX department_counters_clinic_idx
    ON projections.department_counters (clinic_id);

CREATE INDEX department_counters_org_idx
    ON projections.department_counters (organization_id);

-- migrate:down
DROP TABLE IF EXISTS projections.department_counters;
