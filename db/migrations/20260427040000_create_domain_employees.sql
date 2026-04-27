-- migrate:up

CREATE TABLE domain.employees (
    id               UUID PRIMARY KEY,
    zitadel_user_id  TEXT        NOT NULL,
    organization_id  UUID        NOT NULL
        REFERENCES domain.organizations(id) ON DELETE RESTRICT,
    department_id    UUID        NOT NULL
        REFERENCES domain.departments(id) ON DELETE RESTRICT,
    position         TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX employees_org_user_uniq
    ON domain.employees (organization_id, zitadel_user_id);

CREATE INDEX employees_department_id_idx
    ON domain.employees (department_id);

CREATE INDEX employees_zitadel_user_id_idx
    ON domain.employees (zitadel_user_id);

-- migrate:down
DROP TABLE IF EXISTS domain.employees;
