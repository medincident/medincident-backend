-- migrate:up
CREATE TABLE domain.org_dispatchers (
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

CREATE UNIQUE INDEX org_dispatchers_holder_uniq
    ON domain.org_dispatchers (employee_id);

CREATE INDEX org_dispatchers_deputy_employee_id_idx
    ON domain.org_dispatchers (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

-- migrate:down
DROP TABLE IF EXISTS domain.org_dispatchers;
