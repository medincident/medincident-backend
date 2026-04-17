-- migrate:up

-- Read-model copy of domain.org_admins.
CREATE TABLE projections.org_admins (
    organization_id    UUID NOT NULL,
    employee_id        UUID NOT NULL,
    deputy_employee_id UUID,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (organization_id, employee_id)
);

CREATE INDEX org_admins_employee_id_idx
    ON projections.org_admins (employee_id);
CREATE INDEX org_admins_deputy_employee_id_idx
    ON projections.org_admins (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

-- migrate:down
DROP TABLE IF EXISTS projections.org_admins;
