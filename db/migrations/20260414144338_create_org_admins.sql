-- migrate:up
CREATE TABLE domain.org_admins (
    organization_id    UUID NOT NULL
        REFERENCES domain.organizations(id) ON DELETE RESTRICT,
    employee_id        UUID NOT NULL
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    deputy_employee_id UUID
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (organization_id, employee_id)
);

CREATE UNIQUE INDEX org_admins_holder_uniq
    ON domain.org_admins (employee_id);

CREATE INDEX org_admins_deputy_employee_id_idx
    ON domain.org_admins (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

-- migrate:down
DROP TABLE IF EXISTS domain.org_admins;
