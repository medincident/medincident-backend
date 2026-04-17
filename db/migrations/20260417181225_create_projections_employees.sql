-- migrate:up
CREATE TABLE projections.employees (
    id              uuid        PRIMARY KEY,
    zitadel_user_id text        NOT NULL,
    organization_id uuid        NOT NULL,
    clinic_id       uuid        NULL,
    department_id   uuid        NOT NULL,
    position        text        NULL,
    hired_at        timestamptz NOT NULL,
    terminated_at   timestamptz NULL,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX employees_zitadel_user_id_idx
    ON projections.employees (zitadel_user_id);

CREATE INDEX employees_org_dept_idx
    ON projections.employees (organization_id, department_id);

CREATE INDEX employees_active_org_idx
    ON projections.employees (organization_id)
    WHERE terminated_at IS NULL;

CREATE INDEX employees_pending_backfill_idx
    ON projections.employees (department_id)
    WHERE clinic_id IS NULL;

-- migrate:down
DROP TABLE IF EXISTS projections.employees;
