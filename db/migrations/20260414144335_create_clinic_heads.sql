-- migrate:up
CREATE TABLE domain.clinic_heads (
    clinic_id          UUID NOT NULL
        REFERENCES domain.clinics(id) ON DELETE RESTRICT,
    employee_id        UUID NOT NULL
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    deputy_employee_id UUID
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (clinic_id, employee_id)
);

CREATE UNIQUE INDEX clinic_heads_holder_uniq
    ON domain.clinic_heads (employee_id);

CREATE INDEX clinic_heads_deputy_employee_id_idx
    ON domain.clinic_heads (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

-- migrate:down
DROP TABLE IF EXISTS domain.clinic_heads;
