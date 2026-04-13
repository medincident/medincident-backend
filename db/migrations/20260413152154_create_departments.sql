-- migrate:up
CREATE TABLE domain.departments (
    id          UUID PRIMARY KEY,
    clinic_id   UUID        NOT NULL
        REFERENCES domain.clinics(id) ON DELETE RESTRICT,
    name        TEXT        NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX departments_clinic_id_idx
    ON domain.departments (clinic_id);

-- migrate:down
DROP TABLE IF EXISTS domain.departments;
