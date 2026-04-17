-- migrate:up

-- Read-model copy of domain.org_dispatchers.
CREATE TABLE projections.org_dispatchers (
    organization_id    UUID NOT NULL,
    employee_id        UUID NOT NULL,
    deputy_employee_id UUID,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (organization_id, employee_id)
);

CREATE INDEX org_dispatchers_employee_id_idx
    ON projections.org_dispatchers (employee_id);
CREATE INDEX org_dispatchers_deputy_employee_id_idx
    ON projections.org_dispatchers (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

-- migrate:down
DROP TABLE IF EXISTS projections.org_dispatchers;
