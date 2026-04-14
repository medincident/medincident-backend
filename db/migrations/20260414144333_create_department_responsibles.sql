-- migrate:up
CREATE TABLE domain.department_responsibles (
    department_id      UUID NOT NULL
        REFERENCES domain.departments(id) ON DELETE RESTRICT,
    employee_id        UUID NOT NULL
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    deputy_employee_id UUID
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (department_id, employee_id)
);

CREATE UNIQUE INDEX department_responsibles_holder_uniq
    ON domain.department_responsibles (employee_id);

CREATE INDEX department_responsibles_deputy_employee_id_idx
    ON domain.department_responsibles (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

-- migrate:down
DROP TABLE IF EXISTS domain.department_responsibles;
