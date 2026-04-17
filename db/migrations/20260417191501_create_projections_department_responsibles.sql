-- migrate:up

-- Read-model copy of domain.department_responsibles.
CREATE TABLE projections.department_responsibles (
    department_id      UUID NOT NULL,
    employee_id        UUID NOT NULL,
    deputy_employee_id UUID,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (department_id, employee_id)
);

CREATE INDEX department_responsibles_employee_id_idx
    ON projections.department_responsibles (employee_id);
CREATE INDEX department_responsibles_deputy_employee_id_idx
    ON projections.department_responsibles (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

-- migrate:down
DROP TABLE IF EXISTS projections.department_responsibles;
