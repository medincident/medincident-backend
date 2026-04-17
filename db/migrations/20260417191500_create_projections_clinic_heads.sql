-- migrate:up

-- Read-model copy of domain.clinic_heads. Matches the domain schema
-- shape except it lives under the projections.* namespace and has no
-- cross-schema foreign keys (see projections.clinics migration for the
-- rationale).
CREATE TABLE projections.clinic_heads (
    clinic_id          UUID NOT NULL,
    employee_id        UUID NOT NULL,
    deputy_employee_id UUID,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (clinic_id, employee_id)
);

CREATE INDEX clinic_heads_employee_id_idx
    ON projections.clinic_heads (employee_id);
CREATE INDEX clinic_heads_deputy_employee_id_idx
    ON projections.clinic_heads (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

-- migrate:down
DROP TABLE IF EXISTS projections.clinic_heads;
